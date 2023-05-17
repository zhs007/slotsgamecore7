package sgc7game

import (
	"sync"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

var InitGameScenePoolSize int

type gameScenePool struct {
	Pools []*GameScene
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

func newGameScenePool() *gameScenePool {
	return &gameScenePool{
		Pools: make([]*GameScene, 0, InitGameScenePoolSize),
	}
}

type GameScenePoolEx struct {
	Lock     sync.Mutex
	MapPools map[int]map[int]*gameScenePool
}

func (pool *GameScenePoolEx) new(w, h int) *GameScene {
	s, err := NewGameScene(w, h)
	if err != nil {
		goutils.Error("GameScenePoolEx.new:NewGameScene",
			zap.Error(err))

		return nil
	}

	return s
}

func (pool *GameScenePoolEx) Put(scene *GameScene) {
	pool.Lock.Lock()
	defer pool.Lock.Unlock()

	pool.MapPools[scene.Width][scene.Height].put(scene)
}

func (pool *GameScenePoolEx) New(w, h int, isNeedClear bool) *GameScene {
	pool.Lock.Lock()
	defer pool.Lock.Unlock()

	mps, isok := pool.MapPools[w]
	if !isok {
		mps = make(map[int]*gameScenePool)
		pool.MapPools[w] = mps
	}

	p, isok := mps[h]
	if isok {
		gs := p.get()

		if gs != nil {
			if isNeedClear {
				gs.Clear(0)
			}

			return gs
		}
	} else {
		p = newGameScenePool()

		mps[h] = p
	}

	return pool.new(w, h)
}

func NewGameScenePoolEx() *GameScenePoolEx {
	return &GameScenePoolEx{
		MapPools: make(map[int]map[int]*gameScenePool),
	}
}

func init() {
	InitGameScenePoolSize = 16
}