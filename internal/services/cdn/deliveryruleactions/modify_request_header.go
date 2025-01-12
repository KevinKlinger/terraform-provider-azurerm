package deliveryruleactions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func ModifyRequestHeader() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"action": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.HeaderActionAppend),
					string(cdn.HeaderActionDelete),
					string(cdn.HeaderActionOverwrite),
				}, false),
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"value": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func ExpandArmCdnEndpointActionModifyRequestHeader(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		requestHeaderAction := cdn.DeliveryRuleRequestHeaderAction{
			Name: cdn.NameBasicDeliveryRuleActionNameModifyRequestHeader,
			Parameters: &cdn.HeaderActionParameters{
				OdataType:    utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleHeaderActionParameters"),
				HeaderAction: cdn.HeaderAction(item["action"].(string)),
				HeaderName:   utils.String(item["name"].(string)),
			},
		}

		if value := item["value"].(string); value != "" {
			requestHeaderAction.Parameters.Value = utils.String(value)
		}

		output = append(output, requestHeaderAction)
	}

	return &output, nil
}

func FlattenArmCdnEndpointActionModifyRequestHeader(input cdn.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleRequestHeaderAction()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request header action!")
	}

	headerAction := ""
	headerName := ""
	value := ""
	if params := action.Parameters; params != nil {
		headerAction = string(params.HeaderAction)

		if params.HeaderName != nil {
			headerName = *params.HeaderName
		}

		if params.Value != nil {
			value = *params.Value
		}
	}

	return &map[string]interface{}{
		"action": headerAction,
		"name":   headerName,
		"value":  value,
	}, nil
}
