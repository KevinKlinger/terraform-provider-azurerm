package client

import (
	"github.com/Azure/azure-sdk-for-go/services/mixedreality/mgmt/2021-01-01/mixedreality"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/common"
)

type Client struct {
	SpatialAnchorsAccountClient *mixedreality.SpatialAnchorsAccountsClient
}

func NewClient(o *common.ClientOptions) *Client {
	SpatialAnchorsAccountClient := mixedreality.NewSpatialAnchorsAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SpatialAnchorsAccountClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		SpatialAnchorsAccountClient: &SpatialAnchorsAccountClient,
	}
}
