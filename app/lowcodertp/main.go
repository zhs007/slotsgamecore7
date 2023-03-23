package main

import (
	"os"
	"strconv"

	"net/http"
	_ "net/http/pprof"

	_ "github.com/mkevac/debugcharts"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.uber.org/zap"
)

func main() {
	if os.Getenv("PPROF") == "true" {
		go func() {
			// terminal: $ go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap
			// web:
			// 1、http://localhost:8081/ui
			// 2、http://localhost:6060/debug/charts
			// 3、http://localhost:6060/debug/pprof
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}

	strcore := os.Getenv("CORE")
	if strcore == "" {
		strcore = "8"
	}

	strspinnums := os.Getenv("SPINNUMS")
	if strspinnums == "" {
		strspinnums = "10000000"
	}

	gamecfg := os.Getenv("GAMECFG")
	outputPath := os.Getenv("OUTPUTPATH")

	goutils.InitLogger("lowcodertp", sgc7ver.Version,
		"info", true, "./logs")

	icore, err := strconv.Atoi(strcore)
	if err != nil {
		goutils.Error("Getenv(CORE)",
			zap.Error(err))

		return
	}

	ispinnums, err := strconv.ParseInt(strspinnums, 10, 64)
	if err != nil {
		goutils.Error("Getenv(SPINNUMS)",
			zap.Error(err))

		return
	}

	lowcode.StartRTP(gamecfg, icore, ispinnums, outputPath)
}