package validate

import (
	"fmt"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/relay/sdk/2017-04-01/namespaces"
)

func NamespaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := namespaces.ParseNamespaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
