package client

import (
	"github.com/Azure/azure-sdk-for-go/services/iotcentral/mgmt/2018-09-01/iotcentral"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/common"
)

type Client struct {
	AppsClient *iotcentral.AppsClient
}

func NewClient(o *common.ClientOptions) *Client {
	AppsClient := iotcentral.NewAppsClient(o.SubscriptionId)
	o.ConfigureClient(&AppsClient.Client, o.ResourceManagerAuthorizer)
	return &Client{
		AppsClient: &AppsClient,
	}
}
