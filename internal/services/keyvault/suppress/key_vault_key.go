package suppress

import (
	keyVaultParse "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/keyvault/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
)

func DiffSuppressIgnoreKeyVaultKeyVersion(k, old, new string, d *pluginsdk.ResourceData) bool {
	oldKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(old)
	if err != nil {
		return false
	}
	newKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(new)
	if err != nil {
		return false
	}

	return (oldKey.KeyVaultBaseUrl == newKey.KeyVaultBaseUrl) && (oldKey.Name == newKey.Name)
}
