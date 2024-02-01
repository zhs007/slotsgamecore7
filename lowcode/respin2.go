package lowcode

// const Respin2TypeName = "respin2"

// // // 新逻辑如下：
// // // 1. respin有2个可配置逻辑，分别是trigger和retrigger
// // // 2. award 可以修改 RespinNum 和 RetriggerAddRespinNum
// // // 3. award 可以 TriggerRespin2 ，这种会分别执行trigger和retrigger的逻辑

// // // 第一次触发的逻辑
// // type RespinTriggerType int

// // const (
// // 	RTTNone       RespinTriggerType = 0 // 默认值，不需要执行trigger逻辑
// // 	RTTUseTrigger RespinTriggerType = 1 // 需要执行trigger逻辑
// // )

// // func Str2RespinTriggerType(str string) RespinTriggerType {
// // 	if str == "useTrigger" {
// // 		return RTTUseTrigger
// // 	}

// // 	return RTTNone
// // }

// // // 再触发的逻辑
// // type RespinRetriggerType int

// // const (
// // 	RRTNone                     RespinRetriggerType = 0 // 默认值，没有再次触发逻辑，其它组件直接加 RespinNum 即可
// // 	RRTUseRetriggerAddRespinNum RespinRetriggerType = 1 // RetriggerAddRespinNum，这个值会缓存下来，可以修改，等到retrigger时再用这个值加次数
// // 	RRTUseTrigger               RespinRetriggerType = 2 // 用trigger的逻辑，但用retrigger时的scatter数量执行这个逻辑
// // 	RRTUseRetrigger             RespinRetriggerType = 3 // 用retrigger的逻辑
// // )

// // func Str2RespinRetriggerType(str string) RespinRetriggerType {
// // 	if str == "useTrigger" {
// // 		return RRTUseTrigger
// // 	} else if str == "useRetrigger" {
// // 		return RRTUseRetrigger
// // 	} else if str == "useRetriggerAddRespinNum" {
// // 		return RRTUseRetriggerAddRespinNum
// // 	}

// // 	return RRTNone
// // }

// type Respin2Data struct {
// 	BasicComponentData
// 	LastRespinNum         int
// 	CurRespinNum          int
// 	CurAddRespinNum       int
// 	RetriggerAddRespinNum int      // 再次触发时增加的次数
// 	TotalCoinWin          int64    //
// 	TotalCashWin          int64    //
// 	LastTriggerNum        int      // 剩余的触发次数，respin有2种模式，一种是直接增加免费次数，一种是累积整体触发次数
// 	CurTriggerNum         int      // 当前已经触发次数
// 	Awards                []*Award // 当前已经触发次数
// 	TriggerRespinNum      []int    // 配合LastTriggerNum用的respin次数，-1表示用当前的RetriggerAddRespinNum，否则就是具体值
// }

// // OnNewGame -
// func (respin2Data *Respin2Data) OnNewGame() {
// 	respin2Data.BasicComponentData.OnNewGame()

// 	respin2Data.LastRespinNum = 0
// 	respin2Data.CurRespinNum = 0
// 	respin2Data.CurAddRespinNum = 0
// 	respin2Data.TotalCoinWin = 0
// 	respin2Data.TotalCashWin = 0
// 	respin2Data.RetriggerAddRespinNum = 0
// 	respin2Data.LastTriggerNum = 0
// 	respin2Data.CurTriggerNum = 0
// 	respin2Data.Awards = nil
// }

// // OnNewStep -
// func (respin2Data *Respin2Data) OnNewStep() {
// 	respin2Data.BasicComponentData.OnNewStep()

// 	respin2Data.CurAddRespinNum = 0
// }

// // BuildPBComponentData
// func (respin2Data *Respin2Data) BuildPBComponentData() proto.Message {
// 	pbcd := &sgc7pb.Respin2Data{
// 		BasicComponentData:    respin2Data.BuildPBBasicComponentData(),
// 		LastRespinNum:         int32(respin2Data.LastRespinNum),
// 		CurRespinNum:          int32(respin2Data.CurRespinNum),
// 		CurAddRespinNum:       int32(respin2Data.CurAddRespinNum),
// 		TotalCoinWin:          respin2Data.TotalCoinWin,
// 		TotalCashWin:          respin2Data.TotalCashWin,
// 		RetriggerAddRespinNum: int32(respin2Data.RetriggerAddRespinNum),
// 		LastTriggerNum:        int32(respin2Data.LastTriggerNum),
// 		CurTriggerNum:         int32(respin2Data.CurTriggerNum),
// 	}

