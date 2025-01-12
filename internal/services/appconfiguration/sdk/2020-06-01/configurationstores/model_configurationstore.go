package configurationstores

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/identity"
)

type ConfigurationStore struct {
	Id         *string                                 `json:"id,omitempty"`
	Identity   *identity.SystemUserAssignedIdentityMap `json:"identity,omitempty"`
	Location   string                                  `json:"location"`
	Name       *string                                 `json:"name,omitempty"`
	Properties *ConfigurationStoreProperties           `json:"properties,omitempty"`
	Sku        Sku                                     `json:"sku"`
	Tags       *map[string]string                      `json:"tags,omitempty"`
	Type       *string                                 `json:"type,omitempty"`
}
