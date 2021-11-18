package parse

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/relay/sdk/2017-04-01/hybridconnections"
)

func HybridConnectionID(input string) (*hybridconnections.HybridConnectionId, error) {
	return hybridconnections.ParseHybridConnectionID(input)
}
