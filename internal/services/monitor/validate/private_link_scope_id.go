package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/monitor/parse"
)

func PrivateLinkScopeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := parse.PrivateLinkScopeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
