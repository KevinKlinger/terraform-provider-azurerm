package mysql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/tf"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/locks"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/keyvault/client"
	keyVaultParse "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/keyvault/parse"
	keyVaultValidate "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/keyvault/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/mysql/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/mysql/validate"
	resourcesClient "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/resource/client"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func resourceMySQLServerKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySQLServerKeyCreateUpdate,
		Read:   resourceMySQLServerKeyRead,
		Update: resourceMySQLServerKeyCreateUpdate,
		Delete: resourceMySQLServerKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.KeyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerID,
			},

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.NestedItemId,
			},
		},
	}
}

func getMySQLServerKeyName(ctx context.Context, keyVaultsClient *client.Client, resourcesClient *resourcesClient.Client, keyVaultKeyURI string) (*string, error) {
	keyVaultKeyID, err := keyVaultParse.ParseNestedItemID(keyVaultKeyURI)
	if err != nil {
		return nil, err
	}
	keyVaultIDRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyVaultKeyID.KeyVaultBaseUrl)
	if err != nil {
		return nil, err
	}
	keyVaultID, err := keyVaultParse.VaultID(*keyVaultIDRaw)
	if err != nil {
		return nil, err
	}
	return utils.String(fmt.Sprintf("%s_%s_%s", keyVaultID.Name, keyVaultKeyID.Name, keyVaultKeyID.Version)), nil
}

func resourceMySQLServerKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keysClient := meta.(*clients.Client).MySQL.ServerKeysClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverID, err := parse.ServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}
	keyVaultKeyURI := d.Get("key_vault_key_id").(string)
	name, err := getMySQLServerKeyName(ctx, keyVaultsClient, resourcesClient, keyVaultKeyURI)
	if err != nil {
		return fmt.Errorf("cannot compose name for MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
	}

	locks.ByName(serverID.Name, mySQLServerResourceName)
	defer locks.UnlockByName(serverID.Name, mySQLServerResourceName)

	if d.IsNewResource() {
		// This resource is a singleton, but its name can be anything.
		// If you create a new key with different name with the old key, the service will not give you any warning but directly replace the old key with the new key.
		// Therefore sometimes you cannot get the old key using the GET API since you may not know the name of the old key
		resp, err := keysClient.List(ctx, serverID.ResourceGroup, serverID.Name)
		if err != nil {
			return fmt.Errorf("listing existing MySQL Server Keys in Resource Group %q / Server %q: %+v", serverID.ResourceGroup, serverID.Name, err)
		}
		keys := resp.Values()
		if len(keys) > 1 {
			return fmt.Errorf("expecting at most one MySQL Server Key, but got %q", len(keys))
		}
		if len(keys) == 1 && keys[0].ID != nil && *keys[0].ID != "" {
			return tf.ImportAsExistsError("azurerm_mysql_server_key", *keys[0].ID)
		}
	}

	param := mysql.ServerKey{
		ServerKeyProperties: &mysql.ServerKeyProperties{
			ServerKeyType: utils.String("AzureKeyVault"),
			URI:           &keyVaultKeyURI,
		},
	}

	future, err := keysClient.CreateOrUpdate(ctx, serverID.Name, *name, param, serverID.ResourceGroup)
	if err != nil {
		return fmt.Errorf("creating/updating MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
	}
	if err := future.WaitForCompletionRef(ctx, keysClient.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
	}

	resp, err := keysClient.Get(ctx, serverID.ResourceGroup, serverID.Name, *name)
	if err != nil {
		return fmt.Errorf("retrieving MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for MySQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
	}

	d.SetId(*resp.ID)

	return resourceMySQLServerKeyRead(d, meta)
}

func resourceMySQLServerKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	serversClient := meta.(*clients.Client).MySQL.ServersClient
	keysClient := meta.(*clients.Client).MySQL.ServerKeysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KeyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := keysClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] MySQL Server Key %q was not found (Resource Group %q / Server %q)", id.Name, id.ResourceGroup, id.ServerName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving MySQL Server Key %q (Resource Group %q / Server %q): %+v", id.Name, id.ResourceGroup, id.ServerName, err)
	}

	respServer, err := serversClient.Get(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		return fmt.Errorf("cannot get MySQL Server ID: %+v", err)
	}

	d.Set("server_id", respServer.ID)
	if props := resp.ServerKeyProperties; props != nil {
		d.Set("key_vault_key_id", props.URI)
	}

	return nil
}

func resourceMySQLServerKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServerKeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KeyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ServerName, mySQLServerResourceName)
	defer locks.UnlockByName(id.ServerName, mySQLServerResourceName)

	future, err := client.Delete(ctx, id.ServerName, id.Name, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("deleting MySQL Server Key %q (Resource Group %q / Server %q): %+v", id.Name, id.ResourceGroup, id.ServerName, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of MySQL Server Key %q (Resource Group %q / Server %q): %+v", id.Name, id.ResourceGroup, id.ServerName, err)
	}

	return nil
}
