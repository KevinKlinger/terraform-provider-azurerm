package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/securitycenter/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type SecurityCenterSettingResource struct {
}

func TestAccSecurityCenterSetting_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_setting", "test")
	r := SecurityCenterSettingResource{}

	// lintignore:AT001
	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.cfg("MCAS", true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("setting_name").HasValue("MCAS"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cfg("MCAS", false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("setting_name").HasValue("MCAS"),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cfg("WDATP", true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("setting_name").HasValue("WDATP"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cfg("WDATP", false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("setting_name").HasValue("WDATP"),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (SecurityCenterSettingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SecurityCenterSettingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.SettingClient.Get(ctx, id.SettingName)
	if err != nil {
		return nil, fmt.Errorf("reading Security Center Setting (%s): %+v", id.SettingName, err)
	}

	return utils.Bool(resp.Value != nil), nil
}

func (SecurityCenterSettingResource) cfg(settingName string, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_setting" "test" {
  setting_name = "%s"
  enabled      = "%t"
}
`, settingName, enabled)
}
