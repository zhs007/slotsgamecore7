package lowcode

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type FGDataCmdParam struct {
	FGNum                int    `json:"FGNum"`                // FG number
	RespinFirstComponent string `json:"respinFirstComponent"` // like fg-spin
}

// FGDataConfig - configuration for FG
type FGDataConfig struct {
	FGNum                     int            `yaml:"FGNum"`                     // FG number
	FGNumWeight               string         `yaml:"FGNumWeight"`               // FG number weight
	FGNumWithScatterNum       map[int]int    `yaml:"FGNumWithScatterNum"`       // FG number with scatter number
	FGNumWeightWithScatterNum map[int]string `yaml:"FGNumWeightWithScatterNum"` // FG number weight with scatter number
	RespinFirstComponent      string         `yaml:"respinFirstComponent"`      // like fg-spin
	Cmd                       string         `yaml:"cmd"`                       // cmd
}

// BasicWinsConfig - configuration for BasicWins
type MultiFGConfig struct {
	BasicComponentConfig `yaml:",inline"`
	FGData               []*FGDataConfig `yaml:"fgData"`          // wait player select
	TargetSymbolNum      string          `yaml:"targetSymbolNum"` // 这里可以用到一个前面记下的tagSymbolNum值
}

type MultiFG struct {
	*BasicComponent
	Config *MultiFGConfig
}

func (multiFG *MultiFG) parseCmdParam(cmd string, cmdParam string) (*FGDataCmdParam, error) {
	hascmd := false
	for _, v := range multiFG.Config.FGData {
		if v.Cmd == cmd {
			hascmd = true

			break
		}
	}

	if !hascmd {
		goutils.Error("MultiFG.parseCmdParam",
			zap.Error(ErrIvalidCmd))

		return nil, ErrIvalidCmd
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	param := &FGDataCmdParam{}
	err := json.Unmarshal([]byte(cmdParam), param)
	if err != nil {
		goutils.Error("MultiFG.parseCmdParam",
			zap.Error(ErrIvalidCmdParam))

		return nil, err
	}

	return param, nil
}

func (multiFG *MultiFG) genCmdParam(gameProp *GameProperty, fgdata *FGDataConfig, plugin sgc7plugin.IPlugin) (string, *FGDataCmdParam, error) {
	if fgdata.FGNumWeightWithScatterNum != nil {
		sn := gameProp.GetTagInt(multiFG.Config.TargetSymbolNum)

		vw2, err := gameProp.GetIntValWeights(fgdata.FGNumWeightWithScatterNum[sn])
		if err != nil {
			goutils.Error("MultiFG.genCmdParam:GetIntValWeights",
				zap.Int("symbolNum", sn),
				zap.Error(err))

			return "", nil, err
		}

		cr, err := vw2.RandVal(plugin)
		if err != nil {
			goutils.Error("MultiFG.genCmdParam:RandVal",
				zap.Error(err))

			return "", nil, err
		}

		return fgdata.Cmd, &FGDataCmdParam{
			FGNum:                cr.Int(),
			RespinFirstComponent: fgdata.RespinFirstComponent,
		}, nil
	} else if len(fgdata.FGNumWithScatterNum) > 0 {
		sn := gameProp.GetTagInt(multiFG.Config.TargetSymbolNum)

		return fgdata.Cmd, &FGDataCmdParam{
			FGNum:                fgdata.FGNumWithScatterNum[sn],
			RespinFirstComponent: fgdata.RespinFirstComponent,
		}, nil
	} else if fgdata.FGNumWeight != "" {
		vw2, err := gameProp.GetIntValWeights(fgdata.FGNumWeight)
		if err != nil {
			goutils.Error("MultiFG.genCmdParam:GetIntValWeights",
				zap.Error(err))

			return "", nil, err
		}

		cr, err := vw2.RandVal(plugin)
		if err != nil {
			goutils.Error("MultiFG.genCmdParam:RandVal",
				zap.Error(err))

			return "", nil, err
		}

		return fgdata.Cmd, &FGDataCmdParam{
			FGNum:                cr.Int(),
			RespinFirstComponent: fgdata.RespinFirstComponent,
		}, nil
	} else {
		return fgdata.Cmd, &FGDataCmdParam{
			FGNum:                fgdata.FGNum,
			RespinFirstComponent: fgdata.RespinFirstComponent,
		}, nil
	}
}

// Init -
func (multiFG *MultiFG) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MultiFG.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MultiFGConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MultiFG.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	multiFG.Config = cfg

	multiFG.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (multiFG *MultiFG) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[multiFG.Name].(*BasicComponentData)

	if cmd == DefaultCmd {
		lstcmd := []string{}
		lstparam := []string{}

		for _, v := range multiFG.Config.FGData {
			curcmd, curparam, err := multiFG.genCmdParam(gameProp, v, plugin)
			if err != nil {
				goutils.Error("MultiFG.OnPlayGame:genCmdParam",
					zap.Error(err))

				return err
			}

			json := jsoniter.ConfigCompatibleWithStandardLibrary

			buf, err := json.Marshal(curparam)
			if err != nil {
				goutils.Error("MultiFG.OnPlayGame:Marshal",
					zap.Error(err))

				return err
			}

			lstcmd = append(lstcmd, curcmd)
			lstparam = append(lstparam, string(buf))
		}

		curpr.NextCmds = lstcmd
		curpr.NextCmdParams = lstparam
		curpr.IsFinish = false
		curpr.IsWait = true
	} else {
		cmdparam, err := multiFG.parseCmdParam(cmd, param)
		if err != nil {
			goutils.Error("MultiFG.OnPlayGame:parseCmdParam",
				zap.String("cmd", cmd),
				zap.String("param", param),
				zap.Error(err))

			return err
		}

		gameProp.TriggerFG(curpr, gp, cmdparam.FGNum, cmdparam.RespinFirstComponent)
	}

	multiFG.onStepEnd(gameProp, curpr, gp, "")

	gp.AddComponentData(multiFG.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (multiFG *MultiFG) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	return nil
}

// OnStats
func (multiFG *MultiFG) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewMultiFG(name string) IComponent {
	multiFG := &MultiFG{
		BasicComponent: NewBasicComponent(name),
	}

	return multiFG
}
