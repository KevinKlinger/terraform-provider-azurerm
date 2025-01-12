package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
)

type BackupProtectionPolicyFileShareDataSource struct {
}

func TestAccDataSourceBackupPolicyFileShare_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("recovery_vault_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (BackupProtectionPolicyFileShareDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_backup_policy_file_share" "test" {
  name                = azurerm_backup_policy_file_share.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, BackupProtectionPolicyFileShareResource{}.basicDaily(data))
}
