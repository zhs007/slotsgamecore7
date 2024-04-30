package lowcode

type PosComponentData struct {
	Pos []int
}

func (posdata *PosComponentData) Has(target *PosComponentData) bool {
	return HasSamePos(posdata.Pos, target.Pos)
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
