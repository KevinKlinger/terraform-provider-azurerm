package datalake_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/datalake/sdk/datalakeanalytics/2016-11-01/accounts"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type DataLakeAnalyticsAccountResource struct {
}

func TestAccDataLakeAnalyticsAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_account", "test")
	r := DataLakeAnalyticsAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Consumption"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataLakeAnalyticsAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_account", "test")
	r := DataLakeAnalyticsAccountResource{}

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

func TestAccDataLakeAnalyticsAccount_tier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_account", "test")
	r := DataLakeAnalyticsAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tier(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Commitment_100AUHours"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataLakeAnalyticsAccount_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_account", "test")
	r := DataLakeAnalyticsAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (t DataLakeAnalyticsAccountResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := accounts.ParseAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Datalake.AnalyticsAccountsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (DataLakeAnalyticsAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  default_store_account_name = azurerm_data_lake_store.test.name
}
`, DataLakeStoreResource{}.basic(data), strconv.Itoa(data.RandomInteger)[2:17])
}

func (r DataLakeAnalyticsAccountResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "import" {
  name                       = azurerm_data_lake_analytics_account.test.name
  resource_group_name        = azurerm_data_lake_analytics_account.test.resource_group_name
  location                   = azurerm_data_lake_analytics_account.test.location
  default_store_account_name = azurerm_data_lake_analytics_account.test.default_store_account_name
}
`, DataLakeAnalyticsAccountResource{}.basic(data))
}

func (r DataLakeAnalyticsAccountResource) tier(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tier = "Commitment_100AUHours"

  default_store_account_name = azurerm_data_lake_store.test.name
}
`, DataLakeStoreResource{}.basic(data), strconv.Itoa(data.RandomInteger)[2:17])
}

func (r DataLakeAnalyticsAccountResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  default_store_account_name = azurerm_data_lake_store.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, DataLakeStoreResource{}.basic(data), strconv.Itoa(data.RandomInteger)[2:17])
}

func (r DataLakeAnalyticsAccountResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  default_store_account_name = azurerm_data_lake_store.test.name

  tags = {
    environment = "staging"
  }
}
`, DataLakeStoreResource{}.basic(data), strconv.Itoa(data.RandomInteger)[2:17])
}
