package automation

import (
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/automation/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
)

func resourceAutomationVariableInt() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationVariableIntCreateUpdate,
		Read:   resourceAutomationVariableIntRead,
		Update: resourceAutomationVariableIntCreateUpdate,
		Delete: resourceAutomationVariableIntDelete,

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

		Schema: resourceAutomationVariableCommonSchema(pluginsdk.TypeInt, nil),
	}
}

func resourceAutomationVariableIntCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Int")
}

func resourceAutomationVariableIntRead(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Int")
}

func resourceAutomationVariableIntDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Int")
}
