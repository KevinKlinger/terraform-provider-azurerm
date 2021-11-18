package validate

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/analysisservices/sdk/2017-08-01/servers"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
)

func QueryPoolConnectionMode() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(servers.ConnectionModeAll),
		string(servers.ConnectionModeReadOnly),
	}, true)
}
