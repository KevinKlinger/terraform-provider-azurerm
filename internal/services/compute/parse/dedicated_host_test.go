package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/resourceid"
)

var _ resourceid.Formatter = DedicatedHostId{}

func TestDedicatedHostIDFormatter(t *testing.T) {
	actual := NewDedicatedHostID("12345678-1234-9876-4563-123456789012", "resGroup1", "hostGroup1", "host1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/hostGroup1/hosts/host1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestDedicatedHostID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *DedicatedHostId
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
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing HostGroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/",
			Error: true,
		},

		{
			// missing value for HostGroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/",
			Error: true,
		},

		{
			// missing HostName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/hostGroup1/",
			Error: true,
		},

		{
			// missing value for HostName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/hostGroup1/hosts/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/hostGroup1/hosts/host1",
			Expected: &DedicatedHostId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resGroup1",
				HostGroupName:  "hostGroup1",
				HostName:       "host1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.COMPUTE/HOSTGROUPS/HOSTGROUP1/HOSTS/HOST1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := DedicatedHostID(v.Input)
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
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.HostGroupName != v.Expected.HostGroupName {
			t.Fatalf("Expected %q but got %q for HostGroupName", v.Expected.HostGroupName, actual.HostGroupName)
		}
		if actual.HostName != v.Expected.HostName {
			t.Fatalf("Expected %q but got %q for HostName", v.Expected.HostName, actual.HostName)
		}
	}
}
