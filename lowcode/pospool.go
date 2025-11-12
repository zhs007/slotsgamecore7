package lowcode

import (
	"sync"

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
	pool sync.Pool
	size int
}

func NewPosPool(size int) *PosPool {
	return &PosPool{
		pool: sync.Pool{
			New: func() any {
				return &PosData{
					pos: make([]int, 0, size),
				}
			},
		},
		size: size,
	}
}

func (pp *PosPool) Get() *PosData {
	return pp.pool.Get().(*PosData)
}

func (pp *PosPool) Clone(pd *PosData) *PosData {
	npd := pp.pool.Get().(*PosData)

	npd.pos = append(npd.pos, pd.pos...)

	return npd
}

func (pp *PosPool) Put(pd *PosData) {
	pd.pos = pd.pos[:0]

	pp.pool.Put(pd)
}
