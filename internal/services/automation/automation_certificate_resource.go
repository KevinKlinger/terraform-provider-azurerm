package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2018-06-30-preview/automation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/tf"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/automation/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func resourceAutomationCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationCertificateCreateUpdate,
		Read:   resourceAutomationCertificateRead,
		Update: resourceAutomationCertificateCreateUpdate,
		Delete: resourceAutomationCertificateDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CertificateID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"base64": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsBase64,
			},

			"exportable": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
				Optional: true,
			},

			"thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAutomationCertificateCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.CertificateClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Certificate creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("automation_account_name").(string)
	exportable := d.Get("exportable").(bool)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Automation Certificate %q (Account %q / Resource Group %q): %s", name, accountName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_automation_certificate", *existing.ID)
		}
	}

	description := d.Get("description").(string)

	parameters := automation.CertificateCreateOrUpdateParameters{
		Name: &name,
		CertificateCreateOrUpdateProperties: &automation.CertificateCreateOrUpdateProperties{
			Description:  &description,
			IsExportable: &exportable,
		},
	}

	if v, ok := d.GetOk("base64"); ok {
		base64 := v.(string)
		parameters.CertificateCreateOrUpdateProperties.Base64Value = &base64
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, name, parameters); err != nil {
		return fmt.Errorf("creating/updating Certificate %q (Automation Account %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("retrieving Certificate %q (Automation Account %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("ID was nil for Automation Certificate %q (Automation Account %q / Resource Group %q)", name, accountName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceAutomationCertificateRead(d, meta)
}

func resourceAutomationCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.CertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Certificate %q (Automation Account %q / Resource Group %q): %+v", id.Name, id.AutomationAccountName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("automation_account_name", id.AutomationAccountName)

	if props := resp.CertificateProperties; props != nil {
		d.Set("exportable", props.IsExportable)
		d.Set("thumbprint", props.Thumbprint)
		d.Set("description", props.Description)
	}

	return nil
}

func resourceAutomationCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.CertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Certificate %q (Automation Account %q / Resource Group %q): %+v", id.Name, id.AutomationAccountName, id.ResourceGroup, err)
		}
	}

	return nil
}
