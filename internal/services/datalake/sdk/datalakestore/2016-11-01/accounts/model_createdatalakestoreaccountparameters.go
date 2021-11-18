package accounts

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/identity"
)

type CreateDataLakeStoreAccountParameters struct {
	Identity   *identity.SystemAssignedIdentity      `json:"identity,omitempty"`
	Location   string                                `json:"location"`
	Properties *CreateDataLakeStoreAccountProperties `json:"properties,omitempty"`
	Tags       *map[string]string                    `json:"tags,omitempty"`
}
