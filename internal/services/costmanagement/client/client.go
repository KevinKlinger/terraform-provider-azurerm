package client

import (
	"github.com/Azure/azure-sdk-for-go/services/costmanagement/mgmt/2019-10-01/costmanagement"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/common"
)

type Client struct {
	ExportClient *costmanagement.ExportsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ExportClient := costmanagement.NewExportsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExportClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ExportClient: &ExportClient,
	}
}
