package validate

import (
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/pluginsdk"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/tf/validation"
)

func ServiceBusMaxSizeInMegabytes() pluginsdk.SchemaValidateFunc {
	return validation.IntInSlice([]int{
		1024,
		2048,
		3072,
		4096,
		5120,
		10240,
		20480,
		40960,
		81920,
	})
}

func ServiceBusMaxMessageSizeInKilobytes() pluginsdk.SchemaValidateFunc {
	return validation.IntBetween(1024, 102400)
}
