package datalake

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/datalake/sdk/datalakestore/2016-11-01/accounts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
)

func dataSourceDataLakeStoreAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmDateLakeStoreAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"encryption_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"encryption_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"firewall_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"firewall_allow_azure_ips": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmDateLakeStoreAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	subscriptionId := meta.(*clients.Client).Datalake.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := accounts.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if properties := model.Properties; properties != nil {
			d.Set("tier", properties.CurrentTier)

			d.Set("encryption_state", properties.EncryptionState)
			d.Set("firewall_allow_azure_ips", properties.FirewallAllowAzureIps)
			d.Set("firewall_state", properties.FirewallState)

			if config := properties.EncryptionConfig; config != nil {
				d.Set("encryption_type", string(config.Type))
			}
		}

		return tags.FlattenAndSet(d, flattenTags(model.Tags))
	}
	return nil
}
