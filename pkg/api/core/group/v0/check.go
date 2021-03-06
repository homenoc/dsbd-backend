package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"time"
)

func check(input group.Input) error {
	// check
	if input.Question == "" {
		return fmt.Errorf("no data: question")
	}
	if input.Org == "" {
		return fmt.Errorf("no data: org")
	}
	if input.OrgEn == "" {
		return fmt.Errorf("no data: org(english)")
	}
	if input.PostCode == "" {
		return fmt.Errorf("no data: postcode")
	}
	if input.Address == "" {
		return fmt.Errorf("no data: address")
	}
	if input.AddressEn == "" {
		return fmt.Errorf("no data: address(english)")
	}
	if input.Tel == "" {
		return fmt.Errorf("no data: tel")
	}
	if input.Country == "" {
		return fmt.Errorf("no data: country")
	}
	if input.Contract == "" {
		return fmt.Errorf("no data: contract")
	}
	if *input.Student {
		studentExpired, err := time.Parse("2006-01-02", *input.StudentExpired)
		if err != nil {
			return err
		}
		if studentExpired.Unix() < time.Now().Unix() {
			return fmt.Errorf("時間指定に不備があります。")
		}
	}

	return nil
}
