package lowcode

import (
	"fmt"
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

				// return ErrIvalidDefaultScene
			} else {
				cfg.DefaultScene = scene
			}
		}
	}

	return nil
}

func parsePaytable2(n *ast.Node) (*sgc7game.PayTables, error) {
	if n == nil {
		goutils.Error("parsePaytable2",
			goutils.Err(ErrIvalidPayTables))

		return nil, ErrIvalidPayTables
	}

	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parsePaytable2:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	dataPaytables := []map[string]string{}

	err = sonic.Unmarshal(buf, &dataPaytables)
	if err != nil {
		goutils.Error("parsePaytable2:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	paytables := &sgc7game.PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	for _, node := range dataPaytables {
		code, err := goutils.String2Int64(strings.TrimSpace(node["Code"]))
		if err != nil {
			goutils.Error("parsePaytable2:String2Int64:code",
				slog.String("code", node["Code"]),
				goutils.Err(err))

			return nil, err
		}

		arr := []int{}

		n := 1
		for {
			key := fmt.Sprintf("X%v", n)
			v, isok := node[key]
			if !isok {
				break
			}

			i64, err := goutils.String2Int64(v)
			if err != nil {
				goutils.Error("parsePaytable2:String2Int64:X",
					slog.String("X", v),
					goutils.Err(err))

				return nil, err
			}

			arr = append(arr, int(i64))

			n++
		}

		paytables.MapPay[int(code)] = arr

		paytables.MapSymbols[node["Symbol"]] = int(code)
	}

	return paytables, nil
}

func loadPaytables(cfg *Config, paytableData *ast.Node) error {
	paytables, err := parsePaytable2(paytableData)
	if err != nil {
		goutils.Error("loadPaytables:parsePaytable2",
			goutils.Err(err))

		return err
	}

	cfg.Paytables["default"] = "default"
	cfg.MapPaytables["default"] = paytables

	cfg.DefaultPaytables = "default"

	return nil
}

func getLineData2Width(dataLines []map[string]string) int {
	w := 1
	for _, v := range dataLines {
		x := 1
		for ; x < 99; x++ {
			_, isok := v[fmt.Sprintf("R%v", x)]
			if !isok {
				x--

				break
			}
		}

		if w < x {
			w = x
		}
	}

	return w
}

func parseLineData2(n *ast.Node) (*sgc7game.LineData, error) {
	if n == nil {
		goutils.Error("parseLineData2",
			goutils.Err(ErrIvalidReels))

		return nil, ErrIvalidReels
	}

	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseLineData2:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	dataLines := []map[string]string{}

	err = sonic.Unmarshal(buf, &dataLines)
	if err != nil {
		goutils.Error("parseLineData2:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	w := getLineData2Width(dataLines)
	ld := &sgc7game.LineData{}

	for _, linedata := range dataLines {
		arr := []int{}
		for x := 1; x <= w; x++ {
			i64, err := goutils.String2Int64(strings.TrimSpace(linedata[fmt.Sprintf("R%v", x)]))
			if err != nil {
				goutils.Error("parseLineData2:String2Int64",
					goutils.Err(err))

				return nil, err
			}

			arr = append(arr, int(i64))
		}

		ld.Lines = append(ld.Lines, arr)
	}

	return ld, nil
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

func getReel2Width(dataReels []map[string]string) int {
	w := 1
	for _, v := range dataReels {
		x := 1
		for ; x < 99; x++ {
			_, isok := v[fmt.Sprintf("R%v", x)]
			if !isok {
				x--

				break
			}
		}

		if w < x {
			w = x
		}
	}

	return w
}

func parseReel2(dataReels []map[string]string, x int, paytables *sgc7game.PayTables) ([]int, error) {
	arr := []int{}
	for y, obj := range dataReels {
		str, isok := obj[fmt.Sprintf("R%v", x)]
		if !isok {
			return arr, nil
		}

		str = strings.TrimSpace(str)
		if str == "" {
			return arr, nil
		}

		sc, isok := paytables.MapSymbols[str]
		if !isok {
			goutils.Error("parseReel2:MapSymbols",
				slog.Int("x", x),
				slog.Int("y", y),
				goutils.Err(ErrIvalidSymbolInReels))

			return nil, ErrIvalidSymbolInReels
		}

		arr = append(arr, sc)
	}

	return arr, nil
}

func parseReels2(n *ast.Node, paytables *sgc7game.PayTables) (*sgc7game.ReelsData, error) {
	if n == nil {
		goutils.Error("parseReels2",
			goutils.Err(ErrIvalidReels))

		return nil, ErrIvalidReels
	}

	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseReels2:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	dataReels := []map[string]string{}

	err = sonic.Unmarshal(buf, &dataReels)
	if err != nil {
		goutils.Error("parseReels2:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	w := getReel2Width(dataReels)

	reelsd := &sgc7game.ReelsData{}

	for x := 1; x <= w; x++ {
		arr, err := parseReel2(dataReels, x, paytables)
		if err != nil {
			goutils.Error("parseReels2:parseReel2",
				slog.Int("x", x),
				goutils.Err(err))

			return nil, err
		}

		reelsd.Reels = append(reelsd.Reels, arr)
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
			if v.Get("fileJson") != nil {
				ld, err := parseLineData(v.Get("fileJson"), cfg.Width)
				if err != nil {
					goutils.Error("loadOtherList:parseLineData",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.Linedata[name] = name
				cfg.MapLinedate[name] = ld
			} else {
				ld, err := parseLineData2(v.Get("excelJson"))
				if err != nil {
					goutils.Error("loadOtherList:parseLineData2",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.Linedata[name] = name
				cfg.MapLinedate[name] = ld
			}

			if len(cfg.Linedata) == 1 {
				cfg.DefaultLinedata = name
			}
		} else if t == "Reels" {
			if v.Get("fileJson") != nil {
				rd, err := parseReels(v.Get("fileJson"), cfg.GetDefaultPaytables())
				if err != nil {
					goutils.Error("loadOtherList:parseReels",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.Reels[name] = name
				cfg.MapReels[name] = rd
			} else {
				rd, err := parseReels2(v.Get("excelJson"), cfg.GetDefaultPaytables())
				if err != nil {
					goutils.Error("loadOtherList:parseReels2",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.Reels[name] = name
				cfg.MapReels[name] = rd
			}
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
			if v.Get("fileJson") != nil {
				vw2, err := parseReelSetWeights(v.Get("fileJson"))
				if err != nil {
					goutils.Error("loadOtherList:parseReelSetWeights",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.mapReelSetWeights[name] = vw2
			} else {
				vw2, err := parseReelSetWeights2(v.Get("excelJson"))
				if err != nil {
					goutils.Error("loadOtherList:parseReelSetWeights2",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.mapReelSetWeights[name] = vw2
			}
		} else if t == "SymbolWeight" {
			if v.Get("fileJson") != nil {
				vw2, err := parseSymbolWeights(v.Get("fileJson"), cfg.GetDefaultPaytables())
				if err != nil {
					goutils.Error("loadOtherList:parseSymbolWeights",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.mapValWeights[name] = vw2
			} else {
				vw2, err := parseSymbolWeights2(v.Get("excelJson"), cfg.GetDefaultPaytables())
				if err != nil {
					goutils.Error("loadOtherList:parseSymbolWeights2",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.mapValWeights[name] = vw2
			}
		} else if t == "IntValMappingFile" {
			if v.Get("fileJson") != nil {
				vm2, err := parseIntValMappingFile(v.Get("fileJson"))
				if err != nil {
					goutils.Error("loadOtherList:parseSymbolWeights",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.mapIntMapping[name] = vm2
			} else {
				vm2, err := parseIntValMappingFile2(v.Get("excelJson"))
				if err != nil {
					goutils.Error("loadOtherList:parseIntValMappingFile2",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.mapIntMapping[name] = vm2
			}
		} else if t == "StringValWeight" {
			if v.Get("fileJson") != nil {
				vm2, err := parseStrWeights(v.Get("fileJson"))
				if err != nil {
					goutils.Error("loadOtherList:parseStrWeights",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.mapStrWeights[name] = vm2
			} else {
				vm2, err := parseStrWeights2(v.Get("excelJson"))
				if err != nil {
					goutils.Error("loadOtherList:parseStrWeights2",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.mapStrWeights[name] = vm2
			}
		} else if t == "IntValWeight" {
			if v.Get("fileJson") != nil {
				vm2, err := parseIntValWeights(v.Get("fileJson"))
				if err != nil {
					goutils.Error("loadOtherList:parseIntValWeights",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.mapValWeights[name] = vm2
			} else {
				vm2, err := parseIntValWeights2(v.Get("excelJson"))
				if err != nil {
					goutils.Error("loadOtherList:parseIntValWeights2",
						slog.Int("i", i),
						goutils.Err(err))

					return err
				}

				cfg.mapValWeights[name] = vm2
			}
		} else if t == "StringValMappingFile" {
			vm2, err := parseStringValMappingFile2(v.Get("excelJson"))
			if err != nil {
				goutils.Error("loadOtherList:parseStringValMappingFile2",
					slog.Int("i", i),
					goutils.Err(err))

				return err
			}

			cfg.mapIntMapping[name] = vm2

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

		if shape == "custom-node" {
			componentType, err := cell.Get("label").String()
			if err != nil {
				goutils.Error("loadCells:Get:label",
					slog.Int("i", i),
					goutils.Err(err))

				return err
			}

			componentType = strings.ToLower(componentType)

			componentName, err := gJsonMgr.LoadComponent(componentType, cfg, &cell)
			if err != nil {
				goutils.Error("loadCells:LoadComponent",
					slog.Int("i", i),
					goutils.Err(err))

				return err
			}

			mapComponentName[id] = componentName
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
					arr1 := strings.Split(sourcePort, "reelTrigger-component-groups-out-")
					if len(arr1) == 2 {
						ldid.add(arr1[1], source, target)
					} else {
						goutils.Error("loadCells:sourcePort",
							slog.String("sourcePort", sourcePort),
							goutils.Err(ErrUnsupportedLinkType))

						return ErrUnsupportedLinkType
					}
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

	nodeTotalBetInWins := betMethod.Get("totalBetInWins")
	if nodeTotalBetInWins != nil {
		totalBetInWins, err := nodeTotalBetInWins.Int64()
		if err != nil {
			goutils.Error("loadBetMethod:Get:totalBetInWins",
				goutils.Err(err))

			return err
		}

		betcfg.TotalBetInWins = int(totalBetInWins)
		cfg.TotalBetInWins = append(cfg.TotalBetInWins, int(totalBetInWins))
	} else {
		cfg.TotalBetInWins = append(cfg.TotalBetInWins, int(bet))
	}

	err = loadCells(betcfg, betMethod.Get("graph").Get("cells"))
	if err != nil {
		goutils.Error("loadBetMethod:loadCells",
			goutils.Err(err))

		return err
	}

	cfg.MapBetConfigs[int(bet)] = betcfg

	return nil
}

func NewGame2(fn string, funcNewPlugin sgc7plugin.FuncNewPlugin, funcNewRNG FuncNewRNG, funcNewFeatureLevel FuncNewFeatureLevel) (*Game, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("NewGame2:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	return NewGame2WithData(data, funcNewPlugin, funcNewRNG, funcNewFeatureLevel)
}

func NewGame2WithData(data []byte, funcNewPlugin sgc7plugin.FuncNewPlugin, funcNewRNG FuncNewRNG, funcNewFeatureLevel FuncNewFeatureLevel) (*Game, error) {
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
		mapValWeights:     make(map[string]*sgc7game.ValWeights2),
		mapReelSetWeights: make(map[string]*sgc7game.ValWeights2),
		mapStrWeights:     make(map[string]*sgc7game.ValWeights2),
		mapIntMapping:     make(map[string]*sgc7game.ValMapping2),
		MapBetConfigs:     make(map[int]*BetConfig),
	}

	err := loadBasicInfo(cfg, data)
	if err != nil {
		goutils.Error("NewGame2WithData:loadBasicInfo",
			goutils.Err(err))

		return nil, err
	}

	paytableData, err := sonic.Get(data, "repository", "paytableData")
	if err != nil {
		goutils.Error("NewGame2WithData:Get",
			slog.String("key", "repository.paytableList"),
			goutils.Err(err))

		return nil, err
	}

	err = loadPaytables(cfg, &paytableData)
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

	err = game.Init2(cfg, funcNewRNG, funcNewFeatureLevel)
	if err != nil {
		goutils.Error("NewGame2WithData:Init2",
			goutils.Err(err))

		return nil, err
	}

	return game, nil
}
