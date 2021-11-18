package sdk

import "github.com/kevinklinger/terraform-provider-azurerm/v2/internal/resourceid"

// SetID uses the specified ID Formatter to set the Resource ID
func (rmd ResourceMetaData) SetID(formatter resourceid.Formatter) {
	rmd.ResourceData.SetId(formatter.ID())
}
