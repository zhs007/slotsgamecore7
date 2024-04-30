package lowcode

import "github.com/zhs007/goutils"

type PosComponentData struct {
	Pos []int
}

func (posdata *PosComponentData) Has(target *PosComponentData) bool {
	if len(posdata.Pos) == 0 || len(target.Pos) == 0 {
		return false
	}

	for i := 0; i < len(posdata.Pos)/2; i++ {
		if goutils.IndexOfInt2Slice(target.Pos, posdata.Pos[i*2], posdata.Pos[i*2+1], 0) >= 0 {
			return true
		}
	}

	return false
}

func (posdata *PosComponentData) Clear() {
	posdata.Pos = nil
}

func (posdata *PosComponentData) Clone() PosComponentData {
	target := PosComponentData{
		Pos: make([]int, len(posdata.Pos)),
	}

	copy(target.Pos, posdata.Pos)

	return target
}
