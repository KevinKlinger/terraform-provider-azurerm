package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/monitor/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ActionGroupUpgradeV0ToV1{}

type ActionGroupUpgradeV0ToV1 struct{}

func (ActionGroupUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return actionGroupSchemaForV0AndV1()
}

func (ActionGroupUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/actionGroups/actionGroup1
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/actionGroups/actionGroup1
		oldId, err := azure.ParseAzureResourceID(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		groupName := ""
		for key, value := range oldId.Path {
			if strings.EqualFold(key, "actionGroups") {
				groupName = value
				break
			}
		}

		if groupName == "" {
			return rawState, fmt.Errorf("couldn't find the `actionGroups` segment in the old resource id %q", oldId)
		}

		newId := parse.NewActionGroupID(oldId.SubscriptionID, oldId.ResourceGroup, groupName)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId.ID())

		rawState["id"] = newId.ID()

		return rawState, nil
	}
}

func actionGroupSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"short_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"email_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"email_address": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"use_common_alert_schema": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"itsm_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"workspace_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"connection_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"ticket_configuration": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"region": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"azure_app_push_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"email_address": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"sms_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"country_code": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"phone_number": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"webhook_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"service_uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"use_common_alert_schema": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"aad_auth": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"object_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"identifier_uri": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
								},

								"tenant_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"automation_runbook_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"automation_account_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"runbook_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"webhook_resource_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"is_global_runbook": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"service_uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"use_common_alert_schema": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"voice_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"country_code": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"phone_number": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"logic_app_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"resource_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"callback_url": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"use_common_alert_schema": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"azure_function_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"function_app_resource_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"function_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"http_trigger_url": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"use_common_alert_schema": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"arm_role_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"role_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"use_common_alert_schema": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"tags": tags.Schema(),
	}
}
