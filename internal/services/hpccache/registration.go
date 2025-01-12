package hpccache

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "HPC Cache"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Storage",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_hpc_cache":                 resourceHPCCache(),
		"azurerm_hpc_cache_access_policy":   resourceHPCCacheAccessPolicy(),
		"azurerm_hpc_cache_blob_target":     resourceHPCCacheBlobTarget(),
		"azurerm_hpc_cache_blob_nfs_target": resourceHPCCacheBlobNFSTarget(),
		"azurerm_hpc_cache_nfs_target":      resourceHPCCacheNFSTarget(),
	}
}
