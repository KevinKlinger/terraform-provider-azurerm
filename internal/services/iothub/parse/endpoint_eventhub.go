package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
)

type EndpointEventhubId struct {
	SubscriptionId string
	ResourceGroup  string
	IotHubName     string
	EndpointName   string
}

func NewEndpointEventhubID(subscriptionId, resourceGroup, iotHubName, endpointName string) EndpointEventhubId {
	return EndpointEventhubId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		IotHubName:     iotHubName,
		EndpointName:   endpointName,
	}
}

func (id EndpointEventhubId) String() string {
	segments := []string{
		fmt.Sprintf("Endpoint Name %q", id.EndpointName),
		fmt.Sprintf("Iot Hub Name %q", id.IotHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Endpoint Eventhub", segmentsStr)
}

func (id EndpointEventhubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/IotHubs/%s/Endpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotHubName, id.EndpointName)
}

// EndpointEventhubID parses a EndpointEventhub ID into an EndpointEventhubId struct
func EndpointEventhubID(input string) (*EndpointEventhubId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EndpointEventhubId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IotHubName, err = id.PopSegment("IotHubs"); err != nil {
		return nil, err
	}
	if resourceId.EndpointName, err = id.PopSegment("Endpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
