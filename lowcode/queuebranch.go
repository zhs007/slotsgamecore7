package lowcode

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const QueueBranchTypeName = "queueBranch"

// const (
// 	QBDVQueue string = "queue" // 队列数量
// )

type QueueBranchData struct {
	BasicComponentData
	Queue int
}

// OnNewGame -
func (queueBranchData *QueueBranchData) OnNewGame(gameProp *GameProperty, component IComponent) {
	queueBranchData.BasicComponentData.OnNewGame(gameProp, component)
}

// // OnNewStep -
// func (queueBranchData *QueueBranchData) OnNewStep(gameProp *GameProperty, component IComponent) {
// 	queueBranchData.BasicComponentData.OnNewStep(gameProp, component)
// }

// SetConfigIntVal -
func (queueBranchData *QueueBranchData) SetConfigIntVal(key string, val int) {
	if key == CCVQueue {
		queueBranchData.Queue = val
	}
}

// ChgConfigIntVal -
func (queueBranchData *QueueBranchData) ChgConfigIntVal(key string, off int) int {
	if key == CCVQueue {
		queueBranchData.Queue += off

		return queueBranchData.Queue
	}

	return queueBranchData.ChgConfigIntVal(key, off)
}

// Clone
func (queueBranchData *QueueBranchData) Clone() IComponentData {
	target := &QueueBranchData{
		BasicComponentData: queueBranchData.CloneBasicComponentData(),
		Queue:              queueBranchData.Queue,
	}

	return target
}

// BuildPBComponentData
func (queueBranchData *QueueBranchData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.QueueBranchData{
		BasicComponentData: queueBranchData.BuildPBBasicComponentData(),
		Queue:              int32(queueBranchData.Queue),
	}

	return pbcd
}

// GetValEx -
func (queueBranchData *QueueBranchData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVValue {
		return queueBranchData.Queue, true
	}

	return 0, false
}

// // SetVal -
// func (queueBranchData *QueueBranchData) SetVal(key string, val int) {
// 	if key == QBDVQueue {
// 		queueBranchData.Queue = val
// 	}
// }

// QueueBranchConfig - configuration for QueueBranch
type QueueBranchConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	JumpToComponent      string `yaml:"jumpToComponent" json:"jumpToComponent"`
}

// SetLinkComponent
func (cfg *QueueBranchConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
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
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &QueueBranchConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("QueueBranch.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

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
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// queueBranch.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	qbd := cd.(*QueueBranchData)

	if qbd.Queue > 0 {
		qbd.Queue--

		nc := queueBranch.onStepEnd(gameProp, curpr, gp, queueBranch.Config.JumpToComponent)

		return nc, nil
	}

	nc := queueBranch.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (queueBranch *QueueBranch) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	qbd := cd.(*QueueBranchData)

	fmt.Printf("queueBranch %v, got %v\n", queueBranch.GetName(), qbd.Queue)

	return nil
}

// // OnStats
// func (queueBranch *QueueBranch) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// NewComponentData -
func (queueBranch *QueueBranch) NewComponentData() IComponentData {
	return &QueueBranchData{}
}

// GetAllLinkComponents - get all link components
func (queueBranch *QueueBranch) GetAllLinkComponents() []string {
	return []string{queueBranch.Config.DefaultNextComponent, queueBranch.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (queueBranch *QueueBranch) GetNextLinkComponents() []string {
	return []string{queueBranch.Config.DefaultNextComponent, queueBranch.Config.JumpToComponent}
}

func NewQueueBranch(name string) IComponent {
	return &QueueBranch{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// "configuration": {},
type jsonQueueBranch struct {
}

func (jcfg *jsonQueueBranch) build() *QueueBranchConfig {
	cfg := &QueueBranchConfig{}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseQueueBranch(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	_, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseQueueBranch:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	// buf, err := cfg.MarshalJSON()
	// if err != nil {
	// 	goutils.Error("parseQueueBranch:MarshalJSON",
	// 		goutils.Err(err))

	// 	return "", err
	// }

	data := &jsonQueueBranch{}

	// err = sonic.Unmarshal(buf, data)
	// if err != nil {
	// 	goutils.Error("parseQueueBranch:Unmarshal",
	// 		goutils.Err(err))

	// 	return "", err
	// }

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: QueueBranchTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
