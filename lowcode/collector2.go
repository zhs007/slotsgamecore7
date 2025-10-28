package lowcode

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const Collector2TypeName = "collector2"

type Collector2PS struct {
	Value int `json:"value"` // value
}

// SetPublicJson
func (ps *Collector2PS) SetPublicJson(str string) error {
	err := sonic.UnmarshalString(str, ps)
	if err != nil {
		goutils.Error("Collector2PS.SetPublicJson:UnmarshalString",
			goutils.Err(err))

		return err
	}

	return nil
}

// SetPrivateJson
func (ps *Collector2PS) SetPrivateJson(str string) error {
	return nil
}

// GetPublicJson
func (ps *Collector2PS) GetPublicJson() string {
	str, err := sonic.MarshalString(ps)
	if err != nil {
		goutils.Error("Collector2PS.GetPublicJson:MarshalString",
			goutils.Err(err))

		return ""
	}

	return str
}

// GetPrivateJson
func (ps *Collector2PS) GetPrivateJson() string {
	return ""
}

// Clone
func (ps *Collector2PS) Clone() IComponentPS {
	return &Collector2PS{
		Value: ps.Value,
	}
}

type Collector2Data struct {
	BasicComponentData
	Val          int // 当前总值, Current total value
	NewCollector int // 这一个step收集到的, The values collected in this step
	cfg          *Collector2Config
}

func (collectorData *Collector2Data) OnNewGame(gameProp *GameProperty, component IComponent) {
	collectorData.BasicComponentData.OnNewGame(gameProp, component)
	collectorData.Val = 0
	collectorData.NewCollector = 0
}

func (collectorData *Collector2Data) onNewStep() {
	collectorData.NewCollector = 0
}

func (collectorData *Collector2Data) SetConfigIntVal(key string, val int) {
	if key == CCVValueNum {
		collectorData.Val = val
	} else {
		collectorData.BasicComponentData.SetConfigIntVal(key, val)
	}
}

// ChgConfigIntVal -
func (collectorData *Collector2Data) ChgConfigIntVal(key string, off int) int {
	if key == CCVValueNumNow {
		val := collectorData.Val + off

		if val < 0 {
			val = 0
		}

		if val > collectorData.cfg.MaxVal && collectorData.cfg.MaxVal > 0 {
			val = collectorData.cfg.MaxVal

			collectorData.NewCollector = collectorData.cfg.MaxVal - collectorData.Val
		} else {
			collectorData.NewCollector = off
		}

		collectorData.BasicComponentData.SetConfigIntVal(key, val)

		return val
	} else {
		return collectorData.BasicComponentData.ChgConfigIntVal(key, off)
	}
}

func (collectorData *Collector2Data) GetOutput() int {
	return collectorData.Val
}

func (collectorData *Collector2Data) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVValue || key == CCValueNum {
		return collectorData.Val, true
	}

	return 0, false
}

func (collectorData *Collector2Data) Clone() IComponentData {
	target := &Collector2Data{
		BasicComponentData: collectorData.CloneBasicComponentData(),
		Val:                collectorData.Val,
		NewCollector:       collectorData.NewCollector,
	}
	return target
}

func (collectorData *Collector2Data) BuildPBComponentData() proto.Message {
	return &sgc7pb.CollectorData{
		Val:          int32(collectorData.Val),
		NewCollector: int32(collectorData.NewCollector),
	}
}

type Collector2Config struct {
	BasicComponentConfig     `yaml:",inline" json:",inline"`
	MaxVal                   int                 `yaml:"maxVal" json:"maxVal"`
	IsCycle                  bool                `yaml:"isCycle" json:"isCycle"`
	IsPlayerState            bool                `yaml:"isPlayerState" json:"isPlayerState"`
	IsIgnoreBet              bool                `yaml:"isIgnoreBet" json:"isIgnoreBet"`
	IsForceTriggerController bool                `yaml:"isForceTriggerController" json:"isForceTriggerController"`
	MapAwards                map[string][]*Award `yaml:"controllers" json:"controllers"`
}

