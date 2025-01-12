package validate

import (
	"fmt"
	"strings"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
)

func IntegrationAccountSchemaFileName() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", k))
			return
		}

		if !strings.HasSuffix(v, ".xsd") {
			errors = append(errors, fmt.Errorf("%q ends with `.xsd`.", k))
			return
		}

		return
	}
}
