package compute

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func dataSourceImage() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceImageRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{

			"name_regex": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				ConflictsWith: []string{"name"},
			},
			"sort_descending": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"name": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name_regex"},
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"zone_resilient": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"os_disk": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"blob_uri": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"caching": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"managed_disk_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"os_state": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"size_gb": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"data_disk": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"blob_uri": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"caching": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"lun": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"managed_disk_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"size_gb": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceImageRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ImagesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)

	name := d.Get("name").(string)
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if name == "" && !nameRegexOk {
		return fmt.Errorf("[ERROR] either name or name_regex is required")
	}

	var img compute.Image

	if !nameRegexOk {
		var err error
		if img, err = client.Get(ctx, resGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(img.Response) {
				return fmt.Errorf("image %q was not found in resource group %q", name, resGroup)
			}
			return fmt.Errorf("[ERROR] Error making Read request on Azure Image %q (resource group %q): %+v", name, resGroup, err)
		}
	} else {
		r := regexp.MustCompile(nameRegex.(string))

		list := make([]compute.Image, 0)
		resp, err := client.ListByResourceGroupComplete(ctx, resGroup)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response().Response) {
				return fmt.Errorf("No Images were found for Resource Group %q", resGroup)
			}
			return fmt.Errorf("[ERROR] Error getting list of images (resource group %q): %+v", resGroup, err)
		}

		for resp.NotDone() {
			img = resp.Value()
			if r.Match(([]byte)(*img.Name)) {
				list = append(list, img)
			}
			err = resp.NextWithContext(ctx)

			if err != nil {
				return err
			}
		}

		if 1 > len(list) {
			return fmt.Errorf("No Images were found for Resource Group %q", resGroup)
		}

		if len(list) > 1 {
			desc := d.Get("sort_descending").(bool)
			log.Printf("[DEBUG] Image - multiple results found and `sort_descending` is set to: %t", desc)

			sort.Slice(list, func(i, j int) bool {
				return (!desc && *list[i].Name < *list[j].Name) ||
					(desc && *list[i].Name > *list[j].Name)
			})
		}
		img = list[0]
	}

	d.SetId(*img.ID)
	d.Set("name", img.Name)
	d.Set("resource_group_name", resGroup)
	if location := img.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if profile := img.StorageProfile; profile != nil {
		if disk := profile.OsDisk; disk != nil {
			if err := d.Set("os_disk", flattenAzureRmImageOSDisk(disk)); err != nil {
				return fmt.Errorf("[DEBUG] Error setting AzureRM Image OS Disk error: %+v", err)
			}
		}

		if disks := profile.DataDisks; disks != nil {
			if err := d.Set("data_disk", flattenAzureRmImageDataDisks(disks)); err != nil {
				return fmt.Errorf("[DEBUG] Error setting AzureRM Image Data Disks error: %+v", err)
			}
		}

		d.Set("zone_resilient", profile.ZoneResilient)
	}

	return tags.FlattenAndSet(d, img.Tags)
}
