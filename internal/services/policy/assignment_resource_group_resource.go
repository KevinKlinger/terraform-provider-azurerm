package policy

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/sdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/policy/validate"
	resourceValidate "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/resource/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = ResourceGroupAssignmentResource{}

type ResourceGroupAssignmentResource struct {
	base assignmentBaseResource
}

func (r ResourceGroupAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"resource_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourceValidate.ResourceGroupID,
		},
	}
	return r.base.arguments(schema)
}

func (r ResourceGroupAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r ResourceGroupAssignmentResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "resource_group_id")
}

func (r ResourceGroupAssignmentResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r ResourceGroupAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ResourceGroupAssignmentID
}

func (r ResourceGroupAssignmentResource) ModelObject() interface{} {
	return nil
}

func (r ResourceGroupAssignmentResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("resource_group_id")
}

func (r ResourceGroupAssignmentResource) ResourceType() string {
	return "azurerm_resource_group_policy_assignment"
}

func (r ResourceGroupAssignmentResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}
