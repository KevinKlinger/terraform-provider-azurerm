package validate

import (
	"fmt"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/signalr/sdk/2020-05-01/signalr"
)

func ServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := signalr.ParseSignalRID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
