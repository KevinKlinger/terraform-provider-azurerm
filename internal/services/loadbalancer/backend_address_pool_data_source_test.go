package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
)

func TestAccDataSourceBackendAddressPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSourceBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (r LoadBalancerBackendAddressPool) dataSourceBasic(data acceptance.TestData) string {
	resource := r.basicSkuBasic(data)
	return fmt.Sprintf(`
%s

data "azurerm_lb_backend_address_pool" "test" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, resource)
}
