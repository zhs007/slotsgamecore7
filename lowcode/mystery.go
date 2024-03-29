package lowcode

// const MysteryTypeName = "mystery"

// type MysteryData struct {
// 	BasicComponentData
// 	CurMysteryCode int
// }

// // OnNewGame -
// func (mysteryData *MysteryData) OnNewGame(gameProp *GameProperty, component IComponent) {
// 	mysteryData.BasicComponentData.OnNewGame(gameProp, component)
// }

// // OnNewStep -
// func (mysteryData *MysteryData) OnNewStep(gameProp *GameProperty, component IComponent) {
// 	mysteryData.BasicComponentData.OnNewStep(gameProp, component)

// 	mysteryData.CurMysteryCode = -1
// }

// // BuildPBComponentData
// func (mysteryData *MysteryData) BuildPBComponentData() proto.Message {
// 	pbcd := &sgc7pb.MysteryData{
// 		BasicComponentData: mysteryData.BuildPBBasicComponentData(),
// 		CurMysteryCode:     int32(mysteryData.CurMysteryCode),
// 	}

// 	return pbcd
// }

// // MysteryTriggerFeatureConfig - configuration for mystery trigger feature
// type MysteryTriggerFeatureConfig struct {
// 	Symbol               string `yaml:"symbol" json:"symbol"`                             // like LIGHTNING
// 	RespinFirstComponent string `yaml:"respinFirstComponent" json:"respinFirstComponent"` // like lightning
// }

// // MysteryConfig - configuration for Mystery
// type MysteryConfig struct {
// 	BasicComponentConfig   `yaml:",inline" json:",inline"`
// 	MysteryRNG             string                         `yaml:"mysteryRNG" json:"mysteryRNG"` // 强制用已经使用的随机数结果做 Mystery
// 	MysteryWeight          string                         `yaml:"mysteryWeight" json:"mysteryWeight"`
// 	Mystery                string                         `yaml:"mystery" json:"-"`
// 	Mysterys               []string                       `yaml:"mysterys" json:"mysterys"`
// 	MysteryTriggerFeatures []*MysteryTriggerFeatureConfig `yaml:"mysteryTriggerFeatures" json:"mysteryTriggerFeatures"`
// }

// type Mystery struct {
// 	*BasicComponent          `json:"-"`
// 	Config                   *MysteryConfig                       `json:"config"`
// 	MysteryWeights           *sgc7game.ValWeights2                `json:"-"`
// 	MysterySymbols           []int                                `json:"-"`
// 	MapMysteryTriggerFeature map[int]*MysteryTriggerFeatureConfig `json:"-"`
// }

// // maskOtherScene -
// func (mystery *Mystery) maskOtherScene(gameProp *GameProperty, gs *sgc7game.GameScene, symbolCode int) *sgc7game.GameScene {
// 	// cgs := gs.Clone()
// 	cgs := gs.CloneEx(gameProp.PoolScene)

// 	for x, arr := range cgs.Arr {
// 		for y, v := range arr {
// 			if v != symbolCode {
// 				cgs.Arr[x][y] = -1
// 			} else {
// 				cgs.Arr[x][y] = 1
// 			}
// 		}
// 	}

// 	return cgs
// }

// // Init -
// func (mystery *Mystery) Init(fn string, pool *GamePropertyPool) error {
// 	data, err := os.ReadFile(fn)
// 	if err != nil {
// 		goutils.Error("Mystery.Init:ReadFile",
// 			slog.String("fn", fn),
// 			goutils.Err(err))

// 		return err
// 	}

// 	cfg := &MysteryConfig{}

// 	err = yaml.Unmarshal(data, cfg)
// 	if err != nil {
// 		goutils.Error("Mystery.Init:Unmarshal",
// 			slog.String("fn", fn),
// 			goutils.Err(err))

// 		return err
// 	}

// 	return mystery.InitEx(cfg, pool)
// }

// // InitEx -
// func (mystery *Mystery) InitEx(cfg any, pool *GamePropertyPool) error {
// 	mystery.Config = cfg.(*MysteryConfig)
// 	mystery.Config.ComponentType = MysteryTypeName

// 	if mystery.Config.MysteryWeight != "" {
// 		vw2, err := pool.LoadSymbolWeights(mystery.Config.MysteryWeight, "val", "weight", pool.DefaultPaytables, mystery.Config.UseFileMapping)
// 		if err != nil {
// 			goutils.Error("Mystery.Init:LoadSymbolWeights",
// 				slog.String("Weight", mystery.Config.MysteryWeight),
// 				goutils.Err(err))

// 			return err
// 		}

// 		mystery.MysteryWeights = vw2
// 	}

// 	if len(mystery.Config.Mysterys) > 0 {
// 		for _, v := range mystery.Config.Mysterys {
// 			mystery.MysterySymbols = append(mystery.MysterySymbols, pool.DefaultPaytables.MapSymbols[v])
// 		}
// 	} else {
// 		mystery.MysterySymbols = append(mystery.MysterySymbols, pool.DefaultPaytables.MapSymbols[mystery.Config.Mystery])
// 	}

