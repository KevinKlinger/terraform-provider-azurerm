package subscription

import (
	"fmt"
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func dataSourceSubscription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSubscriptionRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*pluginsdk.Schema{
			"subscription_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenant_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location_placement_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"quota_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"spending_limit": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client)
	groupClient := client.Subscription.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId := d.Get("subscription_id").(string)
	if subscriptionId == "" {
		subscriptionId = client.Account.SubscriptionId
	}

	resp, err := groupClient.Get(ctx, subscriptionId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: default tags for Subscription %q was not found", subscriptionId)
		}

		return fmt.Errorf("reading default tags for Subscription: %+v", err)
	}

	d.SetId(*resp.ID)
	d.Set("subscription_id", resp.SubscriptionID)
	d.Set("display_name", resp.DisplayName)
	d.Set("tenant_id", resp.TenantID)
	d.Set("state", resp.State)
	if resp.SubscriptionPolicies != nil {
		d.Set("location_placement_id", resp.SubscriptionPolicies.LocationPlacementID)
		d.Set("quota_id", resp.SubscriptionPolicies.QuotaID)
		d.Set("spending_limit", resp.SubscriptionPolicies.SpendingLimit)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
