package synapse

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/2019-06-01-preview/managedvirtualnetwork"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/tf"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	networkValidate "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/network/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/synapse/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/synapse/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func resourceSynapseManagedPrivateEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseManagedPrivateEndpointCreate,
		Read:   resourceSynapseManagedPrivateEndpointRead,
		Delete: resourceSynapseManagedPrivateEndpointDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedPrivateEndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"target_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"subresource_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.PrivateLinkSubResourceName,
			},
		},
	}
}

func resourceSynapseManagedPrivateEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	workspaceClient := meta.(*clients.Client).Synapse.WorkspaceClient
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	workspace, err := workspaceClient.Get(ctx, workspaceId.ResourceGroup, workspaceId.Name)
	if err != nil {
		return fmt.Errorf("retrieving Synapse workspace %q (Resource Group %q): %+v", workspaceId.Name, workspaceId.ResourceGroup, err)
	}
	if workspace.WorkspaceProperties == nil || workspace.WorkspaceProperties.ManagedVirtualNetwork == nil {
		return fmt.Errorf("empty or nil `ManagedVirtualNetwork` for Synapse workspace %q (Resource Group %q): %+v", workspaceId.Name, workspaceId.ResourceGroup, err)
	}
	virtualNetworkName := *workspace.WorkspaceProperties.ManagedVirtualNetwork

	privateEndpointName := d.Get("name").(string)
	client, err := synapseClient.ManagedPrivateEndpointsClient(workspaceId.Name, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	// check exist
	existing, err := client.Get(ctx, virtualNetworkName, privateEndpointName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Synapse Managed Private Endpoint %q (Workspace %q / Resource Group %q): %+v", privateEndpointName, workspaceId.Name, workspaceId.ResourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_synapse_managed_private_endpoint", *existing.ID)
	}

	// create
	managedPrivateEndpoint := managedvirtualnetwork.ManagedPrivateEndpoint{
		Properties: &managedvirtualnetwork.ManagedPrivateEndpointProperties{
			PrivateLinkResourceID: utils.String(d.Get("target_resource_id").(string)),
			GroupID:               utils.String(d.Get("subresource_name").(string)),
		},
	}
	resp, err := client.Create(ctx, virtualNetworkName, privateEndpointName, managedPrivateEndpoint)
	if err != nil {
		return fmt.Errorf("creating Synapse Managed Private Endpoint %q (Workspace %q / Resource Group %q): %+v", privateEndpointName, workspaceId.Name, workspaceId.ResourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Synapse Managed Private Endpoint %q (Workspace %q / Resource Group %q)", privateEndpointName, workspaceId.Name, workspaceId.ResourceGroup)
	}

	d.SetId(*resp.ID)
	return resourceSynapseManagedPrivateEndpointRead(d, meta)
}

func resourceSynapseManagedPrivateEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.ManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.ManagedPrivateEndpointsClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}
	resp, err := client.Get(ctx, id.ManagedVirtualNetworkName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Synapse Managed Private Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Synapse Managed Private Endpoint (Workspace %q / Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroup, err)
	}

	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID()
	d.Set("synapse_workspace_id", workspaceId)
	d.Set("name", id.Name)

	if props := resp.Properties; props != nil {
		d.Set("target_resource_id", props.PrivateLinkResourceID)
		d.Set("subresource_name", props.GroupID)
	}
	return nil
}

func resourceSynapseManagedPrivateEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.ManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.ManagedPrivateEndpointsClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}
	if _, err := client.Delete(ctx, id.ManagedVirtualNetworkName, id.Name); err != nil {
		return fmt.Errorf("deleting Synapse Managed Private Endpoint %q (Workspace %q / Resource Group %q): %+v", id, id.WorkspaceName, id.ResourceGroup, err)
	}

	return nil
}
