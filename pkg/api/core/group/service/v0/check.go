package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"time"
)

func check(input service.Input) error {
	if input.AveUpstream == 0 {
		return fmt.Errorf("no data: avg upstream")
	}
	if input.MaxUpstream == 0 {
		return fmt.Errorf("no data: max upstream")
	}
	if input.AveDownstream == 0 {
		return fmt.Errorf("no data: avg downstream")
	}
	if input.MaxDownstream == 0 {
		return fmt.Errorf("no data: max downstream")
	}

	return nil
}

func checkJPNIC(input service.Input) error {
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

	return nil
}

func checkJPNICAdminUser(input core.JPNICAdmin) error {
	// check
	if input.Org == "" {
		return fmt.Errorf("failed data: org")
	}
	if input.OrgEn == "" {
		return fmt.Errorf("failed data: org(english)")
	}
	if input.PostCode == "" {
		return fmt.Errorf("failed data: postcode")
	}
	if input.Address == "" {
		return fmt.Errorf("failed data: address")
	}
	if input.AddressEn == "" {
		return fmt.Errorf("failed data: address(english)")
	}
	if input.Tel == "" {
		return fmt.Errorf("failed data: tel")
	}
	return nil
}

func checkJPNICTechUser(input core.JPNICTech) error {
	// check
	if input.Org == "" {
		return fmt.Errorf("failed data: org")
	}
	if input.OrgEn == "" {
		return fmt.Errorf("failed data: org(english)")
	}
	if input.PostCode == "" {
		return fmt.Errorf("failed data: postcode")
	}
	if input.Address == "" {
		return fmt.Errorf("failed data: address")
	}
	if input.AddressEn == "" {
		return fmt.Errorf("failed data: address(english)")
	}
	if input.Tel == "" {
		return fmt.Errorf("failed data: tel")
	}
	return nil
}

func ipCheck(admin, restrict bool, ip service.IPInput) error {

	nowTime := time.Now()

	if ip.Version != 4 && ip.Version != 6 {
		return fmt.Errorf("invalid ip version")
	}
	// 厳格な場合
	if restrict {
		if ip.Name == "" {
			return fmt.Errorf("no network name")
		}
	}
	if !admin {
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
	}

	if ip.Version == 4 {
		if ip.IP == "" {
			return fmt.Errorf("invalid ipv4 address")
		}
		if restrict {
			if ip.Plan == nil {
				return fmt.Errorf("invalid plan data")
			}
		} else {
			// Planの計算

			//for _, tmp := range ip.Plan {
			//	tmp.
			//}
		}
	} else if ip.Version == 6 {
		if ip.IP == "" {
			return fmt.Errorf("invalid ipv6 address")
		}
	} else {
		return fmt.Errorf("invalid ip version")
	}

	return nil
}
