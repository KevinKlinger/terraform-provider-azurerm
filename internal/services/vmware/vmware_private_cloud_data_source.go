package vmware

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/location"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/vmware/sdk/2020-03-20/privateclouds"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
)

func dataSourceVmwarePrivateCloud() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVmwarePrivateCloudRead,

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

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"management_cluster": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"size": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"hosts": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"id": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"network_subnet_cidr": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"internet_connection_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"circuit": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"express_route_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"express_route_private_peering_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"primary_subnet_cidr": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"secondary_subnet_cidr": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"hcx_cloud_manager_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"management_subnet_cidr": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"nsxt_certificate_thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"nsxt_manager_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"provisioning_subnet_cidr": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"vcenter_certificate_thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"vcsa_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"vmotion_subnet_cidr": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceVmwarePrivateCloudRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := privateclouds.NewPrivateCloudID(subscriptionId, resourceGroup, name)
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
		d.Set("location", location.Normalize(model.Location))

		props := model.Properties
		if err := d.Set("management_cluster", flattenPrivateCloudManagementCluster(props.ManagementCluster)); err != nil {
			return fmt.Errorf("setting `management_cluster`: %+v", err)
		}
		d.Set("network_subnet_cidr", props.NetworkBlock)
		if err := d.Set("circuit", flattenPrivateCloudCircuit(props.Circuit)); err != nil {
			return fmt.Errorf("setting `circuit`: %+v", err)
		}

		internetConnectionEnabled := false
		if props.Internet != nil {
			internetConnectionEnabled = *props.Internet == privateclouds.InternetEnumEnabled
		}

		d.Set("internet_connection_enabled", internetConnectionEnabled)
		d.Set("hcx_cloud_manager_endpoint", props.Endpoints.HcxCloudManager)
		d.Set("nsxt_manager_endpoint", props.Endpoints.NsxtManager)
		d.Set("vcsa_endpoint", props.Endpoints.Vcsa)
		d.Set("management_subnet_cidr", props.ManagementNetwork)
		d.Set("nsxt_certificate_thumbprint", props.NsxtCertificateThumbprint)
		d.Set("provisioning_subnet_cidr", props.ProvisioningNetwork)
		d.Set("vcenter_certificate_thumbprint", props.VcenterCertificateThumbprint)
		d.Set("vmotion_subnet_cidr", props.VmotionNetwork)

		d.Set("sku_name", model.Sku.Name)

		if err := tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}
