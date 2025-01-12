package dataprotection

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/identity"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/location"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/dataprotection/legacysdk/dataprotection"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/dataprotection/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func dataSourceDataProtectionBackupVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDataProtectionBackupVaultRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.BackupVaultID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{2,50}$"),
					"DataProtection BackupVault name must be 2 - 50 characters long, contain only letters, numbers and hyphens.).",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": location.SchemaComputed(),

			"datastore_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"redundancy": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceDataProtectionBackupVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewBackupVaultID(subscriptionId, resourceGroup, name)

	resp, err := client.Get(ctx, id.Name, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] DataProtection BackupVault %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupVault (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.Properties; props != nil {
		if props.StorageSettings != nil && len(*props.StorageSettings) > 0 {
			d.Set("datastore_type", (*props.StorageSettings)[0].DatastoreType)
			d.Set("redundancy", (*props.StorageSettings)[0].Type)
		}
	}
	if err := d.Set("identity", dataSourceFlattenBackupVaultDppIdentityDetails(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func dataSourceFlattenBackupVaultDppIdentityDetails(input *dataprotection.DppIdentityDetails) []interface{} {
	var config *identity.ExpandedConfig
	if input != nil {
		principalId := ""
		if input.PrincipalID != nil {
			principalId = *input.PrincipalID
		}

		tenantId := ""
		if input.TenantID != nil {
			tenantId = *input.TenantID
		}
		config = &identity.ExpandedConfig{
			Type:        identity.Type(*input.Type),
			PrincipalId: principalId,
			TenantId:    tenantId,
		}
	}
	return identity.SystemAssigned{}.Flatten(config)
}