// 	return pbcd
// }

// // Respin2LevelConfig - configuration for Respin Level
// type Respin2LevelConfig struct {
// 	LastRespinNum int    `yaml:"lastRespinNum" json:"lastRespinNum"` // 倒数第几局开始
// 	MaxCoinWins   int    `yaml:"maxCoinWins" json:"maxCoinWins"`     // 如果最大获奖低于这个
// 	JumpComponent string `yaml:"jumpComponent" json:"jumpComponent"` // 跳转到这个component
// }

// // Respin2TriggerConfig - configuration for TriggerRespin
// type Respin2TriggerConfig struct {
// 	RespinNum                       int                           `yaml:"respinNum" json:"respinNum"`                                         // 固定次数
// 	RespinNumWeight                 string                        `yaml:"respinNumWeight" json:"respinNumWeight"`                             // respin number weight
// 	RespinNumWithScatterNum         map[int]int                   `yaml:"respinNumWithScatterNum" json:"respinNumWithScatterNum"`             // respin number with scatter number
// 	RespinNumWeightWithScatterNum   map[int]string                `yaml:"respinNumWeightWithScatterNum" json:"respinNumWeightWithScatterNum"` // respin number weight with scatter number
// 	RespinNumWeightWithScatterNumVW map[int]*sgc7game.ValWeights2 `yaml:"-" json:"-"`                                                         // respin number weight with scatter number
// }

// // Respin2Config - configuration for Respin2
// type Respin2Config struct {
// 	BasicComponentConfig `yaml:",inline" json:",inline"`
// 	MainComponent        string                `yaml:"mainComponent" json:"mainComponent"`
// 	IsWinBreak           bool                  `yaml:"isWinBreak" json:"isWinBreak"`
// 	Levels               []*Respin2LevelConfig `yaml:"levels" json:"levels"`
// 	// TriggerTypeStr        string                `yaml:"triggerType" json:"triggerType"`                     // 触发逻辑如何执行
// 	// TriggerType           RespinTriggerType     `yaml:"-" json:"-"`                                         //
// 	// Trigger               *Respin2TriggerConfig `yaml:"trigger" json:"trigger"`                             // 触发逻辑
// 	// RetriggerTypeStr      string                `yaml:"retriggerType" json:"retriggerType"`                 // 再次触发逻辑如何执行
// 	// RetriggerType         RespinRetriggerType   `yaml:"-" json:"-"`                                         //
// 	// Retrigger             *Respin2TriggerConfig `yaml:"retrigger" json:"retrigger"`                         // 再次触发逻辑
// 	// RetriggerAddRespinNum int                   `yaml:"retriggerAddRespinNum" json:"retriggerAddRespinNum"` // RetriggerAddRespinNum的初始值
// }

// type Respin2 struct {
// 	*BasicComponent `json:"-"`
// 	Config          *Respin2Config `json:"config"`
// }

// // // OnNewGame -
// // func (respin *Respin) OnNewGame(gameProp *GameProperty) error {
// // 	cd := gameProp.MapComponentData[respin.Name]

// // 	cd.OnNewGame()

// // 	return nil
// // }

// // OnPlayGame - on playgame
// func (respin2 *Respin2) procLevel(level *Respin2LevelConfig, respin2Data *Respin2Data, gameProp *GameProperty) bool {
// 	if respin2Data.LastRespinNum <= level.LastRespinNum && respin2Data.CoinWin < level.MaxCoinWins {
// 		return true
// 	}

// 	return false
// }

// // OnPlayGame - on playgame
// func (respin2 *Respin2) AddRespinTimes(gameProp *GameProperty, num int) {
// 	cd := gameProp.MapComponentData[respin2.Name].(*RespinData)

// 	cd.LastRespinNum += num
// 	cd.CurAddRespinNum += num
// }

// // Init -
// func (respin2 *Respin2) Init(fn string, pool *GamePropertyPool) error {
// 	data, err := os.ReadFile(fn)
// 	if err != nil {
// 		goutils.Error("Respin2.Init:ReadFile",
// 			zap.String("fn", fn),
// 			zap.Error(err))

// 		return err
// 	}

// 	cfg := &Respin2Config{}

