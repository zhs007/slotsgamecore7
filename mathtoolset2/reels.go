package mathtoolset2

import (
	"io"
	"log/slog"
	"strings"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func LoadReels(reader io.Reader) ([][]string, error) {
	reels := [][]string{}

	// x -> ri
	mapri := make(map[int]int)
	maxri := 0
	isend := []bool{}
	isfirst := true

	err := sgc7game.LoadExcelWithReader(reader, "", func(x int, str string) string {
		header := strings.ToLower(strings.TrimSpace(str))
		if header[0] == 'r' {
			iv, err := goutils.String2Int64(header[1:])
			if err != nil {
				goutils.Error("LoadReels:LoadExcelWithReader:String2Int64",
					slog.String("header", header),
					goutils.Err(err))

				return ""
			}

			if iv <= 0 {
				goutils.Error("LoadReels:LoadExcelWithReader",
					slog.String("info", "check iv"),
					slog.String("header", header),
					goutils.Err(sgc7game.ErrInvalidReelsExcelFile))

				return ""
			}

			mapri[x] = int(iv) - 1
			if int(iv) > maxri {
				maxri = int(iv)
			}
		}

		return header
	}, func(x int, y int, header string, data string) error {
		if isfirst {
			isfirst = false

			if maxri != len(mapri) {
				goutils.Error("LoadReels:LoadExcelWithReader",
					slog.String("info", "check len"),
					slog.Int("maxri", maxri),
					slog.Any("mapri", mapri),
					goutils.Err(sgc7game.ErrInvalidReelsExcelFile))

				return sgc7game.ErrInvalidReelsExcelFile
			}

			if maxri <= 0 {
				goutils.Error("LoadReels:LoadExcelWithReader",
					slog.String("info", "check empty"),
					slog.Int("maxri", maxri),
					slog.Any("mapri", mapri),
					goutils.Err(sgc7game.ErrInvalidReelsExcelFile))

				return sgc7game.ErrInvalidReelsExcelFile
			}

			for i := 0; i < maxri; i++ {
				reels = append(reels, []string{})
				isend = append(isend, false)
			}
		}

		ri, isok := mapri[x]
		if isok {
			data = strings.TrimSpace(data)
			if len(data) > 0 {
				if isok {
					reels[ri] = append(reels[ri], data)
				} else {
					isend[ri] = true
				}
			} else {
				isend[ri] = true
			}
		}

		return nil
	})
	if err != nil {
		goutils.Error("LoadReels:LoadExcelWithReader",
			goutils.Err(err))

		return nil, err
	}

	return reels, nil

}
