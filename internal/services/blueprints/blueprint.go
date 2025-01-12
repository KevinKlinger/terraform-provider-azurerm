package blueprints

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/blueprint/mgmt/2018-11-01-preview/blueprint"
	msiparse "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/msi/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/msi/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/suppress"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func ManagedIdentitySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						// ManagedServiceIdentityTypeNone is not valid; a valid and privileged Identity is required for the service to apply the changes.
						// SystemAssigned type not currently supported - The Portal performs significant activity in temporary escalation of permissions to Owner on the target scope
						// Such activity in the Provider would be brittle
						// string(blueprint.ManagedServiceIdentityTypeSystemAssigned),
						string(blueprint.ManagedServiceIdentityTypeUserAssigned),
					}, true),
					// The first character of value returned by the service is always in lower case - bug?
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"identity_ids": {
					// The API only seems to care about the "key" portion of this struct, which is the ResourceID of the Identity
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validate.UserAssignedIdentityID,
					},
				},

				"principal_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"tenant_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func blueprintAssignmentCreateStateRefreshFunc(ctx context.Context, client *blueprint.AssignmentsClient, scope, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, scope, name)
		if err != nil {
			return nil, "", fmt.Errorf("unable to retrieve Blueprint Assignment %q (Scope %q): %+v", name, scope, err)
		}
		if resp.ProvisioningState == blueprint.Failed {
			return resp, string(resp.ProvisioningState), fmt.Errorf("Blueprint Assignment provisioning entered a Failed state.")
		}

		return resp, string(resp.ProvisioningState), nil
	}
}

func blueprintAssignmentDeleteStateRefreshFunc(ctx context.Context, client *blueprint.AssignmentsClient, scope, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, scope, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "NotFound", nil
			} else {
				return nil, "", fmt.Errorf("unable to retrieve Blueprint Assignment %q (Scope %q): %+v", name, scope, err)
			}
		}

		return resp, string(resp.ProvisioningState), nil
	}
}

func normalizeAssignmentParameterValuesJSON(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}

	var values map[string]*blueprint.ParameterValue
	if err := json.Unmarshal([]byte(jsonString.(string)), &values); err != nil {
		return fmt.Sprintf("unable to parse JSON: %+v", err)
	}

	b, _ := json.Marshal(values)
	return string(b)
}

func normalizeAssignmentResourceGroupValuesJSON(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}

	var values map[string]*blueprint.ResourceGroupValue
	if err := json.Unmarshal([]byte(jsonString.(string)), &values); err != nil {
		return fmt.Sprintf("unable to parse JSON: %+v", err)
	}

	b, _ := json.Marshal(values)
	return string(b)
}

func expandArmBlueprintAssignmentParameters(input string) map[string]*blueprint.ParameterValue {
	var result map[string]*blueprint.ParameterValue
	// the string has been validated by the schema, therefore the error is ignored here, since it will never happen.
	_ = json.Unmarshal([]byte(input), &result)
	return result
}

func expandArmBlueprintAssignmentResourceGroups(input string) map[string]*blueprint.ResourceGroupValue {
	var result map[string]*blueprint.ResourceGroupValue
	// the string has been validated by the schema, therefore the error is ignored here, since it will never happen.
	_ = json.Unmarshal([]byte(input), &result)
	return result
}

func expandArmBlueprintAssignmentIdentity(input []interface{}) (*blueprint.ManagedServiceIdentity, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, fmt.Errorf("Managed Service Identity was empty")
	}

	raw := input[0].(map[string]interface{})

	identity := blueprint.ManagedServiceIdentity{
		Type: blueprint.ManagedServiceIdentityType(raw["type"].(string)),
	}

	identityIdsRaw := raw["identity_ids"].([]interface{})
	identityIds := make(map[string]*blueprint.UserAssignedIdentity)
	for _, v := range identityIdsRaw {
		identityIds[v.(string)] = &blueprint.UserAssignedIdentity{}
	}
	identity.UserAssignedIdentities = identityIds

	return &identity, nil
}

func flattenArmBlueprintAssignmentIdentity(input *blueprint.ManagedServiceIdentity) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	identityIds := make([]string, 0)
	if input.UserAssignedIdentities != nil {
		for key := range input.UserAssignedIdentities {
			parsedId, err := msiparse.UserAssignedIdentityID(key)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
	}

	principalId := ""
	if input.PrincipalID != nil {
		principalId = *input.PrincipalID
	}

	tenantId := ""
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"identity_ids": identityIds,
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}, nil
}

func flattenArmBlueprintAssignmentParameters(input map[string]*blueprint.ParameterValue) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func flattenArmBlueprintAssignmentResourceGroups(input map[string]*blueprint.ResourceGroupValue) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
