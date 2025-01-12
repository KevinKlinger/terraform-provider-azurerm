package network

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func expandNetworkSubResourceID(input []interface{}) *[]network.SubResource {
	results := make([]network.SubResource, 0)
	for _, item := range input {
		id := item.(string)

		results = append(results, network.SubResource{
			ID: utils.String(id),
		})
	}
	return &results
}

func flattenNetworkSubResourceID(input *[]network.SubResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.ID != nil {
			results = append(results, *item.ID)
		}
	}

	return results
}
