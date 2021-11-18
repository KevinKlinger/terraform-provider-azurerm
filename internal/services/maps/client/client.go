package client

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/common"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/maps/sdk/2021-02-01/accounts"
)

type Client struct {
	AccountsClient *accounts.AccountsClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountsClient := accounts.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&accountsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountsClient: &accountsClient,
	}
}
