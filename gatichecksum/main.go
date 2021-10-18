package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"

	goutils "github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/gatiserv"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.uber.org/zap"
)

func copyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func genSBuild(pathRoot string, binFiles []string, srcFiles []string, gamename string, circleComponents [][]string) {
	starttime := time.Now()
	tsPath := starttime.Format("2006-01-02_15:04:05")

	os.MkdirAll(gamename, os.ModePerm)

	dir, err := ioutil.ReadDir(path.Join(gamename, "bins"))
	if err == nil {
		for _, d := range dir {
			os.RemoveAll(path.Join(gamename, "bins", d.Name()))
		}
	}

	os.MkdirAll(path.Join(gamename, "bins"), os.ModePerm)

	// command := fmt.Sprintf("cd %s\nsource buildgati.docker.sh\nsource buildgati.sh", pathRoot)
	command := fmt.Sprintf("cd %s\nsource buildgati.docker.sh", pathRoot)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		goutils.Error("Command",
			zap.String("command", command),
			zap.Error(err))

		return
	}
	ioutil.WriteFile(path.Join(gamename, "backendBuildOutput.txt"), output, 0644)

	copyFile(path.Join(gamename, "bins", "dtgatigame"), path.Join(pathRoot, "gatidocker/dtgatigame/dtgatigame"))

	os.MkdirAll(path.Join(gamename, tsPath, gamename), os.ModePerm)

	strBin := ""
	for _, v := range binFiles {
		fn := path.Join(pathRoot, v)
		cs, err := gatiserv.Checksum(fn)
		if err != nil {
			goutils.Error("Checksum",
				zap.String("fn", fn),
				zap.Error(err))

			return
		}

		strBin += fmt.Sprintf("%s\t%s\n", cs, path.Join("certificationBuildOutput", tsPath, "bins", v))
	}
	ioutil.WriteFile(path.Join(gamename, "shasumBins.txt"), []byte(strBin), 0644)

	strBin = ""
	for _, v := range binFiles {
		fn := path.Join(pathRoot, v)
		cs, err := gatiserv.Checksum(fn)
		if err != nil {
			goutils.Error("Checksum",
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
	ioutil.WriteFile(path.Join(gamename, tsPath, gamename, "sBuildBinData.txt"), []byte(strBin), 0644)

	strSrc := ""
	for _, v := range srcFiles {
		fn := path.Join(pathRoot, v)
		cs, err := gatiserv.Checksum(fn)
		if err != nil {
			goutils.Error("Checksum",
				zap.String("fn", fn),
				zap.Error(err))

			return
		}

		strSrc += fmt.Sprintf("%s\t%s\n", cs, path.Join("certificationBuildOutput", tsPath, gamename, v))
	}
	ioutil.WriteFile(path.Join(gamename, "shasumSourceCode.txt"), []byte(strSrc), 0644)

	strSrc = ""
	for _, v := range srcFiles {
		fn := path.Join(pathRoot, v)
		cs, err := gatiserv.Checksum(fn)
		if err != nil {
			goutils.Error("Checksum",
				zap.String("fn", fn),
				zap.Error(err))

			return
		}

		strSrc += fmt.Sprintf("%s\t%s\n", cs, path.Join("certificationBuildOutput", tsPath, gamename, v))
		finfo, _ := os.Stat(fn)
		linuxFileAttr := finfo.Sys().(*syscall.Stat_t)
		strSrc += fmt.Sprintf("%s\t%s\n", time.Unix(linuxFileAttr.Mtim.Sec, 0).Format("2006-01-02_15:04:05"), path.Join("certificationBuildOutput", tsPath, "bins", v))
	}
	ioutil.WriteFile(path.Join(gamename, tsPath, gamename, "sBuildSrcData.txt"), []byte(strSrc), 0644)

	strCC := ""
	ccid := 1
	for _, v := range circleComponents {
		fn := path.Join(pathRoot, v[0])
		cs, err := gatiserv.Checksum(fn)
		if err != nil {
			goutils.Error("Checksum",
				zap.String("fn", fn),
				zap.Error(err))

			return
		}

		strCC += fmt.Sprintf("%d\t%s\t%s\t%s\n", ccid, v[1], v[2], cs)

		ccid++
	}
	ioutil.WriteFile(path.Join(gamename, "circleComponents.txt"), []byte(strCC), 0644)

	strBuild := "SBuild started.\n"
	strBuild += starttime.Format("2006-01-02 15:04:05")
	strBuild += "\n"
	strBuild += "SBuild finished.\n"
	strBuild += time.Now().Format("2006-01-02 15:04:05")
	strBuild += "\n"
	ioutil.WriteFile(path.Join(gamename, "executionTime.txt"), []byte(strBuild), 0644)
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

func genElemental2() {
	pathRoot := "../../elemental2/"
	binFiles := []string{
		"gatidocker/dtgatigame/dtgatigame",
	}
	srcFiles := []string{
		"cfg/bg.json",
		"cfg/linedata.json",
		"cfg/paytables.json",
		"cfg/fg0.json",
		"cfg/fg1.json",
		"cfg/fg2.json",
		"cfg/fg3.json",
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
		// "utils.go",
		"version.go",
		"freegame0earth_test.go",
		"freegame0earth.go",
		"freegame1fire.go",
		"freegame2air.go",
		"freegame3water.go",
		"gati/main.go",
		"sbuildsrc/plugin.go",
		"sbuildsrc/rng.go",
	}
	circleComponents := [][]string{
		{"gatidocker/dtgatigame/dtgatigame", "dtgatigame", "app/"},
		{"cfg/bg.json", "bg.json", "app/cfg/"},
		{"cfg/linedata.json", "linedata.json", "app/cfg/"},
		{"cfg/paytables.json", "paytables.json", "app/cfg/"},
		{"cfg/fg0.json", "fg0.json", "app/cfg/"},
		{"cfg/fg1.json", "fg1.json", "app/cfg/"},
		{"cfg/fg2.json", "fg2.json", "app/cfg/"},
		{"cfg/fg3.json", "fg3.json", "app/cfg/"},
		{"cfg/rtp96.yaml", "rtp96.yaml", "app/cfg/"},
	}

	genSBuild(pathRoot, binFiles, srcFiles, "elemental2", circleComponents)
}

func genMedusa2() {
	pathRoot := "../../medusa2/"
	binFiles := []string{
		"gatidocker/dtgatigame/dtgatigame",
	}
	srcFiles := []string{
		"cfg/bg.json",
		"cfg/fg.json",
		"cfg/paytables.json",
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
		"gati/main.go",
		"sbuildsrc/plugin.go",
		"sbuildsrc/rng.go",
	}
	circleComponents := [][]string{
		{"gatidocker/dtgatigame/dtgatigame", "dtgatigame", "app/"},
		{"cfg/bg.json", "bg.json", "app/cfg/"},
		{"cfg/fg.json", "fg.json", "app/cfg/"},
		{"cfg/paytables.json", "paytables.json", "app/cfg/"},
		{"cfg/rtp96.yaml", "rtp96.yaml", "app/cfg/"},
	}

	genSBuild(pathRoot, binFiles, srcFiles, "medusa2", circleComponents)
}

func genDualreel() {
	pathRoot := "../../dualreel/"
	binFiles := []string{
		"gatidocker/dtgatigame/dtgatigame",
	}
	srcFiles := []string{
		"cfg/bgleft.json",
		"cfg/bgright.json",
		"cfg/fgleft.json",
		"cfg/fgright.json",
		"cfg/linedata.json",
		"cfg/paytables.json",
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
		"gati/main.go",
		"sbuildsrc/plugin.go",
		"sbuildsrc/rng.go",
	}
	circleComponents := [][]string{
		{"gatidocker/dtgatigame/dtgatigame", "dtgatigame", "app/"},
		{"cfg/bgleft.json", "bgleft.json", "app/cfg/"},
		{"cfg/bgright.json", "bgright.json", "app/cfg/"},
		{"cfg/fgleft.json", "fgleft.json", "app/cfg/"},
		{"cfg/fgright.json", "fgright.json", "app/cfg/"},
		{"cfg/linedata.json", "linedata.json", "app/cfg/"},
		{"cfg/paytables.json", "paytables.json", "app/cfg/"},
		{"cfg/rtp96.yaml", "rtp96.yaml", "app/cfg/"},
	}

	genSBuild(pathRoot, binFiles, srcFiles, "dualreel", circleComponents)
}

func genToysoldier() {
	pathRoot := "../../toysolider/"
	binFiles := []string{
		"gatidocker/dtgatigame/dtgatigame",
	}
	srcFiles := []string{
		"cfg/reels.json",
		"cfg/weights.json",
		"cfg/lines.json",
		"cfg/paytables.json",
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
		"reels.go",
		"weights.go",
		"version.go",
		"gati/main.go",
		"sbuildsrc/plugin.go",
		"sbuildsrc/rng.go",
	}
	circleComponents := [][]string{
		{"gatidocker/dtgatigame/dtgatigame", "dtgatigame", "app/"},
		{"cfg/reels.json", "reels.json", "app/cfg/"},
		{"cfg/weights.json", "weights.json", "app/cfg/"},
		{"cfg/lines.json", "lines.json", "app/cfg/"},
		{"cfg/paytables.json", "paytables.json", "app/cfg/"},
		{"cfg/rtp96.yaml", "rtp96.yaml", "app/cfg/"},
	}

	genSBuild(pathRoot, binFiles, srcFiles, "toysoldier", circleComponents)
}

func genGloryofheroes() {
	pathRoot := "../../gloryofheroes/"
	binFiles := []string{
		"gatidocker/dtgatigame/dtgatigame",
	}
	srcFiles := []string{
		"cfg/reels.json",
		"cfg/stage.json",
		"cfg/paytables.json",
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
		"reels.go",
		"stage.go",
		"version.go",
		"gati/main.go",
		"sbuildsrc/plugin.go",
		"sbuildsrc/rng.go",
	}
	circleComponents := [][]string{
		{"gatidocker/dtgatigame/dtgatigame", "dtgatigame", "app/"},
		{"cfg/reels.json", "reels.json", "app/cfg/"},
		{"cfg/stage.json", "stage.json", "app/cfg/"},
		{"cfg/paytables.json", "paytables.json", "app/cfg/"},
		{"cfg/rtp96.yaml", "rtp96.yaml", "app/cfg/"},
	}

	genSBuild(pathRoot, binFiles, srcFiles, "gloryofheroes", circleComponents)
}

func main() {
	goutils.InitLogger("gatiChecksum", sgc7ver.Version,
		"debug", true, "./logs")

	// genElemental2()
	// genMedusa2()
	// genLotsalines()
	// genDualreel()
	// genToysoldier()
	genGloryofheroes()
}
