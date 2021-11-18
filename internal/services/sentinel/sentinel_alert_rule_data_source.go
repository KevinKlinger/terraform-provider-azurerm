package sentinel

import (
	"fmt"
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	loganalyticsParse "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/loganalytics/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/sentinel/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func dataSourceSentinelAlertRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSentinelAlertRuleRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: loganalyticsValidate.LogAnalyticsWorkspaceID,
			},
		},
	}
}

func dataSourceSentinelAlertRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceID, err := loganalyticsParse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroup, workspaceID.WorkspaceName, name)

	resp, err := client.Get(ctx, workspaceID.ResourceGroup, OperationalInsightsResourceProvider, workspaceID.WorkspaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Sentinel Alert Rule %q was not found", id)
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return nil
}
