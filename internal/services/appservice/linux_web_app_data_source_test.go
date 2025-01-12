package appservice_test

import (
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
)

type LinuxWebAppDataSource struct{}

func TestAccLinuxWebAppDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_linux_web_app", "test")
	d := LinuxWebAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
			),
		},
	})
}

func (LinuxWebAppDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data azurerm_linux_web_app test {
  name                = azurerm_linux_web_app.test.name
  resource_group_name = azurerm_linux_web_app.test.resource_group_name
}

`, LinuxWebAppResource{}.basic(data)) // TODO - Update this to l.complete when that config is written
}
