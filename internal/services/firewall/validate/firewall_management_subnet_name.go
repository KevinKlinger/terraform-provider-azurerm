package validate

import (
	"fmt"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
)

func FirewallManagementSubnetName(v interface{}, k string) (warnings []string, errors []error) {
	parsed, err := azure.ParseAzureResourceID(v.(string))
	if err != nil {
		errors = append(errors, fmt.Errorf("parsing Azure Resource ID %q", v.(string)))
		return warnings, errors
	}
	subnetName := parsed.Path["subnets"]
	if subnetName != "AzureFirewallManagementSubnet" {
		errors = append(errors, fmt.Errorf("The name of the management subnet for %q must be exactly 'AzureFirewallManagementSubnet' to be used for the Azure Firewall resource", k))
	}

	return warnings, errors
}