// 	err = yaml.Unmarshal(data, cfg)
// 	if err != nil {
// 		goutils.Error("Respin2.Init:Unmarshal",
// 			zap.String("fn", fn),
// 			zap.Error(err))

// 		return err
// 	}

// 	return respin2.InitEx(cfg, pool)
// }

// // InitEx -
// func (respin2 *Respin2) InitEx(cfg any, pool *GamePropertyPool) error {
// 	respin2.Config = cfg.(*Respin2Config)
// 	respin2.Config.ComponentType = RespinTypeName

// 	respin2.onInit(&respin2.Config.BasicComponentConfig)

// 	return nil
// }

// // playgame
// func (respin2 *Respin2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
// 	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

// 	respin2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

// 	cd := gameProp.MapComponentData[respin2.Name].(*Respin2Data)

// 	// if cd.CurRespinNum == 0 && cd.LastRespinNum == 0 && respin2.Config.InitRespinNum > 0 {
// 	// 	cd.LastRespinNum = respin2.Config.InitRespinNum
// 	// }

// recheck:
// 	if cd.LastRespinNum == 0 {
// 		if cd.LastTriggerNum > 0 {
// 			respin2.Trigger(gameProp, plugin, curpr, gp)

// 			goto recheck
// 		}

// 		respin2.onStepEnd(gameProp, curpr, gp, respin2.Config.DefaultNextComponent)
// 	} else {
// 		nextComponent := respin2.Config.MainComponent

// 		for _, v := range respin2.Config.Levels {
// 			if respin2.procLevel(v, cd, gameProp) {
// 				nextComponent = v.JumpComponent

// 				break
// 			}
// 		}

// 		if cd.LastRespinNum > 0 {
// 			cd.LastRespinNum--
// 		}

// 		cd.CurRespinNum++

// 		respin2.onStepEnd(gameProp, curpr, gp, nextComponent)
// 	}

// 	// gp.AddComponentData(respin.Name, cd)

// 	return nil
// }

// // OnAsciiGame - outpur to asciigame
// func (respin2 *Respin2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
// 	cd := gameProp.MapComponentData[respin2.Name].(*Respin2Data)

// 	if cd.CurAddRespinNum > 0 {
// 		fmt.Printf("%v last %v, current %v, retrigger %v\n", respin2.Name, cd.LastRespinNum, cd.CurRespinNum, cd.CurAddRespinNum)
// 	} else {
// 		fmt.Printf("%v last %v, current %v\n", respin2.Name, cd.LastRespinNum, cd.CurRespinNum)
// 	}

// 	return nil
// }

// // OnStats
// func (respin2 *Respin2) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	if feature != nil && len(lst) > 0 {

// 		if feature.RespinNumStatus != nil ||
// 			feature.RespinWinStatus != nil {
// 			pbcd, lastpr := findLastPBComponentData(lst, respin2.Name)
// 			if pbcd != nil {
// 				respin2.onStatsWithPBEnding(feature, pbcd, lastpr)
// 			}
// 		}

// 		if feature.RespinStartNumStatus != nil {
// 			pbcd, firstpr := findFirstPBComponentData(lst, respin2.Name)
// 			if pbcd != nil {
// 				respin2.onStatsWithPBStart(feature, pbcd, firstpr)
// 			}
// 		}
// 	}

// 	return false, 0, 0
// }

// // onStatsWithPBEnding -
// func (respin2 *Respin2) onStatsWithPBEnding(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) error {
// 	pbcd, isok := pbComponentData.(*sgc7pb.Respin2Data)
// 	if !isok {
// 		goutils.Error("Respin2.onStatsWithPBEnding",
// 			zap.Error(ErrIvalidProto))

// 		return ErrIvalidProto
// 	}

// 	if feature.RespinNumStatus != nil {
// 		feature.RespinNumStatus.AddStatus(int(pbcd.CurRespinNum))
// 	}

// 	if feature.RespinWinStatus != nil {
// 		feature.RespinWinStatus.AddStatus(int(pbcd.TotalCoinWin))
// 	}

// 	return nil
// }

// // onStatsWithPBEnding -
// func (respin2 *Respin2) onStatsWithPBStart(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) error {
// 	pbcd, isok := pbComponentData.(*sgc7pb.Respin2Data)
// 	if !isok {
// 		goutils.Error("Respin2.onStatsWithPBStart",
// 			zap.Error(ErrIvalidProto))

// 		return ErrIvalidProto
// 	}

