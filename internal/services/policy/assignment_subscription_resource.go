package policy

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/sdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/policy/validate"
	subscriptionValidate "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/subscription/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = SubscriptionAssignmentResource{}

type SubscriptionAssignmentResource struct {
	base assignmentBaseResource
}

func (r SubscriptionAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: subscriptionValidate.SubscriptionID,
		},
	}
	return r.base.arguments(schema)
}

func (r SubscriptionAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r SubscriptionAssignmentResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "subscription_id")
}

func (r SubscriptionAssignmentResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r SubscriptionAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SubscriptionAssignmentID
}

func (r SubscriptionAssignmentResource) ModelObject() interface{} {
	return nil
}

func (r SubscriptionAssignmentResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("subscription_id")
}

func (r SubscriptionAssignmentResource) ResourceType() string {
	return "azurerm_subscription_policy_assignment"
}

func (r SubscriptionAssignmentResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}
