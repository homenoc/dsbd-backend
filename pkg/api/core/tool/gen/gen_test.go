package gen

import "testing"

func TestGenerateUUIDString(t *testing.T) {
	uuid, err := GenerateUUIDString()
	if err != nil {
		t.Fatalf("error: %s\n", err)
	}
	t.Logf("UUID: %s\n", uuid)
}