func (cfg *Collector2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type Collector2 struct {
	*BasicComponent `json:"-"`
	Config          *Collector2Config `json:"config"`
}

func (collector *Collector2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Collector2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))
		return err
	}

	cfg := &Collector2Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Collector2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))
		return err
	}

	return collector.InitEx(cfg, pool)
}

func (collector *Collector2) InitEx(cfg any, pool *GamePropertyPool) error {
	cfg2, isok := cfg.(*Collector2Config)
	if !isok {
		goutils.Error("Collector2.InitEx:cfg type",
			goutils.Err(ErrInvalidComponentConfig))
		return ErrInvalidComponentConfig
	}
	collector.Config = cfg2
	collector.Config.ComponentType = Collector2TypeName

	if collector.Config.MapAwards != nil {
		for _, lst := range collector.Config.MapAwards {
			for _, v := range lst {
				v.Init()
			}
		}
	}

	collector.onInit(&collector.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (collector *Collector2) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	awards, isok := collector.Config.MapAwards[strVal]
	if isok {
		gameProp.procAwards(plugin, awards, curpr, gp)
	}
}

func (collector *Collector2) onAdd(plugin sgc7plugin.IPlugin, startVal int, num int, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams) {
	if num > 0 {
		for i := 1; i <= num; i++ {
			cl := startVal + i
			if cl >= collector.Config.MaxVal {
				collector.onLevelUp(plugin, gameProp, curpr, gp, -1, false)
				break
			} else {
				collector.onLevelUp(plugin, gameProp, curpr, gp, cl, false)
			}
		}
	}
}

func (collector *Collector2) add(plugin sgc7plugin.IPlugin, num int, cd *Collector2Data, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, noProcLevelUp bool) error {
	if num <= 0 {
		return nil
	}
	if collector.Config.MaxVal > 0 {
		if collector.Config.IsCycle {
			if cd.Val+num >= collector.Config.MaxVal {
				cd.NewCollector += num
				oldval := cd.Val
				cd.Val += num
				for {
					if cd.Val >= collector.Config.MaxVal {
						if !noProcLevelUp {
							collector.onAdd(plugin, oldval, collector.Config.MaxVal-oldval, gameProp, curpr, gp)
						}
						if cd.Val == collector.Config.MaxVal {
							cd.Val = 0
							break
						}
						oldval = 0
						cd.Val -= collector.Config.MaxVal
					} else {
						if !noProcLevelUp {
							collector.onAdd(plugin, oldval, cd.Val-oldval, gameProp, curpr, gp)
						}
						break
					}
				}
				return nil
			}
		} else {
			if cd.Val == collector.Config.MaxVal {
				return nil
			}
			if cd.Val > collector.Config.MaxVal {
				goutils.Error("Collector2.add",
					goutils.Err(ErrInvalidCollectorVal))
				return ErrInvalidCollectorVal
			}
			if cd.Val+num >= collector.Config.MaxVal {
				oldval := cd.Val
				cd.NewCollector += num
				cd.Val = collector.Config.MaxVal
				if !noProcLevelUp {
					collector.onAdd(plugin, oldval, cd.Val-oldval, gameProp, curpr, gp)
				}
				return nil
			}
		}
	}
	cd.NewCollector += num
	oldval := cd.Val
	cd.Val += num
	if collector.Config.MaxVal > 0 && cd.Val >= collector.Config.MaxVal {
		goutils.Error("Collector2.add",
			goutils.Err(ErrInvalidCollectorLogic))
		return ErrInvalidCollectorLogic
	}
	if num > 0 && !noProcLevelUp && oldval != cd.Val {
		collector.onAdd(plugin, oldval, num, gameProp, curpr, gp)
	}
	return nil
}

func (collector *Collector2) onLevelUp(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, newLevel int, noProcSPLevel bool) error {
	collector.ProcControllers(gameProp, plugin, curpr, gp, newLevel, "<trigger>")

	if noProcSPLevel {
		return nil
	}

	collector.ProcControllers(gameProp, plugin, curpr, gp, newLevel, fmt.Sprintf("%d", newLevel))

	return nil
}

func (collector *Collector2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ips sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {
	ccd := cd.(*Collector2Data)
	ccd.onNewStep()

	if collector.Config.IsPlayerState {
		ps, isok := ips.(*PlayerState)
		if !isok {
			goutils.Error("Collector2.OnPlayGame:PlayerState",
				goutils.Err(ErrInvalidPlayerState))

			return "", ErrInvalidPlayerState
		}

		betMethod := stake.CashBet / stake.CoinBet
		bmd := ps.GetBetMethodPub(int(betMethod))
		if bmd == nil {
			goutils.Error("Collector2.OnPlayGame:GetBetMethodPub",
				goutils.Err(ErrInvalidPlayerState))

			return "", ErrInvalidPlayerState
		}

		bet := stake.CoinBet
		if collector.Config.IsIgnoreBet {
			bet = -1
		}

		cps := bmd.GetBetCPS(int(bet), collector.GetName())
		if cps == nil {
			goutils.Error("Collector2.OnPlayGame:GetBetCPS",
				goutils.Err(ErrInvalidPlayerState))

			return "", ErrInvalidPlayerState
		}

		cbps, isok := cps.(*Collector2PS)
		if !isok {
			goutils.Error("Collector2.OnPlayGame:Collector2PS",
				goutils.Err(ErrInvalidPlayerState))

			return "", ErrInvalidPlayerState
		}

		if collector.Config.IsForceTriggerController {
			for ci := cbps.Value; ci >= 0; ci-- {
				strCurVal := fmt.Sprintf("%d", ci)
				_, isok := collector.Config.MapAwards[strCurVal]
				if isok {
					collector.ProcControllers(gameProp, plugin, curpr, gp, ci, strCurVal)

					break
				}
			}
		}

		ccd.Val = cbps.Value

		off, isok := ccd.GetConfigIntVal(CCVValueNum)
		if isok {
			err := collector.add(plugin, off, ccd, gameProp, curpr, gp, false)
			if err != nil {
				goutils.Error("Collector2.OnPlayGame:add:off",
					goutils.Err(err))

				return "", err
			}

			ccd.ClearConfigIntVal(CCVValueNum)
		}

		cbps.Value = ccd.Val

		nc := collector.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	off, isok := ccd.GetConfigIntVal(CCVValueNum)
	if isok {
		err := collector.add(plugin, off, ccd, gameProp, curpr, gp, false)
		if err != nil {
			goutils.Error("Collector2.OnPlayGame:add:off",
				goutils.Err(err))
			return "", err
		}
		ccd.ClearConfigIntVal(CCVValueNum)
	}

	nc := collector.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

func (collector *Collector2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	ccd, isok := cd.(*Collector2Data)
	if !isok {
		goutils.Error("Collector2.OnAsciiGame:Collector2Data",
			goutils.Err(ErrInvalidComponentData))
		return ErrInvalidComponentData
	}

	if ccd.NewCollector <= 0 {
		fmt.Printf("%v does not collect new value, the collector value is %v\n", collector.Name, ccd.Val)
	} else {
		fmt.Printf("%v collect %v. the collector value is %v\n", collector.Name, ccd.NewCollector, ccd.Val)
	}
	return nil
}

func (collector *Collector2) NewComponentData() IComponentData {
	return &Collector2Data{
		cfg: collector.Config,
	}
}

func (collector *Collector2) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

func (collector *Collector2) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	collector.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)
	if isOnStepEnd && pr.IsFinish {
		cd := icd.(*Collector2Data)
		s2.ProcStatsIntVal(collector.GetName(), cd.Val)
	}
}

func (collector *Collector2) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, []stats2.Option{stats2.OptIntVal})
}

