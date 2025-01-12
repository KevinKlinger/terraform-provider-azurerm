package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storagesync/mgmt/2020-03-01/storagesync"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/tf"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/storage/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/storage/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func resourceStorageSyncGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageSyncGroupCreate,
		Read:   resourceStorageSyncGroupRead,
		Delete: resourceStorageSyncGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.StorageSyncGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageSyncName,
			},

			"storage_sync_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageSyncId,
			},
		},
	}
}

func resourceStorageSyncGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	ssId, _ := parse.StorageSyncServiceID(d.Get("storage_sync_id").(string))

	existing, err := client.Get(ctx, ssId.ResourceGroup, ssId.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q): %+v", name, ssId.Name, ssId.ResourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_storage_sync_group", *existing.ID)
	}

	if _, err := client.Create(ctx, ssId.ResourceGroup, ssId.Name, name, storagesync.SyncGroupCreateParameters{}); err != nil {
		return fmt.Errorf("creating Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q): %+v", name, ssId.Name, ssId.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, ssId.ResourceGroup, ssId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q): %+v", name, ssId.Name, ssId.ResourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("reading Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q) ID is empty or nil", name, ssId.Name, ssId.ResourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceStorageSyncGroupRead(d, meta)
}

func resourceStorageSyncGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncGroupsClient
	ssClient := meta.(*clients.Client).Storage.SyncServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageSyncGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.StorageSyncServiceName, id.SyncGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Storage Sync Group %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q): %+v", id.SyncGroupName, id.StorageSyncServiceName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)

	ssResp, err := ssClient.Get(ctx, id.ResourceGroup, id.StorageSyncServiceName)
	if err != nil {
		return fmt.Errorf("reading Storage Sync %q (Resource Group %q): %+v", id.StorageSyncServiceName, id.ResourceGroup, err)
	}

	if ssResp.ID == nil || *ssResp.ID == "" {
		return fmt.Errorf("reading Storage Sync %q (Resource Group %q) ID is empty or nil", id.StorageSyncServiceName, id.ResourceGroup)
	}

	d.Set("storage_sync_id", ssResp.ID)

	return nil
}

func resourceStorageSyncGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageSyncGroupID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.StorageSyncServiceName, id.SyncGroupName); err != nil {
		return fmt.Errorf("deleting Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q): %+v", id.SyncGroupName, id.StorageSyncServiceName, id.ResourceGroup, err)
	}

	return nil
}
