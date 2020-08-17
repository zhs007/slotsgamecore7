package sgc7ver

import (
	"io/ioutil"
	"testing"
)

func Test_Version(t *testing.T) {
	data, err := ioutil.ReadFile("../VERSION")
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
