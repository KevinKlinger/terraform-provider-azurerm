package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type SecurityCenterAutoProvisionResource struct {
}

func TestAccSecurityCenterAutoProvision_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_auto_provisioning", "test")
	r := SecurityCenterAutoProvisionResource{}

	// lintignore:AT001
	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.setting("On"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_provision").HasValue("On"),
			),
		},
		data.ImportStep(),
		{
			Config: r.setting("Off"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_provision").HasValue("Off"),
			),
		},
		data.ImportStep(),
	})
}

func (SecurityCenterAutoProvisionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	securityCenterAutoProvisioningName := "default"

	resp, err := clients.SecurityCenter.AutoProvisioningClient.Get(ctx, securityCenterAutoProvisioningName)
	if err != nil {
		return nil, fmt.Errorf("reading Security Center auto provision (%s): %+v", securityCenterAutoProvisioningName, err)
	}

	return utils.Bool(resp.AutoProvisioningSettingProperties != nil), nil
}

func (SecurityCenterAutoProvisionResource) setting(setting string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_auto_provisioning" "test" {
  auto_provision = "%s"
}
`, setting)
}
