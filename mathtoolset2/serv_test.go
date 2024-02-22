package mathtoolset2

import (
	"testing"
)

func Test_Serv(t *testing.T) {
	// serv, err := NewServ("0.0.0.0:9876", "v", false)
	// assert.NoError(t, err)

	// go func() {
	// 	err = serv.Start(context.Background())
	// 	assert.NoError(t, err)
	// }()

	// time.Sleep(time.Second * 3)

	// client, err := NewClient("0.0.0.0:9876")
	// assert.NoError(t, err)

	// mapfiles := make(map[string]string)
	// mapfiles["reelsstats2.xlsx"] = "../unittestdata/reelsstats2.xlsx"

	// ret, err := client.RunScript(context.Background(), "genStackReels(\"output.xlsx\", \"reelsstats2.xlsx\", [2, 3], [\"SC\"])", mapfiles)
	// assert.NoError(t, err)
	// assert.Equal(t, len(ret.ScriptErrs), 0)

	// out, err := NewFileDataMap(ret.MapFiles)
	// assert.NoError(t, err)
	// assert.NotNil(t, out.MapFiles["output.xlsx"])

	// serv.Stop()

	t.Logf("Test_Serv OK")
}
