package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/location"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/sdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/appservice/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/appservice/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type ServicePlanDataSource struct{}

var _ sdk.DataSource = ServicePlanDataSource{}

type ServicePlanDataSourceModel struct {
	Name                      string            `tfschema:"name"`
	ResourceGroup             string            `tfschema:"resource_group_name"`
	Location                  string            `tfschema:"location"`
	Kind                      string            `tfschema:"kind"`
	OSType                    OSType            `tfschema:"os_type"`
	Sku                       string            `tfschema:"sku_name"`
	AppServiceEnvironmentId   string            `tfschema:"app_service_environment_id"`
	PerSiteScaling            bool              `tfschema:"per_site_scaling_enabled"`
	Reserved                  bool              `tfschema:"reserved"`
	NumberOfWorkers           int               `tfschema:"number_of_workers"`
	MaximumElasticWorkerCount int               `tfschema:"maximum_elastic_worker_count"`
	Tags                      map[string]string `tfschema:"tags"`
}

func (r ServicePlanDataSource) ModelObject() interface{} {
	return &ServicePlanDataSourceModel{}
}

func (r ServicePlanDataSource) ResourceType() string {
	return "azurerm_service_plan"
}

func (r ServicePlanDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ServicePlanName,
		},

		"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
	}
}

func (r ServicePlanDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": location.SchemaComputed(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"app_service_environment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"os_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"per_site_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"number_of_workers": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"maximum_elastic_worker_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"reserved": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r ServicePlanDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var servicePlan ServicePlanModel
			if err := metadata.Decode(&servicePlan); err != nil {
				return err
			}

			id := parse.NewServicePlanID(subscriptionId, servicePlan.ResourceGroup, servicePlan.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.ServerfarmName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("%s not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			servicePlan.Location = location.NormalizeNilable(existing.Location)

			if sku := existing.Sku; sku != nil {
				if sku.Name != nil {
					servicePlan.Sku = *sku.Name
					if sku.Capacity != nil {
						servicePlan.NumberOfWorkers = int(*sku.Capacity)
					}
				}
			}

			if props := existing.AppServicePlanProperties; props != nil {
				if props.HyperV != nil && *props.HyperV {
					servicePlan.OSType = OSTypeWindowsContainer
				}
				if props.Reserved != nil && *props.Reserved {
					servicePlan.OSType = OSTypeLinux
				}

				if props.HostingEnvironmentProfile != nil && props.HostingEnvironmentProfile.ID != nil {
					servicePlan.AppServiceEnvironmentId = utils.NormalizeNilableString(props.HostingEnvironmentProfile.ID)
				}

				if v := props.PerSiteScaling; v != nil {
					servicePlan.PerSiteScaling = *v
				}

				if v := props.Reserved; v != nil {
					servicePlan.Reserved = *v
				}

				servicePlan.MaximumElasticWorkerCount = int(utils.NormaliseNilableInt32(props.MaximumElasticWorkerCount))
			}
			servicePlan.Tags = tags.ToTypedObject(existing.Tags)

			metadata.SetID(id)

			return metadata.Encode(&servicePlan)
		},
	}
}
