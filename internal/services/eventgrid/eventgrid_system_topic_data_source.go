package eventgrid

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/eventgrid/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func dataSourceEventGridSystemTopic() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceEventGridSystemTopicRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringMatch(
						regexp.MustCompile("^[-a-zA-Z0-9]{3,128}$"),
						"EventGrid Topics name must be 3 - 128 characters long, contain only letters, numbers and hyphens.",
					),
				),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"identity": IdentitySchemaForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"source_arm_resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"topic_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"metric_arm_resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceEventGridSystemTopicRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.SystemTopicsClient
	subscriptionId := meta.(*clients.Client).EventGrid.DomainsClient.SubscriptionID
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSystemTopicID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Event Grid System Topic '%s' was not found (resource group '%s')", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Event Grid System Topic '%s': %+v", id.Name, err)
	}

	d.SetId(id.ID())
	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.SystemTopicProperties; props != nil {
		d.Set("source_arm_resource_id", props.Source)
		d.Set("topic_type", props.TopicType)
		d.Set("metric_arm_resource_id", props.MetricResourceID)
	}

	if err := d.Set("identity", flattenIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
