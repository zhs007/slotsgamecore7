package lowcode

import (
	"github.com/zhs007/goutils"
)

type PosData struct {
	pos []int
}

func (pd *PosData) Add(x, y int) {
	pd.pos = append(pd.pos, x, y)
}

func (pd *PosData) Has(x, y int) bool {
	return goutils.IndexOfInt2Slice(pd.pos, x, y, 0) >= 0
}

type PosPool struct {
	pool []*PosData
	size int
}

func NewPosPool(size int) *PosPool {
	pp := &PosPool{
		pool: make([]*PosData, 128),
		size: size,
	}

	for i := 0; i < 128; i++ {
		pp.pool[i] = &PosData{
			pos: make([]int, 0, size),
		}
	}

	return pp
}

func (pp *PosPool) Get() *PosData {
	if len(pp.pool) == 0 {
		return &PosData{
			pos: make([]int, 0, pp.size),
		}
	}

	pd := pp.pool[len(pp.pool)-1]
	pp.pool = pp.pool[:len(pp.pool)-1]

	return pd
}

func (pp *PosPool) Clone(pd *PosData) *PosData {
	npd := pp.Get()

	npd.pos = append(npd.pos, pd.pos...)

	return npd
}

func (pp *PosPool) Put(pd *PosData) {
	if pd == nil {
		return
	}

	pd.pos = pd.pos[:0]

	pp.pool = append(pp.pool, pd)
}
