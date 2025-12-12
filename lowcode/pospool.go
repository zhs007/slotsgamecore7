package lowcode

import (
	"slices"

	"github.com/zhs007/goutils"
)

type PosData struct {
	pos []int
}

func (pd *PosData) SetPos(pos []int) {
	pd.pos = pd.pos[:0]

	pd.pos = append(pd.pos, pos...)
}

func (pd *PosData) Len() int {
	return len(pd.pos) / 2
}

func (pd *PosData) IsEmpty() bool {
	return len(pd.pos) == 0
}

func (pd *PosData) Add(x, y int) {
	pd.pos = append(pd.pos, x, y)
}

func (pd *PosData) Has(x, y int) bool {
	return goutils.IndexOfInt2Slice(pd.pos, x, y, 0) >= 0
}

func (pd *PosData) Index(x, y int) int {
	return goutils.IndexOfInt2Slice(pd.pos, x, y, 0)
}

func (pd *PosData) Del(i int) {
	slices.Delete(pd.pos, i*2, i*2+2)
}

func (pd *PosData) Get(i int) (int, int) {
	return pd.pos[i*2], pd.pos[i*2+1]
}

type PosPool struct {
	pool []*PosData
	used []*PosData
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
		pd := &PosData{
			pos: make([]int, 0, pp.size),
		}

		pp.used = append(pp.used, pd)

		return pd
	}

	pd := pp.pool[len(pp.pool)-1]
	pp.pool = pp.pool[:len(pp.pool)-1]

	pp.used = append(pp.used, pd)

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

func (pp *PosPool) Reset() {
	for _, pd := range pp.used {
		pd.pos = pd.pos[:0]
	}

	pp.pool = append(pp.pool, pp.used...)
	pp.used = pp.used[:0]
}
