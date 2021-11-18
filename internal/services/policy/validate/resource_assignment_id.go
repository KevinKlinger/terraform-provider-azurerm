package validate

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	managementGroupValidate "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/managementgroup/validate"
	resourceValidate "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/resource/validate"
	subscriptionValidate "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/subscription/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
)

func ResourceAssignmentId() pluginsdk.SchemaValidateFunc {
	return validation.All(
		validation.None(
			map[string]func(interface{}, string) ([]string, []error){
				"Management Group ID": managementGroupValidate.ManagementGroupID,
				"Resource Group ID":   resourceValidate.ResourceGroupID,
				"Subscription ID":     subscriptionValidate.SubscriptionID,
			},
		),
		azure.ValidateResourceID,
	)
}
