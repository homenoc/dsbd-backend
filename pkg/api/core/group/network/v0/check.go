package v0

import (
	"fmt"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"time"
)

func check(input network.Input) error {
	// check
	if input.RouteV4 == "" && input.RouteV6 == "" {
		return fmt.Errorf("no data: route(v4 or v6)")
	}
	if input.PI {
		if input.ASN == "" {
			return fmt.Errorf("no data: ASN")
		}
	} else {
		if input.Org == "" {
			return fmt.Errorf("no data: Org")
		}
		if input.OrgEn == "" {
			return fmt.Errorf("no data: Org(English)")
		}
		if input.Postcode == "" {
			return fmt.Errorf("no data: postcode")
		}
		if input.Address == "" {
			return fmt.Errorf("no data: Address")
		}
		if input.AddressEn == "" {
			return fmt.Errorf("no data: Address(English)")
		}
		if len(*input.IP) == 0 {
			return fmt.Errorf("no data: ip address data")
		}
	}
	return nil
}

func ipCheck(ip network.IPInput) error {

	nowTime := time.Now()

	if ip.Version != 4 && ip.Version != 6 {
		return fmt.Errorf("invalid ip version")
	}
	if ip.Name == "" {
		return fmt.Errorf("no network name")
	}

	startDate, _ := time.Parse("2006-01-02", ip.StartDate)
	if startDate.UTC().Unix() < nowTime.UTC().Unix() {
		return fmt.Errorf("invalid start Date")
	}

	if ip.EndDate != nil {
		endDate, _ := time.Parse("2006-01-02", *ip.EndDate)
		if endDate.UTC().Unix() < nowTime.UTC().Unix() && startDate.UTC().Unix() >= endDate.UTC().Unix() {
			return fmt.Errorf("invalid end Date")
		}
	}

	if ip.Version == 4 {
		if ip.IP == "" {
			return fmt.Errorf("invalid ipv4 address")
		}
		if ip.Plan == nil {
			return fmt.Errorf("invalid plan data")
		}
	}

	if ip.Version == 6 {
		if ip.IP == "" {
			return fmt.Errorf("invalid ipv6 address")
		}
	}

	return nil
}
