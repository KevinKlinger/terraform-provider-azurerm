package privatedns

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/tf"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/validate"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/services/privatedns/parse"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func resourcePrivateDnsPtrRecord() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateDnsPtrRecordCreateUpdate,
		Read:   resourcePrivateDnsPtrRecordRead,
		Update: resourcePrivateDnsPtrRecordCreateUpdate,
		Delete: resourcePrivateDnsPtrRecordDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PtrRecordID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				// lower-cased due to the broken API https://github.com/Azure/azure-rest-api-specs/issues/6641
				ValidateFunc: validate.LowerCasedString,
			},

			// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/6641
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"zone_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"records": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
			},

			"ttl": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 2147483647),
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourcePrivateDnsPtrRecordCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.RecordSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewPtrRecordID(subscriptionId, d.Get("resource_group_name").(string), d.Get("zone_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.PrivateDnsZoneName, privatedns.PTR, resourceId.PTRName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", resourceId, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_private_dns_ptr_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})

	parameters := privatedns.RecordSet{
		Name: utils.String(resourceId.PTRName),
		RecordSetProperties: &privatedns.RecordSetProperties{
			Metadata:   tags.Expand(t),
			TTL:        &ttl,
			PtrRecords: expandAzureRmPrivateDnsPtrRecords(d),
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.PrivateDnsZoneName, privatedns.PTR, resourceId.PTRName, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
	return resourcePrivateDnsPtrRecordRead(d, meta)
}

func resourcePrivateDnsPtrRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PtrRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Get(ctx, id.ResourceGroup, id.PrivateDnsZoneName, privatedns.PTR, id.PTRName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.PTRName)
	d.Set("zone_name", id.PrivateDnsZoneName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if props := resp.RecordSetProperties; props != nil {
		if err := d.Set("records", flattenAzureRmPrivateDnsPtrRecords(props.PtrRecords)); err != nil {
			return fmt.Errorf("setting `records`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourcePrivateDnsPtrRecordDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PtrRecordID(d.Id())
	if err != nil {
		return err
	}

	if _, err = dnsClient.Delete(ctx, id.ResourceGroup, id.PrivateDnsZoneName, privatedns.PTR, id.PTRName, ""); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func flattenAzureRmPrivateDnsPtrRecords(records *[]privatedns.PtrRecord) []string {
	results := make([]string, 0)

	if records != nil {
		for _, record := range *records {
			if record.Ptrdname == nil {
				continue
			}
			results = append(results, *record.Ptrdname)
		}
	}

	return results
}

func expandAzureRmPrivateDnsPtrRecords(d *pluginsdk.ResourceData) *[]privatedns.PtrRecord {
	recordStrings := d.Get("records").(*pluginsdk.Set).List()
	records := make([]privatedns.PtrRecord, len(recordStrings))

	for i, v := range recordStrings {
		fqdn := v.(string)
		records[i] = privatedns.PtrRecord{
			Ptrdname: &fqdn,
		}
	}

	return &records
}
