package hsm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/hardwaresecuritymodules/mgmt/2018-10-31-preview/hardwaresecuritymodules"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/tf"
	azValidate "github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/location"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/hsm/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/hsm/validate"
	networkValidate "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/network/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func resourceDedicatedHardwareSecurityModule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDedicatedHardwareSecurityModuleCreate,
		Read:   resourceDedicatedHardwareSecurityModuleRead,
		Update: resourceDedicatedHardwareSecurityModuleUpdate,
		Delete: resourceDedicatedHardwareSecurityModuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DedicatedHardwareSecurityModuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DedicatedHardwareSecurityModuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(hardwaresecuritymodules.SafeNetLunaNetworkHSMA790),
				}, false),
			},

			"network_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"network_interface_private_ip_addresses": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: azValidate.IPv4Address,
							},
						},

						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: networkValidate.SubnetID,
						},
					},
				},
			},

			"stamp_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"stamp1",
					"stamp2",
				}, false),
			},

			"zones": azure.SchemaZones(),

			"tags": tags.Schema(),
		},
	}
}

func resourceDedicatedHardwareSecurityModuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Dedicated Hardware Security Module %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_dedicated_hardware_security_module", *existing.ID)
	}

	parameters := hardwaresecuritymodules.DedicatedHsm{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		DedicatedHsmProperties: &hardwaresecuritymodules.DedicatedHsmProperties{
			NetworkProfile: expandDedicatedHsmNetworkProfile(d.Get("network_profile").([]interface{})),
		},
		Sku: &hardwaresecuritymodules.Sku{
			Name: hardwaresecuritymodules.Name(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("stamp_id"); ok {
		parameters.DedicatedHsmProperties.StampID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("zones"); ok {
		parameters.Zones = azure.ExpandZones(v.([]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Dedicated Hardware Security Module %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Dedicated Hardware Security Module %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Dedicated Hardware Security Module %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Dedicated Hardware Security Module %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	return resourceDedicatedHardwareSecurityModuleRead(d, meta)
}

func resourceDedicatedHardwareSecurityModuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DedicatedHardwareSecurityModuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.DedicatedHSMName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Dedicated Hardware Security Module %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Dedicate Hardware Security Module %q (Resource Group %q): %+v", id.DedicatedHSMName, id.ResourceGroup, err)
	}

	d.Set("name", id.DedicatedHSMName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.DedicatedHsmProperties; props != nil {
		if err := d.Set("network_profile", flattenDedicatedHsmNetworkProfile(props.NetworkProfile)); err != nil {
			return fmt.Errorf("setting network_profile: %+v", err)
		}
		d.Set("stamp_id", props.StampID)
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if err := d.Set("zones", resp.Zones); err != nil {
		return fmt.Errorf("setting `zones`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDedicatedHardwareSecurityModuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DedicatedHardwareSecurityModuleID(d.Id())
	if err != nil {
		return err
	}

	parameters := hardwaresecuritymodules.DedicatedHsmPatchParameters{}
	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.DedicatedHSMName, parameters)
	if err != nil {
		return fmt.Errorf("updating Dedicate Hardware Security Module %q (Resource Group %q): %+v", id.DedicatedHSMName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on updating future for Dedicate Hardware Security Module %q (Resource Group %q): %+v", id.DedicatedHSMName, id.ResourceGroup, err)
	}

	return resourceDedicatedHardwareSecurityModuleRead(d, meta)
}

func resourceDedicatedHardwareSecurityModuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DedicatedHardwareSecurityModuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.DedicatedHSMName)
	if err != nil {
		return fmt.Errorf("deleting Dedicated Hardware Security Module %q (Resource Group %q): %+v", id.DedicatedHSMName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Dedicated Hardware Security Module %q (Resource Group %q): %+v", id.DedicatedHSMName, id.ResourceGroup, err)
	}

	return nil
}

func expandDedicatedHsmNetworkProfile(input []interface{}) *hardwaresecuritymodules.NetworkProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := hardwaresecuritymodules.NetworkProfile{
		Subnet: &hardwaresecuritymodules.APIEntityReference{
			ID: utils.String(v["subnet_id"].(string)),
		},
		NetworkInterfaces: expandDedicatedHsmNetworkInterfacePrivateIPAddresses(v["network_interface_private_ip_addresses"].(*pluginsdk.Set).List()),
	}

	return &result
}

func expandDedicatedHsmNetworkInterfacePrivateIPAddresses(input []interface{}) *[]hardwaresecuritymodules.NetworkInterface {
	results := make([]hardwaresecuritymodules.NetworkInterface, 0)

	for _, item := range input {
		if item != nil {
			result := hardwaresecuritymodules.NetworkInterface{
				PrivateIPAddress: utils.String(item.(string)),
			}

			results = append(results, result)
		}
	}

	return &results
}

func flattenDedicatedHsmNetworkProfile(input *hardwaresecuritymodules.NetworkProfile) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var subnetId string
	if input.Subnet != nil && input.Subnet.ID != nil {
		subnetId = *input.Subnet.ID
	}

	return []interface{}{
		map[string]interface{}{
			"network_interface_private_ip_addresses": flattenDedicatedHsmNetworkInterfacePrivateIPAddresses(input.NetworkInterfaces),
			"subnet_id":                              subnetId,
		},
	}
}

func flattenDedicatedHsmNetworkInterfacePrivateIPAddresses(input *[]hardwaresecuritymodules.NetworkInterface) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.PrivateIPAddress != nil {
			results = append(results, *item.PrivateIPAddress)
		}
	}

	return results
}
