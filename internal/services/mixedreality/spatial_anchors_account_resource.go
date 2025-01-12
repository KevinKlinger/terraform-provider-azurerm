package mixedreality

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mixedreality/mgmt/2021-01-01/mixedreality"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/tf"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/mixedreality/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func resourceSpatialAnchorsAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSpatialAnchorsAccountCreateUpdate,
		Read:   resourceSpatialAnchorsAccountRead,
		Update: resourceSpatialAnchorsAccountCreateUpdate,
		Delete: resourceSpatialAnchorsAccountDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpatialAnchorsAccountID(id)
			return err
		}),

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
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[-\w._()]+$`),
					"Spatial Anchors Account name must be 1 - 90 characters long, contain only word characters and underscores.",
				),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_domain": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSpatialAnchorsAccountCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MixedReality.SpatialAnchorsAccountClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Spatial Anchors Account %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_spatial_anchors_account", *existing.ID)
		}
	}

	account := mixedreality.SpatialAnchorsAccount{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	if _, err := client.Create(ctx, resourceGroup, name, account); err != nil {
		return fmt.Errorf("creating/updating Spatial Anchors Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Spatial Anchors Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Spatial Anchors Account %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceSpatialAnchorsAccountRead(d, meta)
}

func resourceSpatialAnchorsAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MixedReality.SpatialAnchorsAccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpatialAnchorsAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Spatial Anchors Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.AccountProperties; props != nil {
		d.Set("account_domain", props.AccountDomain)
		d.Set("account_id", props.AccountID)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSpatialAnchorsAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MixedReality.SpatialAnchorsAccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpatialAnchorsAccountID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting Spatial Anchors Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