// 	if feature.RespinStartNumStatus != nil {
// 		feature.RespinStartNumStatus.AddStatus(int(pbcd.LastRespinNum))
// 	}

// 	return nil
// }

// // NewComponentData -
// func (respin2 *Respin2) NewComponentData() IComponentData {
// 	return &Respin2Data{}
// }

// // EachUsedResults -
// func (respin2 *Respin2) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
// 	pbcd := &sgc7pb.Respin2Data{}

// 	err := pbComponentData.UnmarshalTo(pbcd)
// 	if err != nil {
// 		goutils.Error("Respin2.EachUsedResults:UnmarshalTo",
// 			zap.Error(err))

// 		return
// 	}

// 	for _, v := range pbcd.BasicComponentData.UsedResults {
// 		oneach(pr.Results[v])
// 	}
// }

// // OnPlayGame - on playgame
// func (respin2 *Respin2) OnPlayGameEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
// 	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

// 	cd := gameProp.MapComponentData[respin2.Name].(*Respin2Data)

// 	cd.TotalCashWin += curpr.CashWin
// 	cd.TotalCoinWin += int64(curpr.CoinWin)

// 	if respin2.Config.IsWinBreak && cd.TotalCoinWin > 0 {
// 		cd.LastRespinNum = 0
// 	}

// 	if cd.LastRespinNum == 0 && cd.LastTriggerNum == 0 {
// 		gameProp.removeRespin(respin2.Name)
// 	}

// 	return nil
// }

// // IsRespin -
// func (respin2 *Respin2) IsRespin() bool {
// 	return true
// }

// // SaveRetriggerRespinNum -
// func (respin2 *Respin2) SaveRetriggerRespinNum(gameProp *GameProperty) {
// 	cd := gameProp.MapComponentData[respin2.Name].(*Respin2Data)

// 	cd.RetriggerAddRespinNum = cd.LastRespinNum
// }

// // Trigger -
// func (respin2 *Respin2) Trigger(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams) {
// 	cd := gameProp.MapComponentData[respin2.Name].(*Respin2Data)

// 	n := cd.TriggerRespinNum[cd.CurTriggerNum]
// 	if n <= 0 {
// 		n = cd.RetriggerAddRespinNum

// 		cd.TriggerRespinNum[cd.CurTriggerNum] = n
// 	}

// 	cd.LastRespinNum += n
// 	cd.CurAddRespinNum += n

// 	cd.CurTriggerNum++

// 	if cd.LastTriggerNum > 0 {
// 		cd.LastTriggerNum--
// 	}

// 	for _, v := range cd.Awards {
// 		if v.TriggerIndex == cd.CurTriggerNum {
// 			gameProp.procAward(plugin, v, curpr, gp, true)
// 		}
// 	}
// }

// // AddRetriggerRespinNum -
// func (respin2 *Respin2) AddRetriggerRespinNum(gameProp *GameProperty, num int) {
// 	cd := gameProp.MapComponentData[respin2.Name].(*Respin2Data)

// 	cd.RetriggerAddRespinNum += num
// }

// // AddTriggerAward -
// func (respin2 *Respin2) AddTriggerAward(gameProp *GameProperty, award *Award) {
// 	cd := gameProp.MapComponentData[respin2.Name].(*Respin2Data)

// 	award.TriggerIndex = cd.CurTriggerNum + cd.LastTriggerNum

// 	cd.Awards = append(cd.Awards, award)
// }

// // PushTrigger -
// func (respin2 *Respin2) PushTrigger(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, num int) {
// 	cd := gameProp.MapComponentData[respin2.Name].(*Respin2Data)

// 	cd.LastTriggerNum++

// 	cd.TriggerRespinNum = append(cd.TriggerRespinNum, num)

// 	// 第一次trigger时，需要直接
// 	if cd.LastRespinNum == 0 && cd.CurRespinNum == 0 {
// 		respin2.Trigger(gameProp, plugin, curpr, gp)
// 	}
// }

// // GetLastRespinNum -
// func (respin2 *Respin2) GetLastRespinNum(gameProp *GameProperty) int {
// 	cd := gameProp.MapComponentData[respin2.Name].(*Respin2Data)

// 	return cd.LastRespinNum
// }

// func NewRespin2(name string) IComponent {
// 	return &Respin2{
// 		BasicComponent: NewBasicComponent(name, 0),
// 	}
// }
