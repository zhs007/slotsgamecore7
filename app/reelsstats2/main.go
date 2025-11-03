package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/zhs007/goutils"
	mts2 "github.com/zhs007/slotsgamecore7/mathtoolset2"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
)

// go run ./app/reelsstats2 -i data/bg-reel01.xlsx -o data/bg-reel01-reelsstats2.xlsx

func main() {
	goutils.InitLogger2("reelsstats2", sgc7ver.Version,
		"info", true, "./logs")

	in := flag.String("i", "", "input reels file (xlsx)")
	out := flag.String("o", "", "output reelsstats2 file (xlsx)")
	flag.Parse()

	if *in == "" {
		fmt.Fprintf(os.Stderr, "usage: %s -i input.xlsx [-o output.xlsx]\n", os.Args[0])
		os.Exit(2)
	}

	inpath := *in
	outpath := *out
	if outpath == "" {
		// default: replace ext with -reelsstats2.xlsx
		ext := filepath.Ext(inpath)
		base := strings.TrimSuffix(inpath, ext)
		outpath = base + "-reelsstats2.xlsx"
	}

	f, err := os.Open(inpath)
	if err != nil {
		goutils.Error("reelsstats2:OpenInput",
			slog.String("fn", inpath),
			goutils.Err(err))
		os.Exit(1)
	}
	defer f.Close()

	reels, err := mts2.LoadReels(f)
	if err != nil {
		goutils.Error("reelsstats2:LoadReels",
			slog.String("fn", inpath),
			goutils.Err(err))
		os.Exit(1)
	}

	rss, err := mts2.BuildReelsStats2(reels)
	if err != nil {
		goutils.Error("reelsstats2:BuildReelsStats2",
			goutils.Err(err))
		os.Exit(1)
	}

	err = rss.SaveExcel(outpath)
	if err != nil {
		goutils.Error("reelsstats2:SaveExcel",
			slog.String("fn", outpath),
			goutils.Err(err))
		os.Exit(1)
	}

	goutils.Info("Saved",
		slog.String("output", outpath))
}
