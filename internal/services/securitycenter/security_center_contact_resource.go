package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/tf"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/securitycenter/azuresdkhacks"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

// seems you can only set one contact:
// Invalid security contact name was provided - only 'defaultX' is allowed where X is an index
// Invalid security contact name 'default0' was provided. Expected 'default1'
// Message="Invalid security contact name 'default2' was provided. Expected 'default1'"
const securityCenterContactName = "default1"

func resourceSecurityCenterContact() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSecurityCenterContactCreateUpdate,
		Read:   resourceSecurityCenterContactRead,
		Update: resourceSecurityCenterContactCreateUpdate,
		Delete: resourceSecurityCenterContactDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"email": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"phone": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"alert_notifications": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"alerts_to_admins": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceSecurityCenterContactCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := securityCenterContactName

	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Checking for presence of existing Security Center Contact: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_security_center_contact", *existing.ID)
		}
	}

	contact := security.Contact{
		ContactProperties: &security.ContactProperties{
			Email: utils.String(d.Get("email").(string)),
			Phone: utils.String(d.Get("phone").(string)),
		},
	}

	if alertNotifications := d.Get("alert_notifications").(bool); alertNotifications {
		contact.AlertNotifications = security.On
	} else {
		contact.AlertNotifications = security.Off
	}

	if alertNotifications := d.Get("alerts_to_admins").(bool); alertNotifications {
		contact.AlertsToAdmins = security.AlertsToAdminsOn
	} else {
		contact.AlertsToAdmins = security.AlertsToAdminsOff
	}

	if d.IsNewResource() {
		// TODO: switch back when the Swagger/API bug has been fixed:
		// https://github.com/Azure/azure-rest-api-specs/issues/10717 (an undefined 201)
		if _, err := azuresdkhacks.CreateSecurityCenterContact(client, ctx, name, contact); err != nil {
			return fmt.Errorf("Creating Security Center Contact: %+v", err)
		}

		resp, err := client.Get(ctx, name)
		if err != nil {
			return fmt.Errorf("Reading Security Center Contact: %+v", err)
		}
		if resp.ID == nil {
			return fmt.Errorf("Security Center Contact ID is nil")
		}

		d.SetId(*resp.ID)
	} else if _, err := client.Update(ctx, name, contact); err != nil {
		return fmt.Errorf("Updating Security Center Contact: %+v", err)
	}

	return resourceSecurityCenterContactRead(d, meta)
}

func resourceSecurityCenterContactRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := securityCenterContactName

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Security Center Subscription Contact was not found: %v", err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Reading Security Center Contact: %+v", err)
	}

	if properties := resp.ContactProperties; properties != nil {
		d.Set("email", properties.Email)
		d.Set("phone", properties.Phone)
		d.Set("alert_notifications", properties.AlertNotifications == security.On)
		d.Set("alerts_to_admins", properties.AlertsToAdmins == security.AlertsToAdminsOn)
	}

	return nil
}

func resourceSecurityCenterContactDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := securityCenterContactName

	resp, err := client.Delete(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] Security Center Subscription Contact was not found: %v", err)
			return nil
		}

		return fmt.Errorf("Deleting Security Center Contact: %+v", err)
	}

	return nil
}
