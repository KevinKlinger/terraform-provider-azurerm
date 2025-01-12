package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/resourceid"
)

var _ resourceid.Formatter = ConsumptionBudgetSubscriptionId{}

func TestConsumptionBudgetSubscriptionIDFormatter(t *testing.T) {
	actual := NewConsumptionBudgetSubscriptionID("12345678-1234-9876-4563-123456789012", "budget1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Consumption/budgets/budget1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestConsumptionBudgetSubscriptionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ConsumptionBudgetSubscriptionId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing BudgetName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Consumption/",
			Error: true,
		},

		{
			// missing value for BudgetName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Consumption/budgets/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Consumption/budgets/budget1",
			Expected: &ConsumptionBudgetSubscriptionId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				BudgetName:     "budget1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.CONSUMPTION/BUDGETS/BUDGET1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ConsumptionBudgetSubscriptionID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.BudgetName != v.Expected.BudgetName {
			t.Fatalf("Expected %q but got %q for BudgetName", v.Expected.BudgetName, actual.BudgetName)
		}
	}
}
