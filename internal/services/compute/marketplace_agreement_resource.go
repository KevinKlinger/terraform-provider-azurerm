package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/tf"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/compute/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func resourceMarketplaceAgreement() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMarketplaceAgreementCreateUpdate,
		Read:   resourceMarketplaceAgreementRead,
		Delete: resourceMarketplaceAgreementDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PlanID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"offer": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"plan": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"publisher": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"license_text_link": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"privacy_policy_link": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMarketplaceAgreementCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	offer := d.Get("offer").(string)
	plan := d.Get("plan").(string)
	publisher := d.Get("publisher").(string)

	log.Printf("[DEBUG] Retrieving the Marketplace Terms for Publisher %q / Offer %q / Plan %q", publisher, offer, plan)

	term, err := client.Get(ctx, publisher, offer, plan)
	if err != nil {
		if !utils.ResponseWasNotFound(term.Response) {
			return fmt.Errorf("retrieving the Marketplace Terms for Publisher %q / Offer %q / Plan %q: %s", publisher, offer, plan, err)
		}
	}

	accepted := false
	if props := term.AgreementProperties; props != nil {
		if acc := props.Accepted; acc != nil {
			accepted = *acc
		}
	}

	if accepted {
		agreement, err := client.GetAgreement(ctx, publisher, offer, plan)
		if err != nil {
			if !utils.ResponseWasNotFound(agreement.Response) {
				return fmt.Errorf("retrieving agreement for Publisher %q / Offer %q / Plan %q: %s", publisher, offer, plan, err)
			}
		}
		return tf.ImportAsExistsError("azurerm_marketplace_agreement", *agreement.ID)
	}

	terms, err := client.Get(ctx, publisher, offer, plan)
	if err != nil {
		return fmt.Errorf("retrieving the Marketplace Terms for Publisher %q / Offer %q / Plan %q: %s", publisher, offer, plan, err)
	}
	if terms.AgreementProperties == nil {
		return fmt.Errorf("retrieving the Marketplace Terms for Publisher %q / Offer %q / Plan %q: AgreementProperties was nil", publisher, offer, plan)
	}

	terms.AgreementProperties.Accepted = utils.Bool(true)

	log.Printf("[DEBUG] Accepting the Marketplace Terms for Publisher %q / Offer %q / Plan %q", publisher, offer, plan)
	if _, err := client.Create(ctx, publisher, offer, plan, terms); err != nil {
		return fmt.Errorf("accepting Terms for Publisher %q / Offer %q / Plan %q: %s", publisher, offer, plan, err)
	}
	log.Printf("[DEBUG] Accepted the Marketplace Terms for Publisher %q / Offer %q / Plan %q", publisher, offer, plan)

	agreement, err := client.GetAgreement(ctx, publisher, offer, plan)
	if err != nil {
		return fmt.Errorf("retrieving agreement for Publisher %q / Offer %q / Plan %q: %s", publisher, offer, plan, err)
	}

	d.SetId(*agreement.ID)

	return resourceMarketplaceAgreementRead(d, meta)
}

func resourceMarketplaceAgreementRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PlanID(d.Id())
	if err != nil {
		return err
	}

	term, err := client.Get(ctx, id.AgreementName, id.OfferName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(term.Response) {
			log.Printf("[DEBUG] The Marketplace Terms was not found for Publisher %q / Offer %q / Plan %q", id.AgreementName, id.OfferName, id.Name)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving the Marketplace Terms for Publisher %q / Offer %q / Plan %q: %s", id.AgreementName, id.OfferName, id.Name, err)
	}

	d.Set("publisher", id.AgreementName)
	d.Set("offer", id.OfferName)
	d.Set("plan", id.Name)

	if props := term.AgreementProperties; props != nil {
		if accepted := props.Accepted != nil && *props.Accepted; !accepted {
			// if props.Accepted is not true, the agreement does not exist
			d.SetId("")
		}
		d.Set("license_text_link", props.LicenseTextLink)
		d.Set("privacy_policy_link", props.PrivacyPolicyLink)
	}

	return nil
}

func resourceMarketplaceAgreementDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PlanID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Cancel(ctx, id.AgreementName, id.OfferName, id.Name); err != nil {
		return fmt.Errorf("cancelling agreement for Publisher %q / Offer %q / Plan %q: %s", id.AgreementName, id.OfferName, id.Name, err)
	}

	return nil
}
