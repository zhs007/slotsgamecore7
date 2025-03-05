package lowcode

import "github.com/zhs007/goutils"

type PosComponentData struct {
	Pos []int
}

func (posdata *PosComponentData) MergePosList(pos []int) {
	for i := 0; i < len(pos)/2; i++ {
		posdata.Add(pos[i*2], pos[i*2+1])
	}
}

func (posdata *PosComponentData) Add(x, y int) {
	if goutils.IndexOfInt2Slice(posdata.Pos, x, y, 0) < 0 {
		posdata.Pos = append(posdata.Pos, x, y)
	}
}

func (posdata *PosComponentData) Has(target *PosComponentData) bool {
	return HasSamePos(posdata.Pos, target.Pos)
}

func (posdata *PosComponentData) ClearPos() {
	posdata.Pos = nil
}

func (posdata *PosComponentData) Clone() PosComponentData {
	target := PosComponentData{
		Pos: make([]int, len(posdata.Pos)),
	}

	copy(target.Pos, posdata.Pos)

	return target
}
