package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/mysql/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type MySQLFlexibleServerFirewallRuleResource struct {
}

func TestAccMySQLFlexibleServerFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_firewall_rule", "test")
	r := MySQLFlexibleServerFirewallRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMySQLFlexibleServerFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_firewall_rule", "test")
	r := MySQLFlexibleServerFirewallRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t MySQLFlexibleServerFirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FlexibleServerFirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MySQL.FlexibleServerFirewallRulesClient.Get(ctx, id.ResourceGroup, id.FlexibleServerName, id.FirewallRuleName)
	if err != nil {
		return nil, err
	}

	return utils.Bool(resp.Name != nil), nil
}

func (MySQLFlexibleServerFirewallRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1s"
}

resource "azurerm_mysql_flexible_server_firewall_rule" "test" {
  name                = "acctestfwrule-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_flexible_server.test.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r MySQLFlexibleServerFirewallRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server_firewall_rule" "import" {
  name                = azurerm_mysql_flexible_server_firewall_rule.test.name
  resource_group_name = azurerm_mysql_flexible_server_firewall_rule.test.resource_group_name
  server_name         = azurerm_mysql_flexible_server_firewall_rule.test.server_name
  start_ip_address    = azurerm_mysql_flexible_server_firewall_rule.test.start_ip_address
  end_ip_address      = azurerm_mysql_flexible_server_firewall_rule.test.end_ip_address
}
`, r.basic(data))
}
