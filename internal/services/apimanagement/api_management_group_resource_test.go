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

type ApiManagementGroupResource struct {
}

func TestAccApiManagementGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group", "test")
	r := ApiManagementGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Group"),
				check.That(data.ResourceName).Key("type").HasValue("custom"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group", "test")
	r := ApiManagementGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Group"),
				check.That(data.ResourceName).Key("type").HasValue("custom"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApiManagementGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group", "test")
	r := ApiManagementGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "Test Group", "A test description."),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Group"),
				check.That(data.ResourceName).Key("description").HasValue("A test description."),
				check.That(data.ResourceName).Key("type").HasValue("external"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementGroup_descriptionDisplayNameUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group", "test")
	r := ApiManagementGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "Original Group", "The original description."),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Original Group"),
				check.That(data.ResourceName).Key("description").HasValue("The original description."),
				check.That(data.ResourceName).Key("type").HasValue("external"),
			),
		},
		{
			Config: r.complete(data, "Modified Group", "A modified description."),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Modified Group"),
				check.That(data.ResourceName).Key("description").HasValue("A modified description."),
				check.That(data.ResourceName).Key("type").HasValue("external"),
			),
		},
		{
			Config: r.complete(data, "Original Group", "The original description."),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Original Group"),
				check.That(data.ResourceName).Key("description").HasValue("The original description."),
				check.That(data.ResourceName).Key("type").HasValue("external"),
			),
		},
	})
}

func (ApiManagementGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.GroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.GroupClient.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementGroupResource) basic(data acceptance.TestData) string {
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

  sku_name = "Developer_1"
}

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "Test Group"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_group" "import" {
  name                = azurerm_api_management_group.test.name
  resource_group_name = azurerm_api_management_group.test.resource_group_name
  api_management_name = azurerm_api_management_group.test.api_management_name
  display_name        = azurerm_api_management_group.test.display_name
}
`, r.basic(data))
}

func (ApiManagementGroupResource) complete(data acceptance.TestData, displayName, description string) string {
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

  sku_name = "Developer_1"
}

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "%s"
  description         = "%s"
  type                = "external"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, displayName, description)
}
