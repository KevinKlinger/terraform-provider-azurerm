package identity

import (
	"fmt"

	msivalidate "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/msi/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

var _ Identity = SystemAssignedUserAssigned{}

type SystemAssignedUserAssigned struct{}

func (s SystemAssignedUserAssigned) Expand(input []interface{}) (*ExpandedConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return &ExpandedConfig{
			Type: none,
		}, nil
	}

	v := input[0].(map[string]interface{})

	config := &ExpandedConfig{
		Type: Type(v["type"].(string)),
	}

	identityIds := v["identity_ids"].(*pluginsdk.Set).List()

	if len(identityIds) != 0 {
		if config.Type != userAssigned && config.Type != systemAssignedUserAssigned {
			return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`")
		}
		config.UserAssignedIdentityIds = *utils.ExpandStringSlice(identityIds)
	}

	return config, nil
}

func (s SystemAssignedUserAssigned) Flatten(input *ExpandedConfig) []interface{} {
	if input == nil || input.Type == none {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"type":         input.Type,
			"identity_ids": utils.FlattenStringSlice(&input.UserAssignedIdentityIds),
			"principal_id": input.PrincipalId,
			"tenant_id":    input.TenantId,
		},
	}
}

func (s SystemAssignedUserAssigned) Schema() *pluginsdk.Schema {
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
						string(userAssigned),
						string(systemAssigned),
						string(systemAssignedUserAssigned),
					}, false),
				},
				"identity_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: msivalidate.UserAssignedIdentityID,
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

func (s SystemAssignedUserAssigned) SchemaDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"identity_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
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
