package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/datafactory/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type IntegrationRuntimeAzureResource struct {
}

func TestAccDataFactoryIntegrationRuntimeAzure_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure", "test")
	r := IntegrationRuntimeAzureResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("compute_type").HasValue("General"),
				check.That(data.ResourceName).Key("core_count").HasValue("8"),
				check.That(data.ResourceName).Key("time_to_live_min").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeAzure_computeType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure", "test")
	r := IntegrationRuntimeAzureResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.computeType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("compute_type").HasValue("ComputeOptimized"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeAzure_coreCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure", "test")
	r := IntegrationRuntimeAzureResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.coreCount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("core_count").HasValue("16"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeAzure_timeToLive(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure", "test")
	r := IntegrationRuntimeAzureResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.timeToLive(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("time_to_live_min").HasValue("10"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeAzure_virtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure", "test")
	r := IntegrationRuntimeAzureResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetwork(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeAzure_cleanup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure", "test")
	r := IntegrationRuntimeAzureResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cleanup(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.cleanup(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.cleanup(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (IntegrationRuntimeAzureResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_azure" "test" {
  name                = "azure-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IntegrationRuntimeAzureResource) computeType(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_azure" "test" {
  name                = "azure-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  compute_type        = "ComputeOptimized"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IntegrationRuntimeAzureResource) coreCount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_azure" "test" {
  name                = "azure-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  core_count          = 16
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IntegrationRuntimeAzureResource) timeToLive(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_azure" "test" {
  name                = "azure-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  time_to_live_min    = 10
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IntegrationRuntimeAzureResource) cleanup(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_azure" "test" {
  name                = "azure-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cleanup_enabled     = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enabled)
}

func (IntegrationRuntimeAzureResource) virtualNetwork(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                            = "acctestdf%d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  managed_virtual_network_enabled = true
}

resource "azurerm_data_factory_integration_runtime_azure" "test" {
  name                    = "azure-integration-runtime"
  data_factory_name       = azurerm_data_factory.test.name
  resource_group_name     = azurerm_resource_group.test.name
  location                = "AutoResolve"
  virtual_network_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t IntegrationRuntimeAzureResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IntegrationRuntimeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.IntegrationRuntimesClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory Azure (%s): %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
