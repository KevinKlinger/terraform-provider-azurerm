package consumption

import (
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/consumption/parse"
	subscriptionParse "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/subscription/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
)

func resourceArmConsumptionBudgetSubscription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmConsumptionBudgetSubscriptionCreateUpdate,
		Read:   resourceArmConsumptionBudgetSubscriptionRead,
		Update: resourceArmConsumptionBudgetSubscriptionCreateUpdate,
		Delete: resourceArmConsumptionBudgetSubscriptionDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ConsumptionBudgetSubscriptionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: SchemaConsumptionBudgetSubscriptionResource(),
	}
}

func resourceArmConsumptionBudgetSubscriptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := subscriptionParse.NewSubscriptionId(d.Get("subscription_id").(string))
	id := parse.NewConsumptionBudgetSubscriptionID(subscriptionId.SubscriptionID, d.Get("name").(string))

	err := resourceArmConsumptionBudgetCreateUpdate(d, meta, consumptionBudgetSubscriptionName, subscriptionId.ID())
	if err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceArmConsumptionBudgetSubscriptionRead(d, meta)
}

func resourceArmConsumptionBudgetSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	consumptionBudgetId, err := parse.ConsumptionBudgetSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	subscriptionId := subscriptionParse.NewSubscriptionId(consumptionBudgetId.SubscriptionId)

	err = resourceArmConsumptionBudgetRead(d, meta, subscriptionId.ID(), consumptionBudgetId.BudgetName)
	if err != nil {
		return err
	}

	d.Set("subscription_id", consumptionBudgetId.SubscriptionId)

	return nil
}

func resourceArmConsumptionBudgetSubscriptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	consumptionBudgetId, err := parse.ConsumptionBudgetSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	subscriptionId := subscriptionParse.NewSubscriptionId(consumptionBudgetId.SubscriptionId)

	return resourceArmConsumptionBudgetDelete(d, meta, subscriptionId.ID())
}
