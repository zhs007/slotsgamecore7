package mathtoolset2

import (
	"io"
	"log/slog"
	"strings"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type ReelStats2 struct {
	MapSymbols     map[string]int
	TotalSymbolNum int
}

func (rs2 *ReelStats2) genSymbolsPool() (*SymbolsPool, error) {
	pool := &SymbolsPool{}

	for k, v := range rs2.MapSymbols {
		if strings.Contains(k, "_") {
			arr := strings.Split(k, "_")
			if len(arr) != 2 {
				goutils.Error("ReelStats2.genSymbolsPool:Split",
					slog.Any("arr", arr))

				return nil, ErrUnkonow
			}

			v1, err := goutils.String2Int64(arr[1])
			if err != nil {
				goutils.Error("ReelStats2.genSymbolsPool:Split",
					slog.Any("arr", arr),
					goutils.Err(err))

				return nil, err
			}

			pool.PushEx(arr[0], int(v1), v)
		} else {
			pool.PushEx(k, 1, v)
		}
	}

	return pool, nil
}

type ReelsStats2 struct {
	Reels   []*ReelStats2
	Symbols []string
}

func NewReelsStats2(reelnum int) *ReelsStats2 {
	rss2 := &ReelsStats2{}

	for i := 0; i < reelnum; i++ {
		rs2 := &ReelStats2{
			MapSymbols: make(map[string]int),
		}

		rss2.Reels = append(rss2.Reels, rs2)
	}

	return rss2
}

func getReelID(str string) (int, error) {
	arr := strings.Split(str, "reel")
	if len(arr) == 2 {
		i64, err := goutils.String2Int64(arr[1])
		if err != nil {
			goutils.Error("getReelID",
				goutils.Err(err))

			return -1, err
		}

		return int(i64), nil
	}

	return -1, nil
}

func LoadReelsStats2(reader io.Reader) (*ReelsStats2, error) {
	mapSymbols := make(map[string][]int)
	reelnum := 0
	curSymbol := ""
	err := sgc7game.LoadExcelWithReader(reader, "", func(x int, str string) string {
		head := strings.ToLower(strings.TrimSpace(str))
		rn, err := getReelID(head)
		if err != nil {
			goutils.Error("LoadReelsStats2:LoadExcel:onheader:getReelID",
				goutils.Err(err))
		}

		if rn > reelnum {
			reelnum = rn
		}

		return head
	}, func(x int, y int, header string, data string) error {
		if reelnum <= 0 {
			goutils.Error("LoadReelsStats2:LoadExcel:ondata:reelnum",
				slog.Int("reelnum", reelnum),
				goutils.Err(ErrInvalidReelsStats2File))

			return ErrInvalidReelsStats2File
		}

		data = strings.TrimSpace(data)
		if header == "symbol" {
			curSymbol = data
			mapSymbols[curSymbol] = make([]int, reelnum)
		} else {
			rn, err := getReelID(header)
			if err != nil {
				goutils.Error("LoadReelsStats2:LoadExcel:ondata:getReelID",
					goutils.Err(err))

				return err
			}

			i64, err := goutils.String2Int64(data)
			if err != nil {
				goutils.Error("LoadReelsStats2:LoadExcel:ondata:String2Int64",
					goutils.Err(err))

				return err
			}

			mapSymbols[curSymbol][rn-1] = int(i64)
		}

		return nil
	})
	if err != nil {
		goutils.Error("LoadReelsStats2:LoadExcel",
			goutils.Err(err))

		return nil, err
	}

	rss2 := NewReelsStats2(reelnum)

	for s, arr := range mapSymbols {
		rss2.Symbols = append(rss2.Symbols, s)
		for i, v := range arr {
			rss2.Reels[i].MapSymbols[s] = v
		}
	}

	return rss2, nil
}
