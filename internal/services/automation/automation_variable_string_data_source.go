package automation

import (
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
)

func dataSourceAutomationVariableString() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAutomationVariableStringRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: datasourceAutomationVariableCommonSchema(pluginsdk.TypeString),
	}
}

func dataSourceAutomationVariableStringRead(d *pluginsdk.ResourceData, meta interface{}) error {
	return dataSourceAutomationVariableRead(d, meta, "String")
}
