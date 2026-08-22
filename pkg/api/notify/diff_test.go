package notify

import (
	"testing"
	"time"
)

type sample struct {
	Org       string     `notify:"Org"`
	AveDown   uint       `notify:"平均ダウンロード帯域,unit=Kbps"`
	AddAllow  *bool      `notify:"接続追加許可,true=許可,false=禁止"`
	ASN       *uint      `notify:"ASN"`
	Expired   *time.Time `notify:"有効期限,date"`
	Ignored   string     // no tag → never reported
	unexports string     `notify:"secret"`
}

func ptrBool(b bool) *bool { return &b }
func ptrUint(u uint) *uint { return &u }

func TestDiff(t *testing.T) {
	exp := time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC)
	exp2 := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name         string
		before       sample
		after        sample
		wantContains []string
		wantAbsent   []string
	}{
		{
			name:         "string change",
			before:       sample{Org: "A"},
			after:        sample{Org: "B"},
			wantContains: []string{"Org: A => B\n"},
		},
		{
			name:         "uint with unit",
			before:       sample{AveDown: 100},
			after:        sample{AveDown: 200},
			wantContains: []string{"平均ダウンロード帯域: 100Kbps => 200Kbps\n"},
		},
		{
			name:         "bool label mapping",
			before:       sample{AddAllow: ptrBool(false)},
			after:        sample{AddAllow: ptrBool(true)},
			wantContains: []string{"接続追加許可: 禁止 => 許可\n"},
		},
		{
			name:         "pointer nil to value",
			before:       sample{ASN: nil},
			after:        sample{ASN: ptrUint(65001)},
			wantContains: []string{"ASN: (なし) => 65001\n"},
		},
		{
			name:         "time as date",
			before:       sample{Expired: &exp},
			after:        sample{Expired: &exp2},
			wantContains: []string{"有効期限: 2026-07-11 => 2027-01-01\n"},
		},
		{
			name:       "unchanged and untagged skipped",
			before:     sample{Org: "same", Ignored: "x"},
			after:      sample{Org: "same", Ignored: "y"},
			wantAbsent: []string{"Org", "Ignored", "secret"},
		},
		{
			name:       "both nil pointer skipped",
			before:     sample{ASN: nil},
			after:      sample{ASN: nil},
			wantAbsent: []string{"ASN"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Diff(tc.before, tc.after)
			for _, w := range tc.wantContains {
				if !contains(got, w) {
					t.Errorf("Diff = %q, want to contain %q", got, w)
				}
			}
			for _, w := range tc.wantAbsent {
				if contains(got, w) {
					t.Errorf("Diff = %q, want NOT to contain %q", got, w)
				}
			}
		})
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func TestBlockText(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty diff", in: "", want: "(変更なし)"},
		{name: "whitespace only", in: " \n\t", want: "(変更なし)"},
		{name: "real diff passes through", in: "接続IX: (なし) => ENTERNET\n", want: "接続IX: (なし) => ENTERNET\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := BlockText(tc.in); got != tc.want {
				t.Errorf("BlockText(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
