package client

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/common"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/msi/sdk/2018-11-30/managedidentity"
)

type Client struct {
	UserAssignedIdentitiesClient *managedidentity.ManagedIdentityClient
}

func NewClient(o *common.ClientOptions) *Client {
	UserAssignedIdentitiesClient := managedidentity.NewManagedIdentityClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&UserAssignedIdentitiesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		UserAssignedIdentitiesClient: &UserAssignedIdentitiesClient,
	}
}
