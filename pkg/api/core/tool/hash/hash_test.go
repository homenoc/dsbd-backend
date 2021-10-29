package hash

import "testing"

func TestHashProcess(t *testing.T) {
	baseData := "test1test2test3"
	hash := Generate(baseData)
	t.Logf("Hash: %s", hash)

	if Verify(baseData, hash) {
		t.Logf("OK")
	} else {
		t.Logf("NG")
	}
}
