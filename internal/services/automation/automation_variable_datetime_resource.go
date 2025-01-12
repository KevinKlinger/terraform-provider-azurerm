package automation

import (
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/automation/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
)

func resourceAutomationVariableDateTime() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationVariableDateTimeCreateUpdate,
		Read:   resourceAutomationVariableDateTimeRead,
		Update: resourceAutomationVariableDateTimeCreateUpdate,
		Delete: resourceAutomationVariableDateTimeDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VariableID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceAutomationVariableCommonSchema(pluginsdk.TypeString, validation.IsRFC3339Time),
	}
}

func resourceAutomationVariableDateTimeCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Datetime")
}

func resourceAutomationVariableDateTimeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Datetime")
}

func resourceAutomationVariableDateTimeDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Datetime")
}
