package sgc7game

import (
	"github.com/zhs007/goutils"
)

var InitGameScenePoolSize int

type gameScenePool struct {
	Pools []*GameScene
	used  []*GameScene
}

func (pool *gameScenePool) put(gs *GameScene) {
	pool.Pools = append(pool.Pools, gs)
}

func (pool *gameScenePool) get() *GameScene {
	if len(pool.Pools) > 0 {
		gs := pool.Pools[len(pool.Pools)-1]

		pool.Pools = pool.Pools[0 : len(pool.Pools)-1]

		return gs
	}

	return nil
}

func (pool *gameScenePool) inUsed(sc *GameScene) {
	pool.used = append(pool.used, sc)
}

func (pool *gameScenePool) reset() {
	pool.Pools = append(pool.Pools, pool.used...)
	pool.used = nil
}

func newGameScenePool() *gameScenePool {
	return &gameScenePool{
		Pools: make([]*GameScene, 0, InitGameScenePoolSize),
	}
}

type GameScenePoolEx struct {
	MapPools map[int]map[int]*gameScenePool
}

func (pool *GameScenePoolEx) new(w, h int) *GameScene {
	s, err := NewGameScene(w, h)
	if err != nil {
		goutils.Error("GameScenePoolEx.new:NewGameScene",
			goutils.Err(err))

		return nil
	}

	return s
}

func (pool *GameScenePoolEx) new2(w, h int, v int) *GameScene {
	s, err := NewGameScene2(w, h, v)
	if err != nil {
		goutils.Error("GameScenePoolEx.new2:NewGameScene2",
			goutils.Err(err))

		return nil
	}

	return s
}

func (pool *GameScenePoolEx) Put(scene *GameScene) {
	pool.MapPools[scene.Width][scene.Height].put(scene)
}

func (pool *GameScenePoolEx) Reset() {
	for _, mps := range pool.MapPools {
		for _, p := range mps {
			p.reset()
		}
	}
}

func (pool *GameScenePoolEx) New(w, h int) *GameScene {
	mps, isok := pool.MapPools[w]
	if !isok {
		mps = make(map[int]*gameScenePool)
		pool.MapPools[w] = mps
	}

	p, isok := mps[h]
	if isok {
		gs := p.get()

		if gs != nil {
			p.inUsed(gs)

			return gs
		}
	} else {
		p = newGameScenePool()

		mps[h] = p
	}

	gs := pool.new2(w, h, 0)
	p.inUsed(gs)
	return gs
}

func (pool *GameScenePoolEx) New2(w, h int, v int) *GameScene {
	mps, isok := pool.MapPools[w]
	if !isok {
		mps = make(map[int]*gameScenePool)
		pool.MapPools[w] = mps
	}

	p, isok := mps[h]
	if isok {
		gs := p.get()

		if gs != nil {
			gs.Clear(v)

			p.inUsed(gs)

			return gs
		}
	} else {
		p = newGameScenePool()

		mps[h] = p
	}

	gs := pool.new2(w, h, v)
	p.inUsed(gs)
	return gs
}

func NewGameScenePoolEx() *GameScenePoolEx {
	return &GameScenePoolEx{
		MapPools: make(map[int]map[int]*gameScenePool),
	}
}

func init() {
	InitGameScenePoolSize = 16
}
