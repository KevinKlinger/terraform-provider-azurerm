package billing

import "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Billing"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Billing",
	}
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_billing_enrollment_account_scope": dataSourceBillingEnrollmentAccountScope(),
		"azurerm_billing_mca_account_scope":        dataSourceBillingMCAAccountScope(),
		"azurerm_billing_mpa_account_scope":        dataSourceBillingMPAAccountScope(),
	}
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}
