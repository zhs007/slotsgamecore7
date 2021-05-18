package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/zhs007/slotsgamecore7/gatiserv"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.uber.org/zap"
)

func genSBuild(pathRoot string, binFiles []string, srcFiles []string, gamename string, circleComponents [][]string) {
	starttime := time.Now()
	tsPath := starttime.Format("2006-01-02_15:04:05")

	os.MkdirAll(path.Join(tsPath, gamename), os.ModePerm)

	strBin := ""
	for _, v := range binFiles {
		fn := path.Join(pathRoot, v)
		cs, err := gatiserv.Checksum(fn)
		if err != nil {
			sgc7utils.Error("Checksum",
				zap.String("fn", fn),
				zap.Error(err))

			return
		}

		strBin += fmt.Sprintf("%s\t%s\n", cs, path.Join("certificationBuildOutput", tsPath, "bins", v))
	}
	ioutil.WriteFile(path.Join("./", "shasumBins.txt"), []byte(strBin), 0644)

	strBin = ""
	for _, v := range binFiles {
		fn := path.Join(pathRoot, v)
		cs, err := gatiserv.Checksum(fn)
		if err != nil {
			sgc7utils.Error("Checksum",
				zap.String("fn", fn),
				zap.Error(err))

			return
		}

		strBin += fmt.Sprintf("%s\t%s\n", cs, path.Join("certificationBuildOutput", tsPath, "bins", v))

		finfo, _ := os.Stat(fn)
		linuxFileAttr := finfo.Sys().(*syscall.Stat_t)
		strBin += fmt.Sprintf("  File:\t%s\n", path.Join("certificationBuildOutput", tsPath, "bins", v))
		strBin += fmt.Sprintf("Access:\t%s\n", time.Unix(linuxFileAttr.Atim.Sec, 0).Format("2006-01-02_15:04:05"))
		strBin += fmt.Sprintf("Modify:\t%s\n", time.Unix(linuxFileAttr.Mtim.Sec, 0).Format("2006-01-02_15:04:05"))
		strBin += fmt.Sprintf(" Birth:\t%s\n", time.Unix(linuxFileAttr.Ctim.Sec, 0).Format("2006-01-02_15:04:05"))
	}
	ioutil.WriteFile(path.Join(tsPath, gamename, "sBuildBinData.txt"), []byte(strBin), 0644)

	strSrc := ""
	for _, v := range srcFiles {
		fn := path.Join(pathRoot, v)
		cs, err := gatiserv.Checksum(fn)
		if err != nil {
			sgc7utils.Error("Checksum",
				zap.String("fn", fn),
				zap.Error(err))

			return
		}

		strSrc += fmt.Sprintf("%s\t%s\n", cs, path.Join("certificationBuildOutput", tsPath, gamename, v))
	}
	ioutil.WriteFile(path.Join("./", "shasumSourceCode.txt"), []byte(strSrc), 0644)

	strSrc = ""
	for _, v := range srcFiles {
		fn := path.Join(pathRoot, v)
		cs, err := gatiserv.Checksum(fn)
		if err != nil {
			sgc7utils.Error("Checksum",
				zap.String("fn", fn),
				zap.Error(err))

			return
		}

		strSrc += fmt.Sprintf("%s\t%s\n", cs, path.Join("certificationBuildOutput", tsPath, gamename, v))
		finfo, _ := os.Stat(fn)
		linuxFileAttr := finfo.Sys().(*syscall.Stat_t)
		strSrc += fmt.Sprintf("%s\t%s\n", time.Unix(linuxFileAttr.Mtim.Sec, 0).Format("2006-01-02_15:04:05"), path.Join("certificationBuildOutput", tsPath, "bins", v))
	}
	ioutil.WriteFile(path.Join(tsPath, gamename, "sBuildSrcData.txt"), []byte(strSrc), 0644)

	strCC := ""
	ccid := 1
	for _, v := range circleComponents {
		fn := path.Join(pathRoot, v[0])
		cs, err := gatiserv.Checksum(fn)
		if err != nil {
			sgc7utils.Error("Checksum",
				zap.String("fn", fn),
				zap.Error(err))

			return
		}

		strCC += fmt.Sprintf("%d\t%s\t%s\t%s\n", ccid, v[1], v[2], cs)

		ccid++
	}
	ioutil.WriteFile(path.Join("./", "circleComponents.txt"), []byte(strCC), 0644)

	strBuild := "SBuild started.\n"
	strBuild += starttime.Format("2006-01-02 15:04:05")
	strBuild += "\n"
	strBuild += "SBuild finished.\n"
	strBuild += time.Now().Format("2006-01-02 15:04:05")
	strBuild += "\n"
	ioutil.WriteFile(path.Join("./", "executionTime.txt"), []byte(strBuild), 0644)
}

func genLotsalines() {
	pathRoot := "../../lotsalines/"
	binFiles := []string{
		"gatidocker/dtgatigame/dtgatigame",
	}
	srcFiles := []string{
		"cfg/base.json",
		"cfg/linedata.json",
		"cfg/paytables.json",
		"cfg/respin.json",
		"cfg/rtp96.yaml",
		"basedef.go",
		"basegame.go",
		"config.go",
		"err.go",
		"freegame.go",
		"game.go",
		"gatiservice.go",
		"go.mod",
		"go.sum",
		"utils.go",
		"version.go",
	}
	circleComponents := [][]string{
		{"gatidocker/dtgatigame/dtgatigame", "dtgatigame", "app/"},
		{"cfg/base.json", "base.json", "app/cfg/"},
		{"cfg/linedata.json", "linedata.json", "app/cfg/"},
		{"cfg/paytables.json", "paytables.json", "app/cfg/"},
		{"cfg/respin.json", "respin.json", "app/cfg/"},
		{"cfg/rtp96.yaml", "rtp96.yaml", "app/cfg/"},
	}

	genSBuild(pathRoot, binFiles, srcFiles, "lotsalines", circleComponents)
}

func main() {
	sgc7utils.InitLogger("gatiChecksum", sgc7ver.Version,
		"debug", true, "./logs")

	genLotsalines()

	// ccs, err := gatiserv.GenChecksum([]*gatiserv.GATICriticalComponent{
	// 	{
	// 		ID:       1,
	// 		Name:     "dtgatigame",
	// 		Location: "gitlab.heyalgo.io/slotsgames7/lotsalines",
	// 		Filename: "./gatidocker/dtgatigame/dtgatigame",
	// 	},
	// })
	// if err != nil {
	// 	sgc7utils.Error("GenChecksum", zap.Error(err))

	// 	return
	// }

	// gatiserv.SaveGATIGameInfo(&gatiserv.GATIGameInfo{
	// 	Components: ccs.Components,
	// 	Info: gatiserv.VersionInfo{
	// 		GameTitle:     "lotsalines",
	// 		GameVersion:   dtgame.Version,
	// 		VCSVersion:    "6c9ac1ec7bfb8e914456fb2c5476ae3c4a2f425e",
	// 		BuildChecksum: ccs.Components[1].Checksum,
	// 		BuildTime:     time.Now().Format("2006-01-02 15:04:05"),
	// 	},
	// }, "./gatidocker/dtgatigame/gameinfo.json")
}
