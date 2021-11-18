package parse

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/relay/sdk/2017-04-01/namespaces"
)

func NamespaceID(input string) (*namespaces.NamespaceId, error) {
	return namespaces.ParseNamespaceID(input)
}
