package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"testing"
)

func TestJoinJPNICTech(t *testing.T) {
	err := config.GetConfig("/home/yonedayuto/go/src/github.com/homenoc/dsbd-backend/cmd/backend/con.json")
	if err != nil {
		t.Error(err)
	}

	if err = JoinJPNICTech(1, 1); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteJPNICTech(t *testing.T) {
	err := config.GetConfig("/home/yonedayuto/go/src/github.com/homenoc/dsbd-backend/cmd/backend/con.json")
	if err != nil {
		t.Error(err)
	}

	if err = DeleteJPNICTech(1, 1); err != nil {
		t.Fatal(err)
	}
}

func TestAddJPNICTech(t *testing.T) {
	err := config.GetConfig("/home/yonedayuto/go/src/github.com/homenoc/dsbd-backend/cmd/backend/con.json")
	if err != nil {
		t.Error(err)
	}

	data := core.Service{
		GroupID:           0,
		ServiceTemplateID: nil,
		ServiceTemplate:   nil,
		ServiceComment:    "",
		ServiceNumber:     0,
		Org:               "HomeNOC",
		OrgEn:             "",
		Postcode:          "",
		Address:           "",
		AddressEn:         "",
		ASN:               0,
		RouteV4:           "",
		RouteV6:           "",
		V4Name:            "",
		V6Name:            "",
		AveUpstream:       0,
		MaxUpstream:       0,
		AveDownstream:     0,
		MaxDownstream:     0,
		MaxBandWidthAS:    0,
		IP:                nil,
		Connections:       nil,
		JPNICAdminID:      0,
		JPNICAdmin: core.JPNICAdmin{
			Org:   "HomeNOC",
			OrgEn: "HomeNOC",
		},
		JPNICTech: []core.JPNICTech{
			{
				Org:   "HomeNOC",
				OrgEn: "HomeNOC",
			},
		},
		Open: nil,
		Lock: nil,
	}

	result, err := Create(&data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
