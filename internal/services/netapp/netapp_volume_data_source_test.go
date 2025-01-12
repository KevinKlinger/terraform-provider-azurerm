package netapp_test

import (
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
)

type NetAppVolumeDataSource struct {
}

func TestAccDataSourceNetAppVolume_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_volume", "test")
	r := NetAppVolumeDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("volume_path").Exists(),
				check.That(data.ResourceName).Key("service_level").Exists(),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				check.That(data.ResourceName).Key("storage_quota_in_gb").Exists(),
				check.That(data.ResourceName).Key("protocols.0").Exists(),
				check.That(data.ResourceName).Key("mount_ip_addresses.#").HasValue("1"),
			),
		},
	})
}

func (NetAppVolumeDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_volume" "test" {
  resource_group_name = azurerm_netapp_volume.test.resource_group_name
  account_name        = azurerm_netapp_volume.test.account_name
  pool_name           = azurerm_netapp_volume.test.pool_name
  name                = azurerm_netapp_volume.test.name
}
`, NetAppVolumeResource{}.basic(data))
}
