package client

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/common"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/relay/sdk/2017-04-01/hybridconnections"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/relay/sdk/2017-04-01/namespaces"
)

type Client struct {
	HybridConnectionsClient *hybridconnections.HybridConnectionsClient
	NamespacesClient        *namespaces.NamespacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	hybridConnectionsClient := hybridconnections.NewHybridConnectionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&hybridConnectionsClient.Client, o.ResourceManagerAuthorizer)

	namespacesClient := namespaces.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&namespacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HybridConnectionsClient: &hybridConnectionsClient,
		NamespacesClient:        &namespacesClient,
	}
}
