package configurationstores

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/identity"
)

type ConfigurationStoreUpdateParameters struct {
	Identity   *identity.SystemUserAssignedIdentityMap       `json:"identity,omitempty"`
	Properties *ConfigurationStorePropertiesUpdateParameters `json:"properties,omitempty"`
	Sku        *Sku                                          `json:"sku,omitempty"`
	Tags       *map[string]string                            `json:"tags,omitempty"`
}
