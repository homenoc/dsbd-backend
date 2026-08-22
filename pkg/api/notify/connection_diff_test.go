package notify_test

import (
	"testing"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/notify"
)

// IX 系フィールドの変更が通知本文に載ることの回帰テスト。
func TestDiffConnectionIXFields(t *testing.T) {
	before := core.Connection{ConnectionType: "CC0", ConnectionNumber: 1}
	after := core.Connection{ConnectionType: "IXP", ConnectionNumber: 1,
		IX: "ENTERNET", IXPeerType: "パブリック", IXVlanID: "306"}

	got := notify.Diff(before, after)
	for _, want := range []string{
		"接続ID: CC0 => IXP\n",
		"接続IX:  => ENTERNET\n",
		"IXピア種別:  => パブリック\n",
		"IX VLAN-ID:  => 306\n",
	} {
		if !contains(got, want) {
			t.Errorf("Diff = %q, want to contain %q", got, want)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return len(sub) == 0
}
