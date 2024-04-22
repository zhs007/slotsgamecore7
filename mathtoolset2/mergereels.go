package mathtoolset2

import (
	"io"

	"github.com/zhs007/goutils"
)

func mergeReels(src0 [][]string, src1 [][]string) [][]string {
	if src0 == nil {
		return src1
	}

	for x := range src0 {
		src0[x] = append(src0[x], src1[x]...)
	}

	return src0
}

func MergeReels(readers []io.Reader) ([][]string, error) {
	var trd [][]string

	for _, v := range readers {
		rd, err := LoadReels(v)
		if err != nil {
			goutils.Error("GenStackReels:LoadReelsStats2",
				goutils.Err(err))

			return nil, err
		}

		trd = mergeReels(trd, rd)
	}

	return trd, nil
}