// 	for _, v := range mystery.Config.MysteryTriggerFeatures {
// 		symbolCode := pool.DefaultPaytables.MapSymbols[v.Symbol]

// 		mystery.MapMysteryTriggerFeature[symbolCode] = v
// 	}

// 	mystery.onInit(&mystery.Config.BasicComponentConfig)

// 	return nil
// }

// // playgame
// func (mystery *Mystery) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
// 	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

// 	// mystery.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

// 	mcd := cd.(*MysteryData)

// 	gs := mystery.GetTargetScene3(gameProp, curpr, prs, 0)
// 	if !gs.HasSymbols(mystery.MysterySymbols) {
// 		// mystery.ReTagScene(gameProp, curpr, mcd.TargetSceneIndex, &mcd.BasicComponentData)
// 	} else {
// 		if mystery.MysteryWeights != nil {
// 			if mystery.Config.MysteryRNG != "" {
// 				rng := gameProp.GetTagInt(mystery.Config.MysteryRNG)
// 				cs := mystery.MysteryWeights.Vals[rng]

// 				curmcode := cs.Int()
// 				mcd.CurMysteryCode = curmcode

// 				// gameProp.SetVal(GamePropCurMystery, curmcode)

// 				// sc2 := gs.Clone()
// 				sc2 := gs.CloneEx(gameProp.PoolScene)
// 				for _, v := range mystery.MysterySymbols {
// 					sc2.ReplaceSymbol(v, curmcode)
// 				}

// 				mystery.AddScene(gameProp, curpr, sc2, &mcd.BasicComponentData)

// 				v, isok := mystery.MapMysteryTriggerFeature[curmcode]
// 				if isok {
// 					if v.RespinFirstComponent != "" {
// 						os := mystery.maskOtherScene(gameProp, sc2, curmcode)

// 						gameProp.Respin(curpr, gp, v.RespinFirstComponent, sc2, os)

// 						return v.RespinFirstComponent, nil
// 					}
// 				}
// 			} else {
// 				curm, err := mystery.MysteryWeights.RandVal(plugin)
// 				if err != nil {
// 					goutils.Error("Mystery.OnPlayGame:RandVal",
// 						goutils.Err(err))

// 					return "", err
// 				}

// 				curmcode := curm.Int()

// 				// gameProp.SetVal(GamePropCurMystery, curm.Int())

// 				sc2 := gs.CloneEx(gameProp.PoolScene)
// 				// sc2 := gs.Clone()
// 				for _, v := range mystery.MysterySymbols {
// 					sc2.ReplaceSymbol(v, curm.Int())
// 				}

// 				mystery.AddScene(gameProp, curpr, sc2, &mcd.BasicComponentData)

// 				v, isok := mystery.MapMysteryTriggerFeature[curmcode]
// 				if isok {
// 					if v.RespinFirstComponent != "" {
// 						os := mystery.maskOtherScene(gameProp, sc2, curmcode)

// 						gameProp.Respin(curpr, gp, v.RespinFirstComponent, sc2, os)

// 						return v.RespinFirstComponent, nil
// 					}
// 				}
// 			}
// 		}
// 	}

// 	nc := mystery.onStepEnd(gameProp, curpr, gp, "")

// 	// gp.AddComponentData(mystery.Name, cd)

// 	return nc, nil
// }

// // OnAsciiGame - outpur to asciigame
// func (mystery *Mystery) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
// 	mcd := cd.(*MysteryData)

// 	if len(mcd.UsedScenes) > 0 {
// 		if mystery.MysteryWeights != nil {
// 			fmt.Printf("mystery is %v\n", gameProp.CurPaytables.GetStringFromInt(mcd.CurMysteryCode))
// 			asciigame.OutputScene("after symbols", pr.Scenes[mcd.UsedScenes[0]], mapSymbolColor)
// 		}
// 	}

// 	return nil
// }

// // // OnStats
// // func (mystery *Mystery) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// // 	return false, 0, 0
// // }

// // // OnStatsWithPB -
// // func (mystery *Mystery) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// // 	pbcd, isok := pbComponentData.(*sgc7pb.MysteryData)
// // 	if !isok {
// // 		goutils.Error("Mystery.OnStatsWithPB",
// // 			goutils.Err(ErrIvalidProto))

// // 		return 0, ErrIvalidProto
// // 	}

// // 	return mystery.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// // }

// // NewComponentData -
// func (mystery *Mystery) NewComponentData() IComponentData {
// 	return &MysteryData{}
// }

// // EachUsedResults -
// func (mystery *Mystery) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
// 	pbcd := &sgc7pb.MysteryData{}

// 	err := pbComponentData.UnmarshalTo(pbcd)
// 	if err != nil {
// 		goutils.Error("Mystery.EachUsedResults:UnmarshalTo",
// 			goutils.Err(err))

// 		return
// 	}

// 	for _, v := range pbcd.BasicComponentData.UsedResults {
// 		oneach(pr.Results[v])
// 	}
// }

// func NewMystery(name string) IComponent {
// 	mystery := &Mystery{
// 		BasicComponent:           NewBasicComponent(name, 1),
// 		MapMysteryTriggerFeature: make(map[int]*MysteryTriggerFeatureConfig),
// 	}

// 	return mystery
// }
