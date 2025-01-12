package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
)

type ConsumptionBudgetResourceGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	BudgetName     string
}

func NewConsumptionBudgetResourceGroupID(subscriptionId, resourceGroup, budgetName string) ConsumptionBudgetResourceGroupId {
	return ConsumptionBudgetResourceGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		BudgetName:     budgetName,
	}
}

func (id ConsumptionBudgetResourceGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Budget Name %q", id.BudgetName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Consumption Budget Resource Group", segmentsStr)
}

func (id ConsumptionBudgetResourceGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Consumption/budgets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.BudgetName)
}

// ConsumptionBudgetResourceGroupID parses a ConsumptionBudgetResourceGroup ID into an ConsumptionBudgetResourceGroupId struct
func ConsumptionBudgetResourceGroupID(input string) (*ConsumptionBudgetResourceGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConsumptionBudgetResourceGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.BudgetName, err = id.PopSegment("budgets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
