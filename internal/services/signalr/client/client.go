package client

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/common"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/signalr/sdk/2020-05-01/signalr"
)

type Client struct {
	Client *signalr.SignalRClient
}

func NewClient(o *common.ClientOptions) *Client {
	client := signalr.NewSignalRClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	return &Client{
		Client: &client,
	}
}
