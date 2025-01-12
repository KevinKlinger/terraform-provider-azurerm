package sql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/sql/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type SqlServerResource struct{}

func TestAccSqlServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_server", "test")
	r := SqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccSqlServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_server", "test")
	r := SqlServerResource{}

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

func TestAccSqlServer_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_server", "test")
	r := SqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccSqlServer_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_server", "test")
	r := SqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		{
			Config: r.withTagsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccSqlServer_withIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_server", "test")
	r := SqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").MatchesRegex(validate.UUIDRegExp),
				check.That(data.ResourceName).Key("identity.0.tenant_id").MatchesRegex(validate.UUIDRegExp),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccSqlServer_updateWithIdentityAdded(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_server", "test")
	r := SqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").MatchesRegex(validate.UUIDRegExp),
				check.That(data.ResourceName).Key("identity.0.tenant_id").MatchesRegex(validate.UUIDRegExp),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccSqlServer_updateWithBlobAuditingPolices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_server", "test")
	r := SqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBlobAuditingPolices(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("extended_auditing_policy.0.storage_account_access_key_is_secondary").HasValue("true"),
				check.That(data.ResourceName).Key("extended_auditing_policy.0.retention_in_days").HasValue("6"),
			),
		},
		data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
		{
			Config: r.withBlobAuditingPolicesUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("extended_auditing_policy.0.storage_account_access_key_is_secondary").HasValue("false"),
				check.That(data.ResourceName).Key("extended_auditing_policy.0.retention_in_days").HasValue("11"),
			),
		},
		data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
	})
}

func TestAccSqlServer_threatDetectionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_server", "test")
	r := SqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.threatDetectionPolicy(data, "Enabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.state").HasValue("Enabled"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.retention_days").HasValue("15"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.disabled_alerts.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.email_account_admins").HasValue("true"),
			),
		},
		data.ImportStep("administrator_login_password", "threat_detection_policy.0.storage_account_access_key"),
		{
			Config: r.threatDetectionPolicy(data, "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.state").HasValue("Disabled"),
			),
		},
	})
}

func (r SqlServerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ServerID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Sql.ServersClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Sql Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r SqlServerResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ServerID(state.ID)
	if err != nil {
		return nil, err
	}
	serversClient := client.Sql.ServersClient
	future, err := serversClient.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("deleting Sql Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, serversClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for deletion of Sql Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r SqlServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r SqlServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_server" "import" {
  name                         = azurerm_sql_server.test.name
  resource_group_name          = azurerm_sql_server.test.resource_group_name
  location                     = azurerm_sql_server.test.location
  version                      = azurerm_sql_server.test.version
  administrator_login          = azurerm_sql_server.test.administrator_login
  administrator_login_password = azurerm_sql_server.test.administrator_login_password
}
`, r.basic(data))
}

func (r SqlServerResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r SqlServerResource) withTagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r SqlServerResource) withIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r SqlServerResource) withBlobAuditingPolices(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  extended_auditing_policy {
    storage_account_access_key              = azurerm_storage_account.test.primary_access_key
    storage_endpoint                        = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }
}
`, data.RandomIntOfLength(15), data.Locations.Primary)
}

func (r SqlServerResource) withBlobAuditingPolicesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  extended_auditing_policy {
    storage_account_access_key              = azurerm_storage_account.test2.primary_access_key
    storage_endpoint                        = azurerm_storage_account.test2.primary_blob_endpoint
    storage_account_access_key_is_secondary = false
    retention_in_days                       = 11
  }
}
`, data.RandomIntOfLength(15), data.Locations.Primary)
}

func (r SqlServerResource) threatDetectionPolicy(data acceptance.TestData, state string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "test%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  threat_detection_policy {
    retention_days             = 15
    state                      = "%s"
    disabled_alerts            = ["Sql_Injection"]
    email_account_admins       = true
    storage_account_access_key = azurerm_storage_account.test.primary_access_key
    storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, state)
}
