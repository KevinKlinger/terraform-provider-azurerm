package applicationinsights_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/applicationinsights/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type AppInsightsAPIKey struct {
}

func TestAccApplicationInsightsAPIKey_no_permission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basic(data, "[]", "[]"),
			ExpectError: regexp.MustCompile("The API Key needs to have a Role"),
		},
	})
}

func TestAccApplicationInsightsAPIKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "[]", `["annotations"]`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_permissions.#").HasValue("0"),
				check.That(data.ResourceName).Key("write_permissions.#").HasValue("1"),
			),
		},
		{
			Config:      r.requiresImport(data, "[]", `["annotations"]`),
			ExpectError: acceptance.RequiresImportError("azurerm_application_insights_api_key"),
		},
	})
}

func TestAccApplicationInsightsAPIKey_read_telemetry_permissions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, `["aggregate", "api", "draft", "extendqueries", "search"]`, "[]"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_permissions.#").HasValue("5"),
				check.That(data.ResourceName).Key("write_permissions.#").HasValue("0"),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccApplicationInsightsAPIKey_write_annotations_permission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "[]", `["annotations"]`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_permissions.#").HasValue("0"),
				check.That(data.ResourceName).Key("write_permissions.#").HasValue("1"),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccApplicationInsightsAPIKey_authenticate_permission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, `["agentconfig"]`, "[]"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_permissions.#").HasValue("1"),
				check.That(data.ResourceName).Key("write_permissions.#").HasValue("0"),
			),
		},
		data.ImportStep("api_key"),
	})
}

func TestAccApplicationInsightsAPIKey_full_permissions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_api_key", "test")
	r := AppInsightsAPIKey{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, `["agentconfig", "aggregate", "api", "draft", "extendqueries", "search"]`, `["annotations"]`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_permissions.#").HasValue("6"),
				check.That(data.ResourceName).Key("write_permissions.#").HasValue("1"),
			),
		},
		data.ImportStep("api_key"),
	})
}

func (t AppInsightsAPIKey) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ApiKeyID(state.Attributes["id"])
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppInsights.APIKeysClient.Get(ctx, id.ResourceGroup, id.ComponentName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Application Insights API Key '%s' does not exist", id)
	}

	return utils.Bool(resp.StatusCode != http.StatusNotFound), nil
}

func (AppInsightsAPIKey) basic(data acceptance.TestData, readPerms, writePerms string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "test" {
  name                    = "acctestappinsightsapikey-%d"
  application_insights_id = azurerm_application_insights.test.id
  read_permissions        = %s
  write_permissions       = %s
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, readPerms, writePerms)
}

func (AppInsightsAPIKey) requiresImport(data acceptance.TestData, readPerms, writePerms string) string {
	template := AppInsightsAPIKey{}.basic(data, readPerms, writePerms)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_api_key" "import" {
  name                    = azurerm_application_insights_api_key.test.name
  application_insights_id = azurerm_application_insights_api_key.test.application_insights_id
  read_permissions        = azurerm_application_insights_api_key.test.read_permissions
  write_permissions       = azurerm_application_insights_api_key.test.write_permissions
}
`, template)
}
