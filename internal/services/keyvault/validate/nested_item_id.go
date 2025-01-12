package validate

import (
	"fmt"

	keyVaultParse "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/keyvault/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
)

func NestedItemId(i interface{}, k string) (warnings []string, errors []error) {
	if warnings, errors = validation.StringIsNotEmpty(i, k); len(errors) > 0 {
		return warnings, errors
	}

	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("Expected %s to be a string!", k))
		return warnings, errors
	}

	if _, err := keyVaultParse.ParseNestedItemID(v); err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %s", v, err))
		return warnings, errors
	}

	return warnings, errors
}

func VersionlessNestedItemId(i interface{}, k string) (warnings []string, errors []error) {
	if warnings, errors = validation.StringIsNotEmpty(i, k); len(errors) > 0 {
		return warnings, errors
	}

	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("Expected %s to be a string!", k))
		return warnings, errors
	}

	id, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %s", v, err))
		return warnings, errors
	}

	if id.Version != "" {
		errors = append(errors, fmt.Errorf("expected %s to not have a version", k))
	}

	return warnings, errors
}

func NestedItemIdWithOptionalVersion(i interface{}, k string) (warnings []string, errors []error) {
	if warnings, errors = validation.StringIsNotEmpty(i, k); len(errors) > 0 {
		return warnings, errors
	}

	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("Expected %s to be a string!", k))
		return warnings, errors
	}

	if _, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(v); err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %s", v, err))
		return warnings, errors
	}

	return warnings, errors
}
