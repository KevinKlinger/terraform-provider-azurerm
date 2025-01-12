package validate

import (
	"regexp"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
)

func CdnEndpointCustomDomainName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+(-*[a-zA-Z0-9])*$`), "")
}
