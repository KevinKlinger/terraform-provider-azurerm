package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/apimanagement/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type ApiManagementPropertyResource struct {
}

func TestAccApiManagementProperty_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_property", "test")
	r := ApiManagementPropertyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("TestProperty%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("value").HasValue("Test Value"),
				check.That(data.ResourceName).Key("tags.0").HasValue("tag1"),
				check.That(data.ResourceName).Key("tags.1").HasValue("tag2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementProperty_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_property", "test")
	r := ApiManagementPropertyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("TestProperty%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("value").HasValue("Test Value"),
				check.That(data.ResourceName).Key("tags.0").HasValue("tag1"),
				check.That(data.ResourceName).Key("tags.1").HasValue("tag2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("TestProperty2%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("secret").HasValue("true"),
				check.That(data.ResourceName).Key("tags.0").HasValue("tag3"),
				check.That(data.ResourceName).Key("tags.1").HasValue("tag4"),
			),
		},
		data.ImportStep("value"),
	})
}

func (ApiManagementPropertyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PropertyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.NamedValueClient.Get(ctx, id.ResourceGroup, id.ServiceName, id.NamedValueName)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementPropertyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}

resource "azurerm_api_management_property" "test" {
  name                = "acctestAMProperty-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestProperty%d"
  value               = "Test Value"
  tags                = ["tag1", "tag2"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (ApiManagementPropertyResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}

resource "azurerm_api_management_property" "test" {
  name                = "acctestAMProperty-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestProperty2%d"
  value               = "Test Value2"
  secret              = true
  tags                = ["tag3", "tag4"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
