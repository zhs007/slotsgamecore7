package mathtoolset2

import (
	"io"
	"strings"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type ReelStats2 struct {
	MapSymbols     map[string]int
	TotalSymbolNum int
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
				zap.Error(err))

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
				zap.Error(err))
		}

		if rn > reelnum {
			reelnum = rn
		}

		return head
	}, func(x int, y int, header string, data string) error {
		if reelnum <= 0 {
			goutils.Error("LoadReelsStats2:LoadExcel:ondata:reelnum",
				zap.Int("reelnum", reelnum),
				zap.Error(ErrInvalidReelsStats2File))

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
					zap.Error(err))

				return err
			}

			i64, err := goutils.String2Int64(data)
			if err != nil {
				goutils.Error("LoadReelsStats2:LoadExcel:ondata:String2Int64",
					zap.Error(err))

				return err
			}

			mapSymbols[curSymbol][rn-1] = int(i64)
		}

		return nil
	})
	if err != nil {
		goutils.Error("LoadReelsStats2:LoadExcel",
			zap.Error(err))

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
