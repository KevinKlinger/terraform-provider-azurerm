package bot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/botservice/mgmt/2021-03-01/botservice"
	"github.com/google/uuid"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/acceptance/check"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/bot/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

type BotChannelAlexaResource struct {
}

func testAccBotChannelAlexa_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_alexa", "test")
	r := BotChannelAlexaResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccBotChannelAlexa_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_alexa", "test")
	r := BotChannelAlexaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccBotChannelAlexa_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bot_channel_alexa", "test")
	r := BotChannelAlexaResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t BotChannelAlexaResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ChannelClient.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameAlexaChannel))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotChannelAlexaResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_alexa" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  skill_id            = "amzn1.ask.skill.%s"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), uuid.New().String())
}

func (r BotChannelAlexaResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_alexa" "import" {
  bot_name            = azurerm_bot_channel_alexa.test.bot_name
  location            = azurerm_bot_channel_alexa.test.location
  resource_group_name = azurerm_bot_channel_alexa.test.resource_group_name
  skill_id            = azurerm_bot_channel_alexa.test.skill_id
}
`, r.basic(data))
}

func (BotChannelAlexaResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_alexa" "test" {
  bot_name            = azurerm_bot_channels_registration.test.name
  location            = azurerm_bot_channels_registration.test.location
  resource_group_name = azurerm_resource_group.test.name
  skill_id            = "amzn1.ask.skill.%s"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), uuid.New().String())
}