func (collector *Collector2) IsNeedOnStepEndStats2() bool {
	return true
}

// InitPlayerState -
func (collector *Collector2) InitPlayerState(pool *GamePropertyPool, gameProp *GameProperty, plugin sgc7plugin.IPlugin,
	ps *PlayerState, betMethod int, bet int) error {

	if collector.Config.IsPlayerState {
		bmd := ps.GetBetMethodPub(betMethod)
		if bet <= 0 {
			return nil
		}

		// 如果忽略下注，这时只处理 bet 为 -1 的情况
		if collector.Config.IsIgnoreBet {
			bet = -1
		}

		bps := bmd.GetBetPS(bet)

		cname := collector.GetName()

		_, isok := bps.MapComponentData[cname]
		if !isok {
			str, isok := bps.MapString[cname]
			if isok {
				cps := &Collector2PS{}
				cps.SetPublicJson(str)

				bps.MapComponentData[cname] = cps
			} else {
				cps := &Collector2PS{
					Value: 0,
				}

				bps.MapComponentData[cname] = cps
			}
		}
	}

	return nil
}

// OnUpdateDataWithPlayerState -
func (collector *Collector2) OnUpdateDataWithPlayerState(pool *GamePropertyPool, gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, ps *PlayerState, betMethod int, bet int, cd IComponentData) {
	if collector.Config.IsPlayerState {
		bmd := ps.GetBetMethodPub(betMethod)
		if bet <= 0 {
			return
		}

		// 如果忽略下注，这时只处理 bet 为 -1 的情况
		if collector.Config.IsIgnoreBet {
			bet = -1
		}

		bps := bmd.GetBetPS(bet)

		cname := collector.GetName()

		ips, isok := bps.MapComponentData[cname]
		if !isok {
			goutils.Error("Collector2.OnUpdateDataWithPlayerState:MapComponentData",
				goutils.Err(ErrInvalidPlayerState))

			return
		}

		cps, isok := ips.(*Collector2PS)
		if !isok {
			goutils.Error("Collector2.OnUpdateDataWithPlayerState:Collector2PS",
				goutils.Err(ErrInvalidPlayerState))

			return
		}

		cd2, isok := cd.(*Collector2Data)
		if !isok {
			goutils.Error("Collector2.OnUpdateDataWithPlayerState:Collector2Data",
				goutils.Err(ErrInvalidComponentData))

			return
		}

		// CCVValueNumNow 在这里非常特殊,理论上,这里(playerstate时)什么都不用做,只需要缓存数据即可
		val, isok := cd2.GetConfigIntVal(CCVValueNumNow)
		if isok {
			if val < 0 {
				val = 0
			}

			if val > collector.Config.MaxVal && collector.Config.MaxVal > 0 {
				val = collector.Config.MaxVal
			}

			// 这里不触发onLevelUp,因为数据是从playerstate里同步过来的,实时写回去时只要把值设置正确就行了
			cd2.Val = val

			cd2.ClearConfigIntVal(CCVValueNumNow)

			cps.Value = val
		}
	}
}

