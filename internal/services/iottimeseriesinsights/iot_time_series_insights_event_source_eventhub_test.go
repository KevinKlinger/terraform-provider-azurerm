package iottimeseriesinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/iottimeseriesinsights/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type IoTTimeSeriesInsightsEventSourceEventhubResource struct {
}

func TestAccIoTTimeSeriesInsightsEventSourceEventhub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_time_series_insights_event_source_eventhub", "test")
	r := IoTTimeSeriesInsightsEventSourceEventhubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_key"),
	})
}

func TestAccIoTTimeSeriesInsightsEventSourceEventhub_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_time_series_insights_event_source_eventhub", "test")
	r := IoTTimeSeriesInsightsEventSourceEventhubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_key"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_key"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (IoTTimeSeriesInsightsEventSourceEventhubResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.EventSourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.IoTTimeSeriesInsights.EventSourcesClient.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving IoT Time Series Insights EventSource Eventhub (%q): %+v", id.String(), err)
	}

	return utils.Bool(!utils.ResponseWasNotFound(resp.Response)), nil
}

func (IoTTimeSeriesInsightsEventSourceEventhubResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-tsi-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 7
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_eventhub_authorization_rule" "test" {
  name                = "acctest-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = true
  send   = false
  manage = false
}

resource "azurerm_storage_account" "storage" {
  name                     = "acctestsatsi%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iot_time_series_insights_gen2_environment" "test" {
  name                = "acctest_tsie%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "L1"
  id_properties       = ["id"]

  storage {
    name = azurerm_storage_account.storage.name
    key  = azurerm_storage_account.storage.primary_access_key
  }
}

resource "azurerm_iot_time_series_insights_event_source_eventhub" "test" {
  name                     = "acctest_tsiesi%d"
  location                 = azurerm_resource_group.test.location
  environment_id           = azurerm_iot_time_series_insights_gen2_environment.test.id
  eventhub_name            = azurerm_eventhub.test.name
  namespace_name           = azurerm_eventhub_namespace.test.name
  shared_access_key        = azurerm_eventhub_authorization_rule.test.primary_key
  shared_access_key_name   = azurerm_eventhub_authorization_rule.test.name
  consumer_group_name      = azurerm_eventhub_consumer_group.test.name
  event_source_resource_id = azurerm_eventhub.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (IoTTimeSeriesInsightsEventSourceEventhubResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-tsi-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 7
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_eventhub_authorization_rule" "test" {
  name                = "acctest-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = true
  send   = false
  manage = false
}

resource "azurerm_storage_account" "storage" {
  name                     = "acctestsatsi%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iot_time_series_insights_gen2_environment" "test" {
  name                = "acctest_tsie%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "L1"
  id_properties       = ["id"]

  storage {
    name = azurerm_storage_account.storage.name
    key  = azurerm_storage_account.storage.primary_access_key
  }
}

resource "azurerm_iot_time_series_insights_event_source_eventhub" "test" {
  name                     = "acctest_tsiesi%d"
  location                 = azurerm_resource_group.test.location
  environment_id           = azurerm_iot_time_series_insights_gen2_environment.test.id
  eventhub_name            = azurerm_eventhub.test.name
  namespace_name           = azurerm_eventhub_namespace.test.name
  shared_access_key        = azurerm_eventhub_authorization_rule.test.primary_key
  shared_access_key_name   = azurerm_eventhub_authorization_rule.test.name
  consumer_group_name      = azurerm_eventhub_consumer_group.test.name
  event_source_resource_id = azurerm_eventhub.test.id
  timestamp_property_name  = "test"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}
