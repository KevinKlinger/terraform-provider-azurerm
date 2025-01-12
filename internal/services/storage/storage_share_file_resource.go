package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/tf"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/storage/parse"
	storageValidate "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/storage/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/files"
)

func resourceStorageShareFile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageShareFileCreate,
		Read:   resourceStorageShareFileRead,
		Update: resourceStorageShareFileUpdate,
		Delete: resourceStorageShareFileDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_share_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageShareID,
			},
			"path": {
				Type:         pluginsdk.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "",
				ValidateFunc: storageValidate.StorageShareDirectoryName,
			},

			"content_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "application/octet-stream",
			},

			"content_encoding": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"content_md5": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"content_disposition": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"source": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ForceNew:     true,
			},

			"metadata": MetaDataSchema(),
		},
	}
}

func resourceStorageShareFileCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	storageShareID, err := parse.StorageShareDataPlaneID(d.Get("storage_share_id").(string))
	if err != nil {
		return err
	}

	fileName := d.Get("name").(string)
	path := d.Get("path").(string)

	account, err := storageClient.FindAccount(ctx, storageShareID.AccountName)
	if err != nil {
		return fmt.Errorf("eretrieving Account %q for File %q (Share %q): %s", storageShareID.AccountName, fileName, storageShareID.Name, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Storage Account %q!", storageShareID.AccountName)
	}

	fileSharesClient, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building File Share Directories Client: %s", err)
	}

	share, err := fileSharesClient.Get(ctx, account.ResourceGroup, storageShareID.AccountName, storageShareID.Name)
	if err != nil {
		return fmt.Errorf("retrieving Share %q for File %q: %s", storageShareID.Name, fileName, err)
	}
	if share == nil {
		return fmt.Errorf("unable to locate Storage Share %q", storageShareID.Name)
	}

	client, err := storageClient.FileShareFilesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building File Share Directories Client: %s", err)
	}

	existing, err := client.GetProperties(ctx, storageShareID.AccountName, storageShareID.Name, path, fileName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing File %q (File Share %q / Storage Account %q / Resource Group %q): %s", fileName, storageShareID.Name, storageShareID.AccountName, account.ResourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		id := client.GetResourceID(storageShareID.AccountName, storageShareID.Name, path, fileName)
		return tf.ImportAsExistsError("azurerm_storage_share_file", id)
	}

	input := files.CreateInput{
		MetaData:           ExpandMetaData(d.Get("metadata").(map[string]interface{})),
		ContentType:        utils.String(d.Get("content_type").(string)),
		ContentEncoding:    utils.String(d.Get("content_encoding").(string)),
		ContentDisposition: utils.String(d.Get("content_disposition").(string)),
	}

	if v, ok := d.GetOk("content_md5"); ok {
		input.ContentMD5 = utils.String(v.(string))
	}

	var file *os.File
	if v, ok := d.GetOk("source"); ok {
		file, err = os.Open(v.(string))
		if err != nil {
			return fmt.Errorf("opening file : %s", err)
		}

		info, err := file.Stat()
		if err != nil {
			return fmt.Errorf("'stat'-ing File %q (File Share %q / Account %q): %+v", fileName, storageShareID.Name, storageShareID.AccountName, err)
		}

		input.ContentLength = info.Size()
	}

	if _, err := client.Create(ctx, storageShareID.AccountName, storageShareID.Name, path, fileName, input); err != nil {
		return fmt.Errorf("creating File %q (File Share %q / Account %q): %+v", fileName, storageShareID.Name, storageShareID.AccountName, err)
	}

	if file != nil {
		if err := client.PutFile(ctx, storageShareID.AccountName, storageShareID.Name, path, fileName, file, 4); err != nil {
			return fmt.Errorf("uploading File: %q (File Share %q / Account %q): %+v", fileName, storageShareID.Name, storageShareID.AccountName, err)
		}
	}

	resourceID := client.GetResourceID(storageShareID.AccountName, storageShareID.Name, path, fileName)
	d.SetId(resourceID)

	return resourceStorageShareFileRead(d, meta)
}

func resourceStorageShareFileUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := files.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for File %q (Share %q): %s", id.AccountName, id.FileName, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Storage Account %q!", id.AccountName)
	}

	fileSharesClient, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building File Share Directories Client: %s", err)
	}

	share, err := fileSharesClient.Get(ctx, account.ResourceGroup, id.AccountName, id.ShareName)
	if err != nil {
		return fmt.Errorf("retrieving Share %q for File %q: %s", id.ShareName, id.FileName, err)
	}
	if share == nil {
		return fmt.Errorf("unable to locate Storage Share %q", id.ShareName)
	}

	client, err := storageClient.FileShareFilesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building File Share Files Client: %s", err)
	}

	existing, err := client.GetProperties(ctx, id.AccountName, id.ShareName, id.DirectoryName, id.FileName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing File %q (File Share %q / Storage Account %q / Resource Group %q): %s", id.FileName, id.ShareName, id.AccountName, account.ResourceGroup, err)
		}
	}

	if d.HasChange("content_type") || d.HasChange("content_encoding") || d.HasChange("content_disposition") || d.HasChange("content_md5") {
		input := files.SetPropertiesInput{
			ContentType:        utils.String(d.Get("content_type").(string)),
			ContentEncoding:    utils.String(d.Get("content_encoding").(string)),
			ContentDisposition: utils.String(d.Get("content_disposition").(string)),
			MetaData:           ExpandMetaData(d.Get("metadata").(map[string]interface{})),
		}

		if v, ok := d.GetOk("content_md5"); ok {
			input.ContentMD5 = utils.String(v.(string))
		}

		if _, err := client.SetProperties(ctx, id.AccountName, id.ShareName, id.DirectoryName, id.FileName, input); err != nil {
			return fmt.Errorf("creating File %q (File Share %q / Account %q): %+v", id.FileName, id.ShareName, id.AccountName, err)
		}
	}

	return resourceStorageShareFileRead(d, meta)
}

func resourceStorageShareFileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := files.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for File %q (Share %q): %s", id.AccountName, id.FileName, id.ShareName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Storage Account for Storage Share File %q (Share %s, Account %s) - assuming removed & removing from state", id.FileName, id.ShareName, id.AccountName)
		d.SetId("")
		return nil
	}

	fileSharesClient, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building File Share Directories Client: %s", err)
	}

	share, err := fileSharesClient.Get(ctx, account.ResourceGroup, id.AccountName, id.ShareName)
	if err != nil {
		return fmt.Errorf("retrieving Share %q for File %q: %s", id.ShareName, id.FileName, err)
	}
	if share == nil {
		log.Printf("[WARN] Unable to determine Storage Share for Storage Share File %q (Share %s, Account %s) - assuming removed & removing from state", id.FileName, id.ShareName, id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.FileShareFilesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	props, err := client.GetProperties(ctx, id.AccountName, id.ShareName, id.DirectoryName, id.FileName)
	if err != nil {
		log.Printf("Retrieving Storage File Share file %q (Directory %q / File Share %q / Account %q / Resource Group %q): %s", id.FileName, id.DirectoryName, id.ShareName, id.AccountName, account.ResourceGroup, err)
		d.SetId("")
		return nil
	}

	d.Set("name", id.FileName)
	d.Set("path", id.DirectoryName)
	d.Set("storage_share_id", parse.NewStorageShareDataPlaneId(id.AccountName, storageClient.Environment.StorageEndpointSuffix, id.ShareName).ID())

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %s", err)
	}
	d.Set("content_type", props.ContentType)
	d.Set("content_encoding", props.ContentEncoding)
	d.Set("content_md5", props.ContentMD5)
	d.Set("content_disposition", props.ContentDisposition)

	return nil
}

func resourceStorageShareFileDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := files.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for File %q (Share %q): %s", id.AccountName, id.FileName, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Storage Account %q", id.AccountName)
	}

	client, err := storageClient.FileShareFilesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building File Share File Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	if _, err := client.Delete(ctx, id.AccountName, id.ShareName, id.DirectoryName, id.FileName); err != nil {
		return fmt.Errorf("deleting Storage Share File %q (File Share %q / Account %q / Resource Group %q): %s", id.FileName, id.ShareName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}
