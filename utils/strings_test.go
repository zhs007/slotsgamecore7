package sgc7utils

import (
	"testing"
)

func Test_AppendString(t *testing.T) {

	in := [][]string{
		{"abc", "efg", "hijklmn"},
		{"abc", "", "hijklmn"},
	}

	out := []string{
		"abcefghijklmn",
		"abchijklmn",
	}

	for i, v := range in {
		ret := AppendString(v...)
		if ret != out[i] {
			t.Fatalf("Test_AppendString AppendString \"%s\" != \"%s\" [ %+v ]",
				ret, out[i], in[i])
		}
	}

	t.Logf("Test_AppendString OK")
}
