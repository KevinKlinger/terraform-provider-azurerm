package recoveryservices

import (
	"fmt"
	"log"
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func dataSourceRecoveryServicesVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceRecoveryServicesVaultRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"tags": tags.SchemaDataSource(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRecoveryServicesVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	log.Printf("[DEBUG] Reading Recovery Service Vault %q (resource group %q)", name, resourceGroup)

	vault, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(vault.Response) {
			return fmt.Errorf("Error: Recovery Services Vault %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("making Read request on Recovery Service Vault %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*vault.ID)
	d.Set("name", vault.Name)
	d.Set("location", azure.NormalizeLocation(*vault.Location))
	d.Set("resource_group_name", resourceGroup)

	if sku := vault.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	return tags.FlattenAndSet(d, vault.Tags)
}
