package lowcode

import (
	"log/slog"
	"os"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func loadBasicInfo(cfg *Config, buf []byte) error {
	gameName, err := sonic.Get(buf, "gameName")
	if err != nil {
		goutils.Error("loadBasicInfo:Get",
			slog.String("key", "gameName"),
			goutils.Err(err))

		return err
	}

	cfg.Name, _ = gameName.String()

	lstParam, err := sonic.Get(buf, "parameter")
	if err != nil {
		goutils.Error("loadBasicInfo:Get",
			slog.String("key", "parameter"),
			goutils.Err(err))

		return err
	}

	lst, err := lstParam.ArrayUseNode()
	if err != nil {
		goutils.Error("loadBasicInfo:ArrayUseNode",
			goutils.Err(err))

		return err
	}

	for i, v := range lst {
		str, err := v.Get("name").String()
		if err != nil {
			goutils.Error("loadBasicInfo:name",
				slog.Int("i", i),
				goutils.Err(err))

			return err
		}

		if str == "Width" {
			w, err := v.Get("value").Int64()
			if err != nil {
				goutils.Error("loadBasicInfo:Width",
					slog.Int("i", i),
					goutils.Err(err))

				return ErrIvalidWidth
			}

			cfg.Width = int(w)
		} else if str == "Height" {
			h, err := v.Get("value").Int64()
			if err != nil {
				goutils.Error("loadBasicInfo:Height",
					slog.Int("i", i),
					goutils.Err(err))

				return ErrIvalidHeight
			}

			cfg.Height = int(h)
		} else if str == "Scene" {
			scene, err := v.Get("value").String()
			if err != nil {
				goutils.Error("loadBasicInfo:Scene",
					slog.Int("i", i),
					goutils.Err(err))

				return ErrIvalidDefaultScene
			}

			cfg.DefaultScene = scene
		}
	}

	return nil
}

// func parse2IntSlice(n *ast.Node) ([]int, error) {
// 	arr, err := n.ArrayUseNode()
// 	if err != nil {
// 		goutils.Error("parse2IntSlice:Array",
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	iarr := []int{}

// 	for i, v := range arr {
// 		iv, err := v.Int64()
// 		if err != nil {
// 			goutils.Error("parse2IntSlice:Int64",
// 				slog.Int("i", i),
// 				goutils.Err(err))

// 			return nil, err
// 		}

// 		iarr = append(iarr, int(iv))
// 	}

// 	return iarr, nil
// }

// func parse2StringSlice(n *ast.Node) ([]string, error) {
// 	arr, err := n.ArrayUseNode()
// 	if err != nil {
// 		goutils.Error("parse2StringSlice:Array",
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	strarr := []string{}

// 	for i, v := range arr {
// 		strv, err := v.String()
// 		if err != nil {
// 			goutils.Error("parse2StringSlice:String",
// 				slog.Int("i", i),
// 				goutils.Err(err))

// 			return nil, err
// 		}

// 		strarr = append(strarr, (strv))
// 	}

// 	return strarr, nil
// }

func parsePaytables(n *ast.Node) (*sgc7game.PayTables, error) {
	if n == nil {
		goutils.Error("parsePaytables",
			goutils.Err(ErrIvalidPayTables))

		return nil, ErrIvalidPayTables
	}

	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parsePaytables:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	dataPaytables := []*paytableData{}

	err = sonic.Unmarshal(buf, &dataPaytables)
	if err != nil {
		goutils.Error("parsePaytables:Unmarshal",
			goutils.Err(err))

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
			goutils.Err(err))

		return err
	}

	for i, v := range lst {
		name, err := v.Get("fileName").String()
		if err != nil {
			goutils.Error("loadPaytables:fileName",
				slog.Int("i", i),
				goutils.Err(err))

			return err
		}

		paytables, err := parsePaytables(v.Get("fileJson"))
		if err != nil {
			goutils.Error("loadPaytables:parsePaytables",
				slog.Int("i", i),
				goutils.Err(err))

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

func parseLineData(n *ast.Node, _ int) (*sgc7game.LineData, error) {
	if n == nil {
		goutils.Error("parseLineData",
			goutils.Err(ErrIvalidReels))

		return nil, ErrIvalidReels
	}

	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseLineData:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	dataLines := [][]int{}

	err = sonic.Unmarshal(buf, &dataLines)
	if err != nil {
		goutils.Error("parseLineData:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	return &sgc7game.LineData{
		Lines: dataLines,
	}, nil
}

func parseReels(n *ast.Node, paytables *sgc7game.PayTables) (*sgc7game.ReelsData, error) {
	if n == nil {
		goutils.Error("parseReels",
			goutils.Err(ErrIvalidReels))

		return nil, ErrIvalidReels
	}

	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseReels:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	dataReels := [][]string{}

	err = sonic.Unmarshal(buf, &dataReels)
	if err != nil {
		goutils.Error("parseReels:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	reelsd := &sgc7game.ReelsData{}

	for x, arr := range dataReels {
		reeld := []int{}

		for y, strSym := range arr {
			sc, isok := paytables.MapSymbols[strSym]
			if !isok {
				goutils.Error("parseReels:MapSymbols",
					slog.Int("x", x),
					slog.Int("y", y),
					goutils.Err(ErrIvalidSymbolInReels))

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
			goutils.Err(err))

		return err
	}

	for i, v := range lst {
		name, err := v.Get("fileName").String()
		if err != nil {
			goutils.Error("loadOtherList:fileName",
				slog.Int("i", i),
				goutils.Err(err))

			return err
		}

		t, err := v.Get("type").String()
		if err != nil {
			goutils.Error("loadOtherList:type",
				slog.Int("i", i),
				goutils.Err(err))

			return err
		}

		if t == "Linedata" {
			ld, err := parseLineData(v.Get("fileJson"), cfg.Width)
			if err != nil {
				goutils.Error("loadOtherList:parseLineData",
					slog.Int("i", i),
					goutils.Err(err))

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
					slog.Int("i", i),
					goutils.Err(err))

				return err
			}

			cfg.Reels[name] = name
			cfg.MapReels[name] = rd
		} else if t == "Weights" {
			vw2, err := parseValWeights(v.Get("fileJson"))
			if err != nil {
				goutils.Error("loadOtherList:parseValWeights",
					slog.Int("i", i),
					goutils.Err(err))

				return err
			}

			cfg.mapValWeights[name] = vw2
		} else if t == "ReelSetWeight" {
			vw2, err := parseReelSetWeights(v.Get("fileJson"))
			if err != nil {
				goutils.Error("loadOtherList:parseReelSetWeights",
					slog.Int("i", i),
					goutils.Err(err))

				return err
			}

			cfg.mapReelSetWeights[name] = vw2
		} else if t == "SymbolWeight" {
			vw2, err := parseSymbolWeights(v.Get("fileJson"), cfg.GetDefaultPaytables())
			if err != nil {
				goutils.Error("loadOtherList:parseSymbolWeights",
					slog.Int("i", i),
					goutils.Err(err))

				return err
			}

			cfg.mapValWeights[name] = vw2
		} else if t == "IntValMappingFile" {
			vm2, err := parseIntValMappingFile(v.Get("fileJson"))
			if err != nil {
				goutils.Error("loadOtherList:parseSymbolWeights",
					slog.Int("i", i),
					goutils.Err(err))

				return err
			}

			cfg.mapIntMapping[name] = vm2
		} else if t == "StringValWeight" {
			vm2, err := parseStrWeights(v.Get("fileJson"))
			if err != nil {
				goutils.Error("loadOtherList:parseStrWeights",
					slog.Int("i", i),
					goutils.Err(err))

				return err
			}

			cfg.mapStrWeights[name] = vm2
		} else if t == "IntValWeight" {
			vm2, err := parseIntValWeights(v.Get("fileJson"))
			if err != nil {
				goutils.Error("loadOtherList:parseIntValWeights",
					slog.Int("i", i),
					goutils.Err(err))

				return err
			}

			cfg.mapValWeights[name] = vm2
		} else {
			goutils.Error("loadOtherList",
				slog.Int("i", i),
				slog.String("type", t),
				goutils.Err(ErrUnsupportedOtherList))

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

func loadCells(cfg *BetConfig, cells *ast.Node) error {
	ldid := newLinkData()
	lstStartID := []string{}
	mapComponentName := make(map[string]string)

	lst, err := cells.ArrayUseNode()
	if err != nil {
		goutils.Error("loadCells:ArrayUseNode",
			goutils.Err(err))

		return err
	}

	for i, cell := range lst {
		shape, err := cell.Get("shape").String()
		if err != nil {
			goutils.Error("loadCells:Get:shape",
				slog.Int("i", i),
				goutils.Err(err))

			return err
		}

		id, err := cell.Get("id").String()
		if err != nil {
			goutils.Error("loadCells:Get:id",
				slog.Int("i", i),
				goutils.Err(err))

			return err
		}

		// if shape == "custom-node-width-start" {
		// 	startid = id
		// } else
		if shape == "custom-node" {
			componentType, err := cell.Get("label").String()
			if err != nil {
				goutils.Error("loadCells:Get:label",
					slog.Int("i", i),
					goutils.Err(err))

				return err
			}

			componentType = strings.ToLower(componentType)

			if componentType == "weightreels" {
				componentName, err := parseWeightReels(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseWeightReels",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "basicreels" {
				componentName, err := parseBasicReels(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseBasicReels",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "scattertrigger" {
				componentName, err := parseScatterTrigger(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseScatterTrigger",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "linestrigger" {
				componentName, err := parseLinesTrigger(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseLinesTrigger",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "waystrigger" {
				componentName, err := parseWaysTrigger(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseWaysTrigger",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "movesymbol" {
				componentName, err := parseMoveSymbol(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseMoveSymbol",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "respin" {
				componentName, err := parseRespin(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseRespin",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "symbolcollection" {
				componentName, err := parseSymbolCollection2(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseSymbolCollection2",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "removesymbols" {
				componentName, err := parseRemoveSymbols(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseRemoveSymbols",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "dropdownsymbols" {
				componentName, err := parseDropDownSymbols(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseDropDownSymbols",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "refillsymbols" {
				componentName, err := parseRefillSymbols(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseRefillSymbols",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "collector" {
				componentName, err := parseCollector(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseCollector",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "queuebranch" || componentType == "delayqueue" {
				componentName, err := parseQueueBranch(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseQueueBranch",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "replacesymbolgroup" {
				componentName, err := parseReplaceSymbolGroup(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseReplaceSymbolGroup",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "rollsymbol" {
				componentName, err := parseRollSymbol(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseRollSymbol",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "mask" {
				componentName, err := parseMask(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseMask",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "replacereelwithmask" {
				componentName, err := parseReplaceReelWithMask(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseReplaceReelWithMask",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "piggybank" {
				componentName, err := parsePiggyBank(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parsePiggyBank",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "addsymbols" {
				componentName, err := parseAddSymbols(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseAddSymbols",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "intvalmapping" {
				componentName, err := parseIntValMapping(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseIntValMapping",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "weightbranch" {
				componentName, err := parseWeightBranch(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseWeightBranch",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "clustertrigger" {
				componentName, err := parseClusterTrigger(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseClusterTrigger",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "gengigasymbol" {
				componentName, err := parseGenGigaSymbol(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseGenGigaSymbol",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "winresultcache" {
				componentName, err := parseWinResultCache(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseWinResultCache",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "gensymbolvalswithwinresult" {
				componentName, err := parseGenSymbolValsWithPos(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseGenSymbolValsWithPos",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "checksymbolvals" {
				componentName, err := parseCheckSymbolVals(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseCheckSymbolVals",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "positioncollection" {
				componentName, err := parsePositionCollection(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parsePositionCollection",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "chgsymbolvals" {
				componentName, err := parseChgSymbolVals(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseChgSymbolVals",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "chgsymbols" {
				componentName, err := parseChgSymbols(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseChgSymbols",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "gensymbolvalswithsymbol" {
				componentName, err := parseGenSymbolValsWithSymbol(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseGenSymbolValsWithSymbol",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else if componentType == "symbolvalswins" {
				componentName, err := parseSymbolValWins(cfg, &cell)
				if err != nil {
					goutils.Error("loadCells:parseSymbolValWins",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				mapComponentName[id] = componentName
			} else {
				goutils.Error("loadCells:ErrUnsupportedComponentType",
					slog.String("componentType", componentType),
					goutils.Err(ErrUnsupportedComponentType))

				return ErrUnsupportedComponentType
			}
		} else if shape == "edge" {
			source, err := cell.Get("source").Get("cell").String()
			if err != nil {
				goutils.Error("loadCells:edge:source:cell",
					goutils.Err(err))

				return err
			}

			sourcePort, err := cell.Get("source").Get("port").String()
			if err != nil {
				goutils.Error("loadCells:edge:source:port",
					goutils.Err(err))

				return err
			}

			target, err := cell.Get("target").Get("cell").String()
			if err != nil {
				goutils.Error("loadCells:edge:target",
					goutils.Err(err))

				return err
			}

			if sourcePort == "jump-component-groups-out" {
				ldid.add("jump", source, target)
			} else if sourcePort == "component-groups-out" {
				ldid.add("next", source, target)
			} else if sourcePort == "loop-component-groups-out" {
				ldid.add("loop", source, target)
			} else if sourcePort == "foreach-component-groups-out" {
				ldid.add("foreach", source, target)
			} else if sourcePort == "start-out" {
				lstStartID = append(lstStartID, target)
			} else {
				arr := strings.Split(sourcePort, "vals-component-groups-out-")
				if len(arr) == 2 {
					ldid.add(arr[1], source, target)
				} else {
					goutils.Error("loadCells:sourcePort",
						slog.String("sourcePort", sourcePort),
						goutils.Err(ErrUnsupportedLinkType))

					return ErrUnsupportedLinkType
				}
			}
		}
	}

	if len(lstStartID) > 0 {
		cfg.Start = mapComponentName[lstStartID[0]]
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
	// 				slog.String("label", k),
	// 				goutils.Err(ErrIvalidTriggerLabel))

	// 			return ErrIvalidTriggerLabel
	// 		}

	// 		basicWinsCfg.BeforMain = append(basicWinsCfg.BeforMain, cfg)
	// 	}

	// 	for _, k := range basicWinsCfg.AfterMainTriggerName {
	// 		cfg, isok := mapTrigger[k]
	// 		if !isok {
	// 			goutils.Error("loadCells:AfterMain",
	// 				slog.String("label", k),
	// 				goutils.Err(ErrIvalidTriggerLabel))

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
			goutils.Err(err))

		return err
	}

	betcfg := &BetConfig{
		Bet:            int(bet),
		TotalBetInWins: int(bet),
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	cfg.Bets = append(cfg.Bets, int(bet))
	// cfg.TotalBetInWins = append(cfg.TotalBetInWins, int(bet))

	err = loadCells(betcfg, betMethod.Get("graph").Get("cells"))
	if err != nil {
		goutils.Error("loadBetMethod:loadCells",
			goutils.Err(err))

		return err
	}

	cfg.MapBetConfigs[int(bet)] = betcfg

	return nil
}

func NewGame2(fn string, funcNewPlugin sgc7plugin.FuncNewPlugin, funcNewRNG FuncNewRNG) (*Game, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("NewGame2:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	return NewGame2WithData(data, funcNewPlugin, funcNewRNG)
}

// func NewGame2ForRTP(bet int, fn string, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
// 	data, err := os.ReadFile(fn)
// 	if err != nil {
// 		goutils.Error("NewGame2:ReadFile",
// 			slog.String("fn", fn),
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	return NewGame2WithData(data, funcNewPlugin)
// }

// func NewGame3(fn string, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
// 	if strings.Contains(fn, ".json") {
// 		return NewGame2(fn, funcNewPlugin)
// 	}

// 	return NewGameEx(fn, funcNewPlugin)
// }

// func NewGame3ForRTP(bet int, fn string, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
// 	if strings.Contains(fn, ".json") {
// 		return NewGame2(fn, funcNewPlugin)
// 	}

// 	return NewGameExForRTP(bet, fn, funcNewPlugin)
// }

func NewGame2WithData(data []byte, funcNewPlugin sgc7plugin.FuncNewPlugin, funcNewRNG FuncNewRNG) (*Game, error) {
	game := &Game{
		BasicGame:    sgc7game.NewBasicGame(funcNewPlugin),
		MgrComponent: NewComponentMgr(),
	}

	cfg := &Config{
		Paytables:    make(map[string]string),
		MapPaytables: make(map[string]*sgc7game.PayTables),
		Linedata:     make(map[string]string),
		MapLinedate:  make(map[string]*sgc7game.LineData),
		Reels:        make(map[string]string),
		MapReels:     make(map[string]*sgc7game.ReelsData),
		// mapConfig:         make(map[string]IComponentConfig),
		// StartComponents: make(map[int]string),
		// mapBasicConfig:    make(map[string]*BasicComponentConfig),
		mapValWeights:     make(map[string]*sgc7game.ValWeights2),
		mapReelSetWeights: make(map[string]*sgc7game.ValWeights2),
		mapStrWeights:     make(map[string]*sgc7game.ValWeights2),
		mapIntMapping:     make(map[string]*sgc7game.ValMapping2),
		MapBetConfigs:     make(map[int]*BetConfig),
		// mapBetConfig:    make(map[int]*BetDataConfig),
	}

	err := loadBasicInfo(cfg, data)
	if err != nil {
		goutils.Error("NewGame2WithData:loadBasicInfo",
			goutils.Err(err))

		return nil, err
	}

	lstPaytables, err := sonic.Get(data, "repository", "paytableList")
	if err != nil {
		goutils.Error("NewGame2WithData:Get",
			slog.String("key", "repository.paytableList"),
			goutils.Err(err))

		return nil, err
	}

	err = loadPaytables(cfg, &lstPaytables)
	if err != nil {
		goutils.Error("NewGame2WithData:loadPaytables",
			goutils.Err(err))

		return nil, err
	}

	lstOther, err := sonic.Get(data, "repository", "otherList")
	if err != nil {
		goutils.Error("NewGame2WithData:Get",
			slog.String("key", "repository.otherList"),
			goutils.Err(err))

		return nil, err
	}

	err = loadOtherList(cfg, &lstOther)
	if err != nil {
		goutils.Error("NewGame2WithData:loadOtherList",
			goutils.Err(err))

		return nil, err
	}

	// cfgGameMod := &GameModConfig{}
	// cfgGameMod.Type = "bg"
	// cfg.GameMods = append(cfg.GameMods, cfgGameMod)

	betMethodNode, err := sonic.Get(data, "betMethod")
	if err != nil {
		goutils.Error("NewGame2WithData:Get",
			slog.String("key", "betMethod"),
			goutils.Err(err))

		return nil, err
	}

	lstBetMethod, err := betMethodNode.ArrayUseNode()
	if err != nil {
		goutils.Error("NewGame2WithData:Get",
			slog.String("key", "lstBetMethod.ArrayUseNode"),
			goutils.Err(err))

		return nil, err
	}

	for _, betMethod := range lstBetMethod {
		// betMethod := lstBetMethod.Index(i) //sonic.Get(data, "betMethod", 0)
		// if err != nil {
		// 	goutils.Error("NewGame2WithData:Get",
		// 		slog.String("key", "betMethod[0]"),
		// 		goutils.Err(err))

		// 	return nil, err
		// }

		err = loadBetMethod(cfg, &betMethod)
		if err != nil {
			goutils.Error("NewGame2WithData:loadBetMethod",
				goutils.Err(err))

			return nil, err
		}
	}

	// cfg.RTP = &RTPConfig{}

	err = game.Init2(cfg, funcNewRNG)
	if err != nil {
		goutils.Error("NewGame2WithData:Init2",
			goutils.Err(err))

		return nil, err
	}

	return game, nil
}
