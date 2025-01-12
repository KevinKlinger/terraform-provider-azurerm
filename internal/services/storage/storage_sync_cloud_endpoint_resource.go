package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storagesync/mgmt/2020-03-01/storagesync"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/tf"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/storage/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/storage/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func resourceStorageSyncCloudEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageSyncCloudEndpointCreate,
		Read:   resourceStorageSyncCloudEndpointRead,
		Delete: resourceStorageSyncCloudEndpointDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.StorageSyncCloudEndpointID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(45 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageSyncName,
			},

			"storage_sync_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageSyncGroupID,
			},

			"file_share_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageShareName,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountID,
			},

			"storage_account_tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceStorageSyncCloudEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.CloudEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	storagesyncGroupId, _ := parse.StorageSyncGroupID(d.Get("storage_sync_group_id").(string))

	existing, err := client.Get(ctx, storagesyncGroupId.ResourceGroup, storagesyncGroupId.StorageSyncServiceName, storagesyncGroupId.SyncGroupName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Storage Sync Cloud Endpoint %q (Storage Sync Group %q / Storage Sync Name %q / Resource Group %q): %+v", name, storagesyncGroupId.SyncGroupName, storagesyncGroupId.StorageSyncServiceName, storagesyncGroupId.ResourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_storage_sync_cloud_endpoint", *existing.ID)
	}

	parameters := storagesync.CloudEndpointCreateParameters{
		CloudEndpointCreateParametersProperties: &storagesync.CloudEndpointCreateParametersProperties{
			StorageAccountResourceID: utils.String(d.Get("storage_account_id").(string)),
			AzureFileShareName:       utils.String(d.Get("file_share_name").(string)),
		},
	}

	tenantId := meta.(*clients.Client).Account.TenantId
	if v, ok := d.GetOk("storage_account_tenant_id"); ok {
		tenantId = v.(string)
	}
	parameters.CloudEndpointCreateParametersProperties.StorageAccountTenantID = &tenantId

	future, err := client.Create(ctx, storagesyncGroupId.ResourceGroup, storagesyncGroupId.StorageSyncServiceName, storagesyncGroupId.SyncGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Storage Sync Cloud Endpoint %q (Storage Sync Group %q / Storage Sync %q / Resource Group %q): %+v", name, storagesyncGroupId.SyncGroupName, storagesyncGroupId.StorageSyncServiceName, storagesyncGroupId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for Storage Sync Cloud Endpoint %q to be created: %+v", name, err)
	}

	resp, err := client.Get(ctx, storagesyncGroupId.ResourceGroup, storagesyncGroupId.StorageSyncServiceName, storagesyncGroupId.SyncGroupName, name)
	if err != nil {
		return fmt.Errorf("retrieving Storage Sync Cloud Endpoint %q (Storage Sync Group %q / Storage Sync %q / Resource Group %q): %+v", name, storagesyncGroupId.SyncGroupName, storagesyncGroupId.StorageSyncServiceName, storagesyncGroupId.ResourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("reading Storage Sync Cloud Endpoint %q (Storage Sync Group %q / Storage Sync %q / Resource Group %q) ID is nil or empty", name, storagesyncGroupId.SyncGroupName, storagesyncGroupId.StorageSyncServiceName, storagesyncGroupId.ResourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceStorageSyncCloudEndpointRead(d, meta)
}

func resourceStorageSyncCloudEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.CloudEndpointsClient
	gpClient := meta.(*clients.Client).Storage.SyncGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageSyncCloudEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.StorageSyncServiceName, id.SyncGroupName, id.CloudEndpointName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Storage Sync Cloud Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Storage Sync Cloud Endpoint %q (Storage Sync Group %q / Storage Sync %q / Resource Group %q): %+v", id.CloudEndpointName, id.SyncGroupName, id.StorageSyncServiceName, id.ResourceGroup, err)
	}
	d.Set("name", resp.Name)

	gpResp, err := gpClient.Get(ctx, id.ResourceGroup, id.StorageSyncServiceName, id.SyncGroupName)
	if err != nil {
		return fmt.Errorf("reading Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q): %+v", id.SyncGroupName, id.StorageSyncServiceName, id.ResourceGroup, err)
	}

	if gpResp.ID == nil || *gpResp.ID == "" {
		return fmt.Errorf("reading Storage Sync Group %q (Resource Group %q) ID is empty or nil", id.SyncGroupName, id.ResourceGroup)
	}

	d.Set("storage_sync_group_id", gpResp.ID)
	if props := resp.CloudEndpointProperties; props != nil {
		d.Set("storage_account_id", props.StorageAccountResourceID)
		d.Set("file_share_name", props.AzureFileShareName)
		d.Set("storage_account_tenant_id", props.StorageAccountTenantID)
	}

	return nil
}

func resourceStorageSyncCloudEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.CloudEndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageSyncCloudEndpointID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.StorageSyncServiceName, id.SyncGroupName, id.CloudEndpointName)
	if err != nil {
		return fmt.Errorf("deleting Storage Sync Cloud Endpoint %q (Storage Sync Group %q / Storage Sync %q / Resource Group %q): %+v", id.CloudEndpointName, id.SyncGroupName, id.StorageSyncServiceName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of Storage Sync Cloud Endpoint %q (Storage Sync Group %q / Storage Sync %q / Resource Group %q): %+v", id.CloudEndpointName, id.SyncGroupName, id.StorageSyncServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}
