package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/recoveryservices/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type BackupProtectionContainerStorageAccountResource struct {
}

func TestAccBackupProtectionContainerStorageAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_container_storage_account", "test")
	r := BackupProtectionContainerStorageAccountResource{}

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

func (t BackupProtectionContainerStorageAccountResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ProtectionContainerID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.RecoveryServices.BackupProtectionContainersClient.Get(ctx, id.VaultName, id.ResourceGroup, id.BackupFabricName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading site recovery protection container (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (BackupProtectionContainerStorageAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "testvlt" {
  name                = "acctest-vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = true
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_backup_container_storage_account" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.testvlt.name
  storage_account_id  = azurerm_storage_account.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
