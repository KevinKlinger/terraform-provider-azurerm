package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
)

type LogicAppStandardId struct {
	SubscriptionId string
	ResourceGroup  string
	SiteName       string
}

func NewLogicAppStandardID(subscriptionId, resourceGroup, siteName string) LogicAppStandardId {
	return LogicAppStandardId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SiteName:       siteName,
	}
}

func (id LogicAppStandardId) String() string {
	segments := []string{
		fmt.Sprintf("Site Name %q", id.SiteName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Logic App Standard", segmentsStr)
}

func (id LogicAppStandardId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName)
}

// LogicAppStandardID parses a LogicAppStandard ID into an LogicAppStandardId struct
func LogicAppStandardID(input string) (*LogicAppStandardId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogicAppStandardId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
