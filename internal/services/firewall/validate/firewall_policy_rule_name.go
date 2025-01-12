package validate

import (
	"regexp"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
)

func FirewallPolicyRuleName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^\W_][\w-.]*[\w]$`),
		"The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.")
}
