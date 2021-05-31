package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	ipv4 "github.com/homenoc/dsbd-backend/pkg/api/core/template/ipv4"
	ipv6 "github.com/homenoc/dsbd-backend/pkg/api/core/template/ipv6"
	dbIPv4Template "github.com/homenoc/dsbd-backend/pkg/api/store/template/ipv4/v0"
	dbIPv6Template "github.com/homenoc/dsbd-backend/pkg/api/store/template/ipv6/v0"
	"log"
	"strings"
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
	if input.Name == "" {
		return fmt.Errorf("failed data: [jpnic admin] name")
	}
	if input.NameEn == "" {
		return fmt.Errorf("failed data: [jpnic admin] name(english)")
	}
	if input.Mail == "" || !strings.Contains(input.Mail, "@") {
		return fmt.Errorf("failed data: mail")
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
	if input.Name == "" {
		return fmt.Errorf("failed data: [jpnic tech] name")
	}
	if input.NameEn == "" {
		return fmt.Errorf("failed data: [jpnic tech] name(english)")
	}
	if input.Mail == "" || !strings.Contains(input.Mail, "@") {
		return fmt.Errorf("failed data: mail")
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
			resultIPv4Template := dbIPv4Template.Get(ipv4.Subnet, &core.IPv4Template{Subnet: ip.IP})
			if resultIPv4Template.Err != nil {
				log.Println(resultIPv4Template.Err)
				return resultIPv4Template.Err
			}
			if len(resultIPv4Template.IPv4) == 0 {
				return fmt.Errorf("Invalid IP address or subnet ")
			}

			var after uint = 0
			var halfYear uint = 0
			var oneYear uint = 0

			// Planの計算
			for _, tmp := range ip.Plan {
				after += tmp.After
				halfYear += tmp.HalfYear
				oneYear += tmp.OneYear
			}

			if after < (resultIPv4Template.IPv4[0].Quantity/4) || after > resultIPv4Template.IPv4[0].Quantity {
				return fmt.Errorf("address count error: (after)")
			}
			if halfYear < (resultIPv4Template.IPv4[0].Quantity/4) || halfYear > resultIPv4Template.IPv4[0].Quantity {
				return fmt.Errorf("address count error: (half year)")
			}
			if oneYear < (resultIPv4Template.IPv4[0].Quantity/2) || oneYear > resultIPv4Template.IPv4[0].Quantity {
				return fmt.Errorf("address count error: (one year)")
			}

		}
	} else if ip.Version == 6 {
		if ip.IP == "" {
			return fmt.Errorf("invalid ipv6 address")
		}
		if restrict {
			resultIPv6Template := dbIPv6Template.Get(ipv6.Subnet, &core.IPv6Template{Subnet: ip.IP})
			if resultIPv6Template.Err != nil {
				log.Println(resultIPv6Template.Err)
				return resultIPv6Template.Err
			}
			if len(resultIPv6Template.IPv6) == 0 {
				return fmt.Errorf("Invalid IP address or subnet ")
			}
		}
	} else {
		return fmt.Errorf("invalid ip version")
	}

	return nil
}
