package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
)

type DatabaseId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
	Name           string
}

func NewDatabaseID(subscriptionId, resourceGroup, clusterName, name string) DatabaseId {
	return DatabaseId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ClusterName:    clusterName,
		Name:           name,
	}
}

func (id DatabaseId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Database", segmentsStr)
}

func (id DatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/Clusters/%s/Databases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.Name)
}

// DatabaseID parses a Database ID into an DatabaseId struct
func DatabaseID(input string) (*DatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DatabaseId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ClusterName, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("Databases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
