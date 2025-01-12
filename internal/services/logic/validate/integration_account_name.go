package validate

import (
	"regexp"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
)

func IntegrationAccountName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[\w-().]{1,80}$`), `Integration name can contain only letters, numbers, '_','-', '(', ')' or '.'`,
	)
}
