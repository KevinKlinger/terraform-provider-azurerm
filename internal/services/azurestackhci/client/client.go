package client

import (
	"github.com/Azure/azure-sdk-for-go/services/azurestackhci/mgmt/2020-10-01/azurestackhci"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/common"
)

type Client struct {
	ClusterClient *azurestackhci.ClustersClient
}

func NewClient(o *common.ClientOptions) *Client {
	clusterClient := azurestackhci.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&clusterClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ClusterClient: &clusterClient,
	}
}
