package parse

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/msi/sdk/2018-11-30/managedidentity"
)

// these are here primarily to enable migration over

func UserAssignedIdentityID(input string) (*managedidentity.UserAssignedIdentitiesId, error) {
	return managedidentity.ParseUserAssignedIdentitiesID(input)
}

func UserAssignedIdentityIDInsensitively(input string) (*managedidentity.UserAssignedIdentitiesId, error) {
	return managedidentity.ParseUserAssignedIdentitiesIDInsensitively(input)
}
