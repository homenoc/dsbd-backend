package v0

import (
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
