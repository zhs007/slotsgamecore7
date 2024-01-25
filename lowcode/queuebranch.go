package lowcode

import (
	"fmt"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const QueueBranchTypeName = "queueBranch"

const (
	QBDVQueue string = "queue" // 队列数量
)

type QueueBranchData struct {
	BasicComponentData
	Queue int
}

// OnNewGame -
func (queueBranchData *QueueBranchData) OnNewGame() {
	queueBranchData.BasicComponentData.OnNewGame()
}

// OnNewStep -
func (queueBranchData *QueueBranchData) OnNewStep() {
	queueBranchData.BasicComponentData.OnNewStep()
}

// BuildPBComponentData
func (queueBranchData *QueueBranchData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.QueueBranchData{
		BasicComponentData: queueBranchData.BuildPBBasicComponentData(),
		Queue:              int32(queueBranchData.Queue),
	}

	return pbcd
}

// GetVal -
func (queueBranchData *QueueBranchData) GetVal(key string) int {
	if key == QBDVQueue {
		return queueBranchData.Queue
	}

	return 0
}

// SetVal -
func (queueBranchData *QueueBranchData) SetVal(key string, val int) {
	if key == QBDVQueue {
		queueBranchData.Queue = val
	}
}

// QueueBranchConfig - configuration for QueueBranch
type QueueBranchConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	JumpToComponent      string `yaml:"jumpToComponent" json:"jumpToComponent"`
}

type QueueBranch struct {
	*BasicComponent `json:"-"`
	Config          *QueueBranchConfig `json:"config"`
}

// Init -
func (queueBranch *QueueBranch) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("QueueBranch.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &QueueBranchConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("QueueBranch.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return queueBranch.InitEx(cfg, pool)
}

// InitEx -
func (queueBranch *QueueBranch) InitEx(cfg any, pool *GamePropertyPool) error {
	queueBranch.Config = cfg.(*QueueBranchConfig)
	queueBranch.Config.ComponentType = QueueBranchTypeName

	queueBranch.onInit(&queueBranch.Config.BasicComponentConfig)

	return nil
}

// playgame
func (queueBranch *QueueBranch) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	queueBranch.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	qbd := gameProp.MapComponentData[queueBranch.Name].(*QueueBranchData)

	if qbd.Queue > 0 {
		qbd.Queue--

		queueBranch.onStepEnd(gameProp, curpr, gp, queueBranch.Config.JumpToComponent)
	} else {
		queueBranch.onStepEnd(gameProp, curpr, gp, "")
	}

	return nil
}

// OnAsciiGame - outpur to asciigame
func (queueBranch *QueueBranch) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	qbd := gameProp.MapComponentData[queueBranch.Name].(*QueueBranchData)

	fmt.Printf("queueBranch %v, got %v\n", queueBranch.GetName(), qbd.Queue)

	return nil
}

// OnStats
func (queueBranch *QueueBranch) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewQueueBranch(name string) IComponent {
	return &QueueBranch{
		BasicComponent: NewBasicComponent(name, 0),
	}
}
