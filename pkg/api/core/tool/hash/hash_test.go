package hash

import "testing"

func TestHashProcess(t *testing.T) {
	baseData := "test1test2test3"
	hash := Generate(baseData)
	t.Logf("Hash: %s\n", hash)

	if Verify(baseData, hash) {
		t.Logf("OK\n")
	} else {
		t.Logf("NG\n")
	}
}
