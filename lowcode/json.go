package lowcode

import (
	"os"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
)

func loadBasicInfo(cfg *Config, buf []byte) error {
	gameName, err := sonic.Get(buf, "gameName")
	if err != nil {
		goutils.Error("loadBasicInfo:Get",
			zap.String("key", "gameName"),
			zap.Error(err))

		return err
	}

	cfg.Name, _ = gameName.String()

	lstParam, err := sonic.Get(buf, "parameter")
	if err != nil {
		goutils.Error("loadBasicInfo:Get",
			zap.String("key", "parameter"),
			zap.Error(err))

		return err
	}

	lst, err := lstParam.ArrayUseNode()
	if err != nil {
		goutils.Error("loadBasicInfo:ArrayUseNode",
			zap.Error(err))

		return err
	}

	for i, v := range lst {
		str, err := v.Get("name").String()
		if err != nil {
			goutils.Error("loadBasicInfo:name",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}

		if str == "Width" {
			w, err := v.Get("value").Int64()
			if err != nil {
				goutils.Error("loadBasicInfo:Width",
					zap.Int("i", i),
					zap.Error(err))

				return ErrIvalidWidth
			}

			cfg.Width = int(w)
		} else if str == "Height" {
			h, err := v.Get("value").Int64()
			if err != nil {
				goutils.Error("loadBasicInfo:Height",
					zap.Int("i", i),
					zap.Error(err))

				return ErrIvalidHeight
			}

			cfg.Height = int(h)
		} else if str == "Scene" {
			scene, err := v.Get("value").String()
			if err != nil {
				goutils.Error("loadBasicInfo:Scene",
					zap.Int("i", i),
					zap.Error(err))

				return ErrIvalidDefaultScene
			}

			cfg.DefaultScene = scene
		}
	}

	return nil
}

func parse2IntSlice(n *ast.Node) ([]int, error) {
	arr, err := n.ArrayUseNode()
	if err != nil {
		goutils.Error("parse2IntSlice:Array",
			zap.Error(err))

		return nil, err
	}

	iarr := []int{}

	for i, v := range arr {
		iv, err := v.Int64()
		if err != nil {
			goutils.Error("parse2IntSlice:Int64",
				zap.Int("i", i),
				zap.Error(err))

			return nil, err
		}

		iarr = append(iarr, int(iv))
	}

	return iarr, nil
}

func parse2StringSlice(n *ast.Node) ([]string, error) {
	arr, err := n.ArrayUseNode()
	if err != nil {
		goutils.Error("parse2StringSlice:Array",
			zap.Error(err))

		return nil, err
	}

	strarr := []string{}

	for i, v := range arr {
		strv, err := v.String()
		if err != nil {
			goutils.Error("parse2StringSlice:String",
				zap.Int("i", i),
				zap.Error(err))

			return nil, err
		}

		strarr = append(strarr, (strv))
	}

	return strarr, nil
}

func parsePaytables(n *ast.Node) (*sgc7game.PayTables, error) {
	if n == nil {
		goutils.Error("parsePaytables",
			zap.Error(ErrIvalidPayTables))

		return nil, ErrIvalidPayTables
	}

	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parsePaytables:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	dataPaytables := []*paytableData{}

	err = sonic.Unmarshal(buf, &dataPaytables)
	if err != nil {
		goutils.Error("parsePaytables:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	paytables := &sgc7game.PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	for _, node := range dataPaytables {
		paytables.MapPay[node.Code] = node.Data
		paytables.MapSymbols[node.Symbol] = node.Code
	}

	return paytables, nil
}

func loadPaytables(cfg *Config, lstPaytables *ast.Node) error {
	lst, err := lstPaytables.ArrayUseNode()
	if err != nil {
		goutils.Error("loadPaytables:ArrayUseNode",
			zap.Error(err))

		return err
	}

	for i, v := range lst {
		name, err := v.Get("fileName").String()
		if err != nil {
			goutils.Error("loadPaytables:fileName",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}

		paytables, err := parsePaytables(v.Get("fileJson"))
		if err != nil {
			goutils.Error("loadPaytables:parsePaytables",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}

		cfg.Paytables[name] = name
		cfg.MapPaytables[name] = paytables

		if i == 0 {
			cfg.DefaultPaytables = name
		}
	}

	return nil
}

func parseLineData(n *ast.Node, width int) (*sgc7game.LineData, error) {
	if n == nil {
		goutils.Error("parseLineData",
			zap.Error(ErrIvalidReels))

		return nil, ErrIvalidReels
	}

	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseLineData:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	dataLines := [][]int{}

	err = sonic.Unmarshal(buf, &dataLines)
	if err != nil {
		goutils.Error("parseLineData:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	return &sgc7game.LineData{
		Lines: dataLines,
	}, nil
}

func parseReels(n *ast.Node, paytables *sgc7game.PayTables) (*sgc7game.ReelsData, error) {
	if n == nil {
		goutils.Error("parseReels",
			zap.Error(ErrIvalidReels))

		return nil, ErrIvalidReels
	}

	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseReels:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	dataReels := [][]string{}

	err = sonic.Unmarshal(buf, &dataReels)
	if err != nil {
		goutils.Error("parseReels:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	reelsd := &sgc7game.ReelsData{}

	for x, arr := range dataReels {
		reeld := []int{}

		for y, strSym := range arr {
			sc, isok := paytables.MapSymbols[strSym]
			if !isok {
				goutils.Error("parseReels:MapSymbols",
					zap.Int("x", x),
					zap.Int("y", y),
					zap.Error(ErrIvalidSymbolInReels))

				return nil, ErrIvalidSymbolInReels
			}

			reeld = append(reeld, sc)
		}

		reelsd.Reels = append(reelsd.Reels, reeld)
	}

	return reelsd, nil
}

func loadOtherList(cfg *Config, lstOther *ast.Node) error {
	lst, err := lstOther.ArrayUseNode()
	if err != nil {
		goutils.Error("loadOtherList:ArrayUseNode",
			zap.Error(err))

		return err
	}

	for i, v := range lst {
		name, err := v.Get("fileName").String()
		if err != nil {
			goutils.Error("loadOtherList:fileName",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}

		t, err := v.Get("type").String()
		if err != nil {
			goutils.Error("loadOtherList:type",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}

		if t == "Linedata" {
			ld, err := parseLineData(v.Get("fileJson"), cfg.Width)
			if err != nil {
				goutils.Error("loadOtherList:parseLineData",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}

			cfg.Linedata[name] = name
			cfg.MapLinedate[name] = ld

			if len(cfg.Linedata) == 1 {
				cfg.DefaultLinedata = name
			}
		} else if t == "Reels" {
			rd, err := parseReels(v.Get("fileJson"), cfg.GetDefaultPaytables())
			if err != nil {
				goutils.Error("loadOtherList:parseReels",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}

			cfg.Reels[name] = name
			cfg.MapReels[name] = rd
		} else if t == "Weights" {
			vw2, err := parseValWeights(v.Get("fileJson"))
			if err != nil {
				goutils.Error("loadOtherList:parseValWeights",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}

			cfg.mapValWeights[name] = vw2
		} else if t == "ReelSetWeight" {
			vw2, err := parseReelSetWeights(v.Get("fileJson"))
			if err != nil {
				goutils.Error("loadOtherList:parseReelSetWeights",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}

			cfg.mapReelSetWeights[name] = vw2
		} else if t == "SymbolWeight" {
			vw2, err := parseSymbolWeights(v.Get("fileJson"), cfg.GetDefaultPaytables())
			if err != nil {
				goutils.Error("loadOtherList:parseSymbolWeights",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}

			cfg.mapValWeights[name] = vw2
		} else if t == "IntValMappingFile" {
			vm2, err := parseIntValMappingFile(v.Get("fileJson"))
			if err != nil {
				goutils.Error("loadOtherList:parseSymbolWeights",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}

			cfg.mapIntMapping[name] = vm2
		} else if t == "StringValWeight" {
			vm2, err := parseStrWeights(v.Get("fileJson"))
			if err != nil {
				goutils.Error("loadOtherList:parseStrWeights",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}

			cfg.mapStrWeights[name] = vm2
		} else {
			goutils.Error("loadOtherList",
				zap.Int("i", i),
				zap.String("type", t),
				zap.Error(ErrUnsupportedOtherList))

			return ErrUnsupportedOtherList
		}
	}

	return nil
}

type linkData struct {
	mapLinks map[string][][]string
}

func (ld *linkData) add(linktype string, src string, target string) {
	_, isok := ld.mapLinks[linktype]
	if !isok {
		ld.mapLinks[linktype] = [][]string{}
	}

	ld.mapLinks[linktype] = append(ld.mapLinks[linktype], []string{src, target})
}

func newLinkData() *linkData {
	return &linkData{
		mapLinks: make(map[string][][]string),
	}
}

func loadCells(cfg *Config, bet int, cells *ast.Node) error {
	// linkScene := [][]string{}
	// linkOtherScene := [][]string{}
	// linkComponent := [][]string{}
	// jumpComponent := [][]string{}
	// loopComponent := [][]string{}
	ldid := newLinkData()
	// lstStart := []string{}
	lstStartID := []string{}
	// mapTrigger := make(map[string]*TriggerFeatureConfig)
	// mapTriggerID := make(map[string]*TriggerFeatureConfig)
	// lstBasicWins := []*BasicWinsConfig{}
	mapComponentName := make(map[string]string)

	lst, err := cells.ArrayUseNode()
	if err != nil {
		goutils.Error("loadCells:ArrayUseNode",
			zap.Error(err))

		return err
	}

	// startid := ""

	for i, cell := range lst {
		shape, err := cell.Get("shape").String()
		if err != nil {
			goutils.Error("loadCells:Get:shape",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}

		id, err := cell.Get("id").String()
		if err != nil {
			goutils.Error("loadCells:Get:id",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}

		// if shape == "custom-node-width-start" {
		// 	startid = id
		// } else
		if shape == "custom-node" {
			componentType, err := cell.Get("label").String()
			if err != nil {
				goutils.Error("loadCells:Get:label",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}

			componentType = strings.ToLower(componentType)

			if componentType == "weightreels" {
				componentName, err := parseWeightReels(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseWeightReels",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "basicreels" {
				componentName, err := parseBasicReels(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseBasicReels",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "scattertrigger" {
				componentName, err := parseScatterTrigger(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseScatterTrigger",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "linestrigger" {
				componentName, err := parseLinesTrigger(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseLinesTrigger",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "waystrigger" {
				componentName, err := parseWaysTrigger(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseWaysTrigger",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "movesymbol" {
				componentName, err := parseMoveSymbol(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseMoveSymbol",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "respin" {
				componentName, err := parseRespin(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseRespin",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "symbolcollection" {
				componentName, err := parseSymbolCollection2(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseSymbolCollection2",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "removesymbols" {
				componentName, err := parseRemoveSymbols(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseRemoveSymbols",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "dropdownsymbols" {
				componentName, err := parseDropDownSymbols(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseDropDownSymbols",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "refillsymbols" {
				componentName, err := parseRefillSymbols(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseRefillSymbols",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "collector" {
				componentName, err := parseCollector(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseCollector",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "queuebranch" || componentType == "delayqueue" {
				componentName, err := parseQueueBranch(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseQueueBranch",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "replacesymbolgroup" {
				componentName, err := parseReplaceSymbolGroup(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseReplaceSymbolGroup",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "rollsymbol" {
				componentName, err := parseRollSymbol(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseRollSymbol",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "mask" {
				componentName, err := parseMask(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseMask",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "replacereelwithmask" {
				componentName, err := parseReplaceReelWithMask(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseReplaceReelWithMask",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "piggybank" {
				componentName, err := parsePiggyBank(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parsePiggyBank",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "addsymbols" {
				componentName, err := parseAddSymbols(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseAddSymbols",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "intvalmapping" {
				componentName, err := parseIntValMapping(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseIntValMapping",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "weightbranch" {
				componentName, err := parseWeightBranch(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseWeightBranch",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "clustertrigger" {
				componentName, err := parseClusterTrigger(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseClusterTrigger",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapComponentName[id] = componentName
			} else {
				goutils.Error("loadCells:ErrUnsupportedComponentType",
					zap.String("componentType", componentType),
					zap.Error(ErrUnsupportedComponentType))

				return ErrUnsupportedComponentType
			}
		} else if shape == "edge" {
			source, err := cell.Get("source").Get("cell").String()
			if err != nil {
				goutils.Error("loadCells:edge:source:cell",
					zap.Error(err))

				return err
			}

			sourcePort, err := cell.Get("source").Get("port").String()
			if err != nil {
				goutils.Error("loadCells:edge:source:port",
					zap.Error(err))

				return err
			}

			target, err := cell.Get("target").Get("cell").String()
			if err != nil {
				goutils.Error("loadCells:edge:target",
					zap.Error(err))

				return err
			}

			// if source == startid {
			// 	lstStart = append(lstStart, mapComponentName[target])
			// } else {
			if sourcePort == "jump-component-groups-out" {
				ldid.add("jump", source, target)
				// jumpComponent = append(jumpComponent, []string{mapComponentName[source], mapComponentName[target]})
			} else if sourcePort == "component-groups-out" {
				ldid.add("next", source, target)
				// linkComponent = append(linkComponent, []string{mapComponentName[source], mapComponentName[target]})
			} else if sourcePort == "loop-component-groups-out" {
				ldid.add("loop", source, target)
				// loopComponent = append(loopComponent, []string{mapComponentName[source], mapComponentName[target]})
			} else if sourcePort == "foreach-component-groups-out" {
				ldid.add("foreach", source, target)
				// loopComponent = append(loopComponent, []string{mapComponentName[source], mapComponentName[target]})
			} else if sourcePort == "start-out" {
				lstStartID = append(lstStartID, target)
				// ld.add("foreach", mapComponentName[source], mapComponentName[target])
				// loopComponent = append(loopComponent, []string{mapComponentName[source], mapComponentName[target]})
			} else {
				arr := strings.Split(sourcePort, "vals-component-groups-out-")
				if len(arr) == 2 {
					ldid.add(arr[1], source, target)
				} else {
					goutils.Error("loadCells:sourcePort",
						zap.String("sourcePort", sourcePort),
						zap.Error(ErrUnsupportedLinkType))

					return ErrUnsupportedLinkType
				}
			}
			// }
		}
	}

	if len(lstStartID) > 0 {
		cfg.StartComponents[bet] = mapComponentName[lstStartID[0]]
	}

	for lt, arr := range ldid.mapLinks {
		for _, cld := range arr {
			icfg, isok := cfg.mapConfig[mapComponentName[cld[0]]]
			if isok {
				icfg.SetLinkComponent(lt, mapComponentName[cld[1]])
			}
		}
	}

	// for _, arr := range linkComponent {
	// 	icfg, isok := cfg.mapConfig[arr[0]]
	// 	if isok {
	// 		icfg.SetLinkComponent("next", arr[1])
	// 	}
	// }

	// for _, arr := range jumpComponent {
	// 	icfg, isok := cfg.mapConfig[arr[0]]
	// 	if isok {
	// 		icfg.SetLinkComponent("jump", arr[1])
	// 	}
	// }

	// for _, arr := range loopComponent {
	// 	icfg, isok := cfg.mapConfig[arr[0]]
	// 	if isok {
	// 		icfg.SetLinkComponent("loop", arr[1])
	// 	}
	// }

	// for _, arr := range linkScene {
	// 	sourceCfg, isok0 := cfg.mapBasicConfig[arr[0]]
	// 	if isok0 {
	// 		sourceCfg.TagScenes = append(sourceCfg.TagScenes, arr[0])
	// 	}

	// 	targetCfg := cfg.mapBasicConfig[arr[1]]
	// 	if targetCfg != nil {
	// 		targetCfg.TargetScene = arr[0]
	// 	}

	// 	triggerCfg := mapTriggerID[arr[1]]
	// 	if triggerCfg != nil {
	// 		triggerCfg.TargetScene = arr[0]
	// 	}
	// }

	// for _, arr := range linkOtherScene {
	// 	sourceCfg := cfg.mapBasicConfig[arr[0]]
	// 	if sourceCfg != nil {
	// 		sourceCfg.TagOtherScenes = append(sourceCfg.TagOtherScenes, arr[0])
	// 	}

	// 	targetCfg := cfg.mapBasicConfig[arr[1]]
	// 	if targetCfg != nil {
	// 		targetCfg.TargetOtherScene = arr[0]
	// 	}
	// }

	// for _, basicWinsCfg := range lstBasicWins {
	// 	for _, k := range basicWinsCfg.BeforMainTriggerName {
	// 		cfg, isok := mapTrigger[k]
	// 		if !isok {
	// 			goutils.Error("loadCells:BeforMain",
	// 				zap.String("label", k),
	// 				zap.Error(ErrIvalidTriggerLabel))

	// 			return ErrIvalidTriggerLabel
	// 		}

	// 		basicWinsCfg.BeforMain = append(basicWinsCfg.BeforMain, cfg)
	// 	}

	// 	for _, k := range basicWinsCfg.AfterMainTriggerName {
	// 		cfg, isok := mapTrigger[k]
	// 		if !isok {
	// 			goutils.Error("loadCells:AfterMain",
	// 				zap.String("label", k),
	// 				zap.Error(ErrIvalidTriggerLabel))

	// 			return ErrIvalidTriggerLabel
	// 		}

	// 		basicWinsCfg.AfterMain = append(basicWinsCfg.AfterMain, cfg)
	// 	}
	// }

	return nil
}

func loadBetMethod(cfg *Config, betMethod *ast.Node) error {
	bet, err := betMethod.Get("bet").Int64()
	if err != nil {
		goutils.Error("loadBetMethod:Get:bet",
			zap.Error(err))

		return err
	}

	cfg.Bets = append(cfg.Bets, int(bet))
	cfg.TotalBetInWins = append(cfg.TotalBetInWins, int(bet))

	err = loadCells(cfg, int(bet), betMethod.Get("graph").Get("cells"))
	if err != nil {
		goutils.Error("loadBetMethod:loadCells",
			zap.Error(err))

		return err
	}

	return nil
}

func NewGame2(fn string, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("NewGame2:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	return NewGame2WithData(data, funcNewPlugin)
}

func NewGame2ForRTP(bet int, fn string, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("NewGame2:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	return NewGame2WithData(data, funcNewPlugin)
}

func NewGame3(fn string, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
	if strings.Contains(fn, ".json") {
		return NewGame2(fn, funcNewPlugin)
	}

	return NewGameEx(fn, funcNewPlugin)
}

func NewGame3ForRTP(bet int, fn string, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
	if strings.Contains(fn, ".json") {
		return NewGame2(fn, funcNewPlugin)
	}

	return NewGameExForRTP(bet, fn, funcNewPlugin)
}

func NewGame2WithData(data []byte, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
	game := &Game{
		BasicGame:    sgc7game.NewBasicGame(funcNewPlugin),
		MgrComponent: NewComponentMgr(),
	}

	cfg := &Config{
		Paytables:         make(map[string]string),
		MapPaytables:      make(map[string]*sgc7game.PayTables),
		Linedata:          make(map[string]string),
		MapLinedate:       make(map[string]*sgc7game.LineData),
		Reels:             make(map[string]string),
		MapReels:          make(map[string]*sgc7game.ReelsData),
		mapConfig:         make(map[string]IComponentConfig),
		StartComponents:   make(map[int]string),
		mapBasicConfig:    make(map[string]*BasicComponentConfig),
		mapValWeights:     make(map[string]*sgc7game.ValWeights2),
		mapReelSetWeights: make(map[string]*sgc7game.ValWeights2),
		mapStrWeights:     make(map[string]*sgc7game.ValWeights2),
		mapIntMapping:     make(map[string]*sgc7game.ValMapping2),
		// mapBetConfig:    make(map[int]*BetDataConfig),
	}

	err := loadBasicInfo(cfg, data)
	if err != nil {
		goutils.Error("NewGame2WithData:loadBasicInfo",
			zap.Error(err))

		return nil, err
	}

	lstPaytables, err := sonic.Get(data, "repository", "paytableList")
	if err != nil {
		goutils.Error("NewGame2WithData:Get",
			zap.String("key", "repository.paytableList"),
			zap.Error(err))

		return nil, err
	}

	err = loadPaytables(cfg, &lstPaytables)
	if err != nil {
		goutils.Error("NewGame2WithData:loadPaytables",
			zap.Error(err))

		return nil, err
	}

	lstOther, err := sonic.Get(data, "repository", "otherList")
	if err != nil {
		goutils.Error("NewGame2WithData:Get",
			zap.String("key", "repository.otherList"),
			zap.Error(err))

		return nil, err
	}

	err = loadOtherList(cfg, &lstOther)
	if err != nil {
		goutils.Error("NewGame2WithData:loadOtherList",
			zap.Error(err))

		return nil, err
	}

	cfgGameMod := &GameModConfig{}
	cfgGameMod.Type = "bg"
	cfg.GameMods = append(cfg.GameMods, cfgGameMod)

	betMethod, err := sonic.Get(data, "betMethod", 0)
	if err != nil {
		goutils.Error("NewGame2WithData:Get",
			zap.String("key", "betMethod[0]"),
			zap.Error(err))

		return nil, err
	}

	err = loadBetMethod(cfg, &betMethod)
	if err != nil {
		goutils.Error("NewGame2WithData:loadBetMethod",
			zap.Error(err))

		return nil, err
	}

	cfg.RTP = &RTPConfig{}

	err = game.Init2(cfg)
	if err != nil {
		goutils.Error("NewGame2WithData:Init2",
			zap.Error(err))

		return nil, err
	}

	return game, nil
}
