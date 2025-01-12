package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/logic/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/logic/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
)

func resourceLogicAppTriggerHttpRequest() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppTriggerHttpRequestCreateUpdate,
		Read:   resourceLogicAppTriggerHttpRequestRead,
		Update: resourceLogicAppTriggerHttpRequestCreateUpdate,
		Delete: resourceLogicAppTriggerHttpRequestDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.TriggerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
			relativePath := diff.Get("relative_path").(string)
			if relativePath != "" {
				method := diff.Get("method").(string)
				if method == "" {
					return fmt.Errorf("`method` must be specified when `relative_path` is set.")
				}
			}

			return nil
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"logic_app_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"schema": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"method": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					http.MethodDelete,
					http.MethodGet,
					http.MethodPatch,
					http.MethodPost,
					http.MethodPut,
				}, false),
			},

			"relative_path": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.TriggerHttpRequestRelativePath,
			},

			"callback_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLogicAppTriggerHttpRequestCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	schemaRaw := d.Get("schema").(string)
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(schemaRaw), &schema); err != nil {
		return fmt.Errorf("unmarshalling JSON from Schema: %+v", err)
	}

	inputs := map[string]interface{}{
		"schema": schema,
	}

	if v, ok := d.GetOk("method"); ok {
		inputs["method"] = v.(string)
	}

	if v, ok := d.GetOk("relative_path"); ok {
		inputs["relativePath"] = v.(string)
	}

	trigger := map[string]interface{}{
		"inputs": inputs,
		"kind":   "Http",
		"type":   "Request",
	}

	workflowId, err := parse.WorkflowID(d.Get("logic_app_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewTriggerID(workflowId.SubscriptionId, workflowId.ResourceGroup, workflowId.Name, d.Get("name").(string))

	if err := resourceLogicAppTriggerUpdate(d, meta, *workflowId, id, trigger, "azurerm_logic_app_trigger_http_request"); err != nil {
		return err
	}

	return resourceLogicAppTriggerHttpRequestRead(d, meta)
}

func resourceLogicAppTriggerHttpRequestRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	t, app, url, err := retrieveLogicAppHttpTrigger(d, meta, id.ResourceGroup, id.WorkflowName, id.Name)
	if err != nil {
		return err
	}

	if t == nil {
		log.Printf("[DEBUG] Logic App %q (Resource Group %q) does not contain Trigger %q - removing from state", id.WorkflowName, id.ResourceGroup, id.Name)
		d.SetId("")
		return nil
	}

	trigger := *t

	d.Set("name", id.Name)
	d.Set("logic_app_id", app.ID)
	d.Set("callback_url", url)

	v := trigger["inputs"]
	if v == nil {
		return fmt.Errorf("`inputs` was nil for HTTP Trigger %s", id)
	}

	inputs, ok := v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("parsing `inputs` for HTTP Trigger %s", id)
	}

	if method := inputs["method"]; method != nil {
		d.Set("method", method.(string))
	}

	if relativePath := inputs["relativePath"]; relativePath != nil {
		d.Set("relative_path", relativePath.(string))
	}

	if schemaRaw := inputs["schema"]; schemaRaw != nil {
		schema, err := json.Marshal(schemaRaw)
		if err != nil {
			return fmt.Errorf("serializing the Schema to JSON: %+v", err)
		}

		d.Set("schema", string(schema))
	}

	return nil
}

func resourceLogicAppTriggerHttpRequestDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	err = resourceLogicAppTriggerRemove(d, meta, id.ResourceGroup, id.WorkflowName, id.Name)
	if err != nil {
		return fmt.Errorf("removing Trigger %s: %+v", id, err)
	}

	return nil
}