func NewCollector2(name string) IComponent {
	collector := &Collector2{
		BasicComponent: NewBasicComponent(name, 1),
	}
	return collector
}

// "maxVal": 20,
// "isCycle": false,
// "isPlayerState": true,
// "isIgnoreBet": true,
// "isForceTriggerController": true

type jsonCollector2 struct {
	MaxVal                   int  `json:"maxVal"`
	IsCycle                  bool `json:"isCycle"`
	IsPlayerState            bool `json:"isPlayerState"`
	IsIgnoreBet              bool `json:"isIgnoreBet"`
	IsForceTriggerController bool `json:"isForceTriggerController"`
}

func (jcfg *jsonCollector2) build() *Collector2Config {
	cfg := &Collector2Config{
		MaxVal:                   jcfg.MaxVal,
		IsCycle:                  jcfg.IsCycle,
		IsPlayerState:            jcfg.IsPlayerState,
		IsIgnoreBet:              jcfg.IsIgnoreBet,
		IsForceTriggerController: jcfg.IsForceTriggerController,
	}
	return cfg
}

func parseCollector2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseCollector2:getConfigInCell",
			goutils.Err(err))
		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseCollector2:MarshalJSON",
			goutils.Err(err))
		return "", err
	}

	data := &jsonCollector2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseCollector2:Unmarshal",
			goutils.Err(err))
		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseSymbolValsSP:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: Collector2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
