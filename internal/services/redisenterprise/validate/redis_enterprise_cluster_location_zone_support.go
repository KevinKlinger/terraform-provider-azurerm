package validate

import (
	"fmt"

	"github.com/kevinklinger/terraform-provider-azurerm/v2/helpers/azure"
	"github.com/kevinklinger/terraform-provider-azurerm/v2/internal/location"
)

// RedisEnterpriseClusterLocationZoneSupport - validates that the passed location supports zones or not
func RedisEnterpriseClusterLocationZoneSupport(input string) error {
	location := location.Normalize(input)
	invalidLocations := invalidRedisEnterpriseClusterZoneLocations()

	for _, str := range invalidLocations {
		if location == str {
			return fmt.Errorf("'Zones' are not currently supported in the %s regions, got %q", azure.QuotedStringSlice(friendlyInvalidRedisEnterpriseClusterZoneLocations()), location)
		}
	}

	return nil
}

func invalidRedisEnterpriseClusterZoneLocations() []string {
	var invalidZone []string

	for _, v := range friendlyInvalidRedisEnterpriseClusterZoneLocations() {
		invalidZone = append(invalidZone, location.Normalize(v))
	}

	return invalidZone
}

func friendlyInvalidRedisEnterpriseClusterZoneLocations() []string {
	return []string{
		"Central US EUAP",
		"West US",
		"Australia Southeast",
		"East Asia",
		"UK West",
		"Central India",
		"South India",
	}
}
