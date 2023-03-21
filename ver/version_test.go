package sgc7ver

import (
	"os"
	"testing"
)

func Test_Version(t *testing.T) {
	data, err := os.ReadFile("../VERSION")
	if err != nil {
		t.Fatalf("Test_Version ReadFile error %+v",
			err)
	}

	if string(data) != Version {
		t.Fatalf("Test_Version VERSION error %s != %s",
			string(data),
			Version)
	}

	t.Logf("Test_Version OK")
}
