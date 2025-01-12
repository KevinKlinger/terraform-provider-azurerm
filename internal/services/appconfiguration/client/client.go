package client

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/common"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/appconfiguration/sdk/1.0/appconfiguration"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/appconfiguration/sdk/2020-06-01/configurationstores"
)

type Client struct {
	ConfigurationStoresClient *configurationstores.ConfigurationStoresClient
	tokenFunc                 func(endpoint string) (autorest.Authorizer, error)
	configureClientFunc       func(c *autorest.Client, authorizer autorest.Authorizer)
}

func (c Client) DataPlaneClient(ctx context.Context, configurationStoreId string) (*appconfiguration.BaseClient, error) {
	appConfigId, err := configurationstores.ParseConfigurationStoreID(configurationStoreId)
	if err != nil {
		return nil, err
	}

	// TODO: caching all of this
	appConfig, err := c.ConfigurationStoresClient.Get(ctx, *appConfigId)
	if err != nil {
		// TODO: if not found etc
		return nil, err
	}

	if appConfig.Model == nil || appConfig.Model.Properties == nil || appConfig.Model.Properties.Endpoint == nil {
		return nil, fmt.Errorf("endpoint was nil")
	}

	endpoint := *appConfig.Model.Properties.Endpoint
	appConfigAuth, err := c.tokenFunc(endpoint)
	if err != nil {
		return nil, fmt.Errorf("obtaining auth token for %q: %+v", endpoint, err)
	}

	client := appconfiguration.NewWithoutDefaults("", endpoint)
	c.configureClientFunc(&client.Client, appConfigAuth)
	return &client, nil
}

func NewClient(o *common.ClientOptions) *Client {
	configurationStores := configurationstores.NewConfigurationStoresClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&configurationStores.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationStoresClient: &configurationStores,
		tokenFunc:                 o.TokenFunc,
		configureClientFunc:       o.ConfigureClient,
	}
}
