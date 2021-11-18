package resource

import (
	"fmt"
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/resource/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/resource/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func dataSourceTemplateSpecVersion() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceTemplateSpecVersionRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		//lintignore:S033
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.TemplateSpecName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"version": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.TemplateSpecVersionName,
			},

			"template_body": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceTemplateSpecVersionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TemplateSpecsVersionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewTemplateSpecVersionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string), d.Get("version").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.TemplateSpecName, id.VersionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Templatespec %q, with version name %q (resource group %q) was not found: %+v", id.TemplateSpecName, id.VersionName, id.ResourceGroup, err)
		}
		return fmt.Errorf("reading Templatespec %q, with version name %q (resource group %q): %+v", id.TemplateSpecName, id.VersionName, id.ResourceGroup, err)
	}

	templateBody := "{}"
	if props := resp.VersionProperties; props != nil && props.Template != nil {
		templateBodyRaw, err := flattenTemplateDeploymentBody(props.Template)
		if err != nil {
			return err
		}

		templateBody = *templateBodyRaw
	}
	d.Set("template_body", templateBody)

	d.SetId(id.ID())

	return tags.FlattenAndSet(d, resp.Tags)
}
