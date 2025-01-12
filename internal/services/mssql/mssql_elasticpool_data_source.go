package mssql

import (
	"fmt"
	"time"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/clients"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tags"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/timeouts"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/utils"
)

func dataSourceMsSqlElasticpool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMsSqlElasticpoolRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"server_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"max_size_bytes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"max_size_gb": {
				Type:     pluginsdk.TypeFloat,
				Computed: true,
			},

			"per_db_min_capacity": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"per_db_max_capacity": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"sku": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"capacity": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"tier": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"family": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),

			"zone_redundant": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"license_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMsSqlElasticpoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	elasticPoolName := d.Get("name").(string)
	serverName := d.Get("server_name").(string)

	resp, err := client.Get(ctx, resourceGroup, serverName, elasticPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Elasticpool %q (Resource Group %q, SQL Server %q) was not found", elasticPoolName, resourceGroup, serverName)
		}

		return fmt.Errorf("making Read request on AzureRM Elasticpool %s (Resource Group %q, SQL Server %q): %+v", elasticPoolName, resourceGroup, serverName, err)
	}

	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}
	d.Set("name", elasticPoolName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("server_name", serverName)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ElasticPoolProperties; props != nil {
		d.Set("max_size_gb", float64(*props.MaxSizeBytes/int64(1073741824)))
		d.Set("max_size_bytes", props.MaxSizeBytes)

		d.Set("zone_redundant", props.ZoneRedundant)
		d.Set("license_type", props.LicenseType)

		if perDbSettings := props.PerDatabaseSettings; perDbSettings != nil {
			d.Set("per_db_min_capacity", perDbSettings.MinCapacity)
			d.Set("per_db_max_capacity", perDbSettings.MaxCapacity)
		}
	}

	if err := d.Set("sku", flattenMsSqlElasticPoolSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
