package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/jinzhu/gorm"
	"testing"
)

func TestUpdatePartField(t *testing.T) {
	err := config.GetConfig("/home/yonedayuto/go/src/github.com/homenoc/dsbd-backend/cmd/backend/con.json")
	if err != nil {
		t.Error(err)
	}
	testTemplate := core.NOC{
		Model:     gorm.Model{ID: 1},
		Notice:    nil,
		BGPRouter: nil,
		TunnelEndPointRouter: []*core.TunnelEndPointRouter{
			{
				Model:    gorm.Model{ID: 1},
				HostName: "noc01Test",
			},
		},
		Name:      "nocTest",
		Location:  "",
		Bandwidth: "",
		Enable:    nil,
		Comment:   "",
	}
	if err := Update(noc.UpdateAll, testTemplate); err != nil {
		t.Fatal(err)
	}
}
