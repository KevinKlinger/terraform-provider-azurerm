package types

import (
	"context"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
)

type TestResource interface {
	Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error)
}

type TestResourceVerifyingRemoved interface {
	TestResource
	Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error)
}
