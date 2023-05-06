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

type RespinDataCmdParam struct {
	RespinNum       int    `json:"RespinNum"`       // respin number
	RespinComponent string `json:"respinComponent"` // like fg-spin
}

// RespinDataConfig - configuration for MultiRespin
type RespinDataConfig struct {
	RespinNum                     int            `yaml:"respinNum"`                     // respin number
	RespinNumWeight               string         `yaml:"respinNumWeight"`               // respin number weight
	RespinNumWithScatterNum       map[int]int    `yaml:"respinNumWithScatterNum"`       // respin number with scatter number
	RespinNumWeightWithScatterNum map[int]string `yaml:"respinNumWeightWithScatterNum"` // respin number weight with scatter number
	RespinComponent               string         `yaml:"respinComponent"`               // like fg-spin
	Cmd                           string         `yaml:"cmd"`                           // cmd
}

// BasicWinsConfig - configuration for BasicWins
type MultiRespinConfig struct {
	BasicComponentConfig `yaml:",inline"`
	RespinData           []*RespinDataConfig `yaml:"respinData"`      // wait player select
	TargetSymbolNum      string              `yaml:"targetSymbolNum"` // 这里可以用到一个前面记下的tagSymbolNum值
}

type MultiRespin struct {
	*BasicComponent
	Config *MultiRespinConfig
}

func (multiRespin *MultiRespin) parseCmdParam(cmd string, cmdParam string) (*RespinDataCmdParam, error) {
	hascmd := false
	for _, v := range multiRespin.Config.RespinData {
		if v.Cmd == cmd {
			hascmd = true

			break
		}
	}

	if !hascmd {
		goutils.Error("MultiRespin.parseCmdParam",
			zap.Error(ErrIvalidCmd))

		return nil, ErrIvalidCmd
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	param := &RespinDataCmdParam{}
	err := json.Unmarshal([]byte(cmdParam), param)
	if err != nil {
		goutils.Error("MultiRespin.parseCmdParam",
			zap.Error(ErrIvalidCmdParam))

		return nil, err
	}

	return param, nil
}

func (multiRespin *MultiRespin) genCmdParam(gameProp *GameProperty, fgdata *RespinDataConfig, plugin sgc7plugin.IPlugin) (string, *RespinDataCmdParam, error) {
	if fgdata.RespinNumWeightWithScatterNum != nil {
		sn := gameProp.GetTagInt(multiRespin.Config.TargetSymbolNum)

		vw2, err := gameProp.GetIntValWeights(fgdata.RespinNumWeightWithScatterNum[sn])
		if err != nil {
			goutils.Error("MultiRespin.genCmdParam:GetIntValWeights",
				zap.Int("symbolNum", sn),
				zap.Error(err))

			return "", nil, err
		}

		cr, err := vw2.RandVal(plugin)
		if err != nil {
			goutils.Error("MultiRespin.genCmdParam:RandVal",
				zap.Error(err))

			return "", nil, err
		}

		return fgdata.Cmd, &RespinDataCmdParam{
			RespinNum:       cr.Int(),
			RespinComponent: fgdata.RespinComponent,
		}, nil
	} else if len(fgdata.RespinNumWithScatterNum) > 0 {
		sn := gameProp.GetTagInt(multiRespin.Config.TargetSymbolNum)

		return fgdata.Cmd, &RespinDataCmdParam{
			RespinNum:       fgdata.RespinNumWithScatterNum[sn],
			RespinComponent: fgdata.RespinComponent,
		}, nil
	} else if fgdata.RespinNumWeight != "" {
		vw2, err := gameProp.GetIntValWeights(fgdata.RespinNumWeight)
		if err != nil {
			goutils.Error("MultiRespin.genCmdParam:GetIntValWeights",
				zap.Error(err))

			return "", nil, err
		}

		cr, err := vw2.RandVal(plugin)
		if err != nil {
			goutils.Error("MultiFG.genCmdParam:RandVal",
				zap.Error(err))

			return "", nil, err
		}

		return fgdata.Cmd, &RespinDataCmdParam{
			RespinNum:       cr.Int(),
			RespinComponent: fgdata.RespinComponent,
		}, nil
	} else {
		return fgdata.Cmd, &RespinDataCmdParam{
			RespinNum:       fgdata.RespinNum,
			RespinComponent: fgdata.RespinComponent,
		}, nil
	}
}

// Init -
func (multiRespin *MultiRespin) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MultiRespin.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MultiRespinConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MultiRespin.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	multiRespin.Config = cfg

	multiRespin.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (multiRespin *MultiRespin) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	// cd := gameProp.MapComponentData[multiRespin.Name].(*BasicComponentData)

	if cmd == DefaultCmd {
		lstcmd := []string{}
		lstparam := []string{}

		for _, v := range multiRespin.Config.RespinData {
			curcmd, curparam, err := multiRespin.genCmdParam(gameProp, v, plugin)
			if err != nil {
				goutils.Error("MultiRespin.OnPlayGame:genCmdParam",
					zap.Error(err))

				return err
			}

			json := jsoniter.ConfigCompatibleWithStandardLibrary

			buf, err := json.Marshal(curparam)
			if err != nil {
				goutils.Error("MultiRespin.OnPlayGame:Marshal",
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
		cmdparam, err := multiRespin.parseCmdParam(cmd, param)
		if err != nil {
			goutils.Error("MultiFG.OnPlayGame:parseCmdParam",
				zap.String("cmd", cmd),
				zap.String("param", param),
				zap.Error(err))

			return err
		}

		gameProp.TriggerRespin(curpr, gp, cmdparam.RespinNum, cmdparam.RespinComponent)
	}

	multiRespin.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(multiRespin.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (multiRespin *MultiRespin) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	return nil
}

// OnStats
func (multiRespin *MultiRespin) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewMultiRespin(name string) IComponent {
	multiFG := &MultiRespin{
		BasicComponent: NewBasicComponent(name),
	}

	return multiFG
}
