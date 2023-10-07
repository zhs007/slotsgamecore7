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
		}
	}

	return nil
}

func parseBasicReels(cell *ast.Node) (*BasicReelsConfig, error) {
	componentValues := cell.Get("componentValues")
	if componentValues == nil {
		goutils.Error("parseBasicReels:componentValues",
			zap.Error(ErrNoComponentValues))

		return nil, ErrNoComponentValues
	}

	buf, err := componentValues.MarshalJSON()
	if err != nil {
		goutils.Error("parseBasicReels:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	data := &basicReelsData{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseBasicReels:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	return data.build(), nil
}

func parseTriggerFeatureConfig(cell *ast.Node) (string, *TriggerFeatureConfig, error) {
	componentValues := cell.Get("componentValues")
	if componentValues == nil {
		goutils.Error("parseTriggerFeatureConfig:componentValues",
			zap.Error(ErrNoComponentValues))

		return "", nil, ErrNoComponentValues
	}

	buf, err := componentValues.MarshalJSON()
	if err != nil {
		goutils.Error("parseTriggerFeatureConfig:MarshalJSON",
			zap.Error(err))

		return "", nil, err
	}

	data := &triggerFeatureData{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseTriggerFeatureConfig:Unmarshal",
			zap.Error(err))

		return "", nil, err
	}

	return data.Label, data.build(), nil
}

func parseSymbolMulti(cell *ast.Node) (*SymbolMultiConfig, error) {
	componentValues := cell.Get("componentValues")
	if componentValues == nil {
		goutils.Error("parseSymbolMulti:componentValues",
			zap.Error(ErrNoComponentValues))

		return nil, ErrNoComponentValues
	}

	buf, err := componentValues.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolMulti:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	data := &symbolMultiData{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSymbolMulti:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	return data.build(), nil
}

func parseSymbolVal(cell *ast.Node) (*SymbolValConfig, error) {
	componentValues := cell.Get("componentValues")
	if componentValues == nil {
		goutils.Error("parseSymbolVal:componentValues",
			zap.Error(ErrNoComponentValues))

		return nil, ErrNoComponentValues
	}

	buf, err := componentValues.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolVal:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	data := &symbolValData{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSymbolVal:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	return data.build(), nil
}

func parseSymbolVal2(cell *ast.Node) (*SymbolVal2Config, error) {
	componentValues := cell.Get("componentValues")
	if componentValues == nil {
		goutils.Error("parseSymbolVal2:componentValues",
			zap.Error(ErrNoComponentValues))

		return nil, ErrNoComponentValues
	}

	buf, err := componentValues.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolVal2:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	data := &symbolVal2Data{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSymbolVal2:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	return data.build(), nil
}

func parseBasicWins(cell *ast.Node) (*BasicWinsConfig, error) {
	componentValues := cell.Get("componentValues")
	if componentValues == nil {
		goutils.Error("parseBasicWins:componentValues",
			zap.Error(ErrNoComponentValues))

		return nil, ErrNoComponentValues
	}

	buf, err := componentValues.MarshalJSON()
	if err != nil {
		goutils.Error("parseBasicWins:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	data := &basicWinsData{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseBasicWins:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	return data.build(), nil
}

func parseBookOf(cell *ast.Node) (*BookOfConfig, error) {
	componentValues := cell.Get("componentValues")
	if componentValues == nil {
		goutils.Error("parseBookOf:componentValues",
			zap.Error(ErrNoComponentValues))

		return nil, ErrNoComponentValues
	}

	buf, err := componentValues.MarshalJSON()
	if err != nil {
		goutils.Error("parseBookOf:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	data := &bookOfData{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseBookOf:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	return data.build(), nil
}

func loadCells(cfg *Config, bet int, cells *ast.Node) error {
	linkScene := [][]string{}
	linkOtherScene := [][]string{}
	linkComponent := [][]string{}
	lstStart := []string{}
	mapTrigger := make(map[string]*TriggerFeatureConfig)
	mapTriggerID := make(map[string]*TriggerFeatureConfig)
	lstBasicWins := []*BasicWinsConfig{}

	lst, err := cells.ArrayUseNode()
	if err != nil {
		goutils.Error("loadCells:ArrayUseNode",
			zap.Error(err))

		return err
	}

	startid := ""

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

		if shape == "custom-node-width-start" {
			startid = id
		} else if shape == "custom-node" {
			componentType, err := cell.Get("label").String()
			if err != nil {
				goutils.Error("loadCells:Get:label",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}

			if componentType == "basicReels" {
				componentCfg, err := parseBasicReels(&cell)
				if err != nil {
					goutils.Error("loadCells:parseBasicReels",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				cfg.mapConfig[id] = componentCfg
				cfg.mapBasicConfig[id] = &componentCfg.BasicComponentConfig

				ccfg := &ComponentConfig{
					Name: id,
					Type: "basicReels",
				}

				cfg.GameMods[0].Components = append(cfg.GameMods[0].Components, ccfg)
			} else if componentType == "basicWins" {
				componentCfg, err := parseBasicWins(&cell)
				if err != nil {
					goutils.Error("loadCells:parseBasicWins",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				cfg.mapConfig[id] = componentCfg
				cfg.mapBasicConfig[id] = &componentCfg.BasicComponentConfig

				ccfg := &ComponentConfig{
					Name: id,
					Type: "basicWins",
				}

				cfg.GameMods[0].Components = append(cfg.GameMods[0].Components, ccfg)

				lstBasicWins = append(lstBasicWins, componentCfg)
			} else if componentType == "symbolMulti" {
				componentCfg, err := parseSymbolMulti(&cell)
				if err != nil {
					goutils.Error("loadCells:parseSymbolMulti",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				cfg.mapConfig[id] = componentCfg
				cfg.mapBasicConfig[id] = &componentCfg.BasicComponentConfig

				ccfg := &ComponentConfig{
					Name: id,
					Type: "symbolMulti",
				}

				cfg.GameMods[0].Components = append(cfg.GameMods[0].Components, ccfg)
			} else if componentType == "symbolVal" {
				componentCfg, err := parseSymbolVal(&cell)
				if err != nil {
					goutils.Error("loadCells:parseSymbolVal",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				cfg.mapConfig[id] = componentCfg
				cfg.mapBasicConfig[id] = &componentCfg.BasicComponentConfig

				ccfg := &ComponentConfig{
					Name: id,
					Type: "symbolVal",
				}

				cfg.GameMods[0].Components = append(cfg.GameMods[0].Components, ccfg)
			} else if componentType == "symbolVal2" {
				componentCfg, err := parseSymbolVal2(&cell)
				if err != nil {
					goutils.Error("loadCells:parseSymbolVal2",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				cfg.mapConfig[id] = componentCfg
				cfg.mapBasicConfig[id] = &componentCfg.BasicComponentConfig

				ccfg := &ComponentConfig{
					Name: id,
					Type: "symbolVal2",
				}

				cfg.GameMods[0].Components = append(cfg.GameMods[0].Components, ccfg)
			} else if componentType == "bookOf" {
				componentCfg, err := parseBookOf(&cell)
				if err != nil {
					goutils.Error("loadCells:parseBookOf",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				cfg.mapConfig[id] = componentCfg
				cfg.mapBasicConfig[id] = &componentCfg.BasicComponentConfig

				ccfg := &ComponentConfig{
					Name: id,
					Type: "bookOf",
				}

				cfg.GameMods[0].Components = append(cfg.GameMods[0].Components, ccfg)
			} else if componentType == "basicWins-trigger" {
				triggerName, triggerCfg, err := parseTriggerFeatureConfig(&cell)
				if err != nil {
					goutils.Error("loadCells:parseTriggerFeatureConfig",
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				mapTrigger[triggerName] = triggerCfg

				mapTriggerID[id] = triggerCfg
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

			if source == startid {
				lstStart = append(lstStart, target)
			} else {
				if strings.Contains(sourcePort, "component") {
					linkComponent = append(linkComponent, []string{source, target})
				} else if strings.Contains(sourcePort, "-Scenes-") {
					linkScene = append(linkScene, []string{source, target})
				} else if strings.Contains(sourcePort, "-OtherScenes-") {
					linkOtherScene = append(linkOtherScene, []string{source, target})
				}
			}
		}
	}

	if len(lstStart) > 0 {
		cfg.StartComponents[bet] = lstStart[0]
	}

	for _, arr := range linkComponent {
		basicCfg, isok := cfg.mapBasicConfig[arr[0]]
		if isok {
			basicCfg.DefaultNextComponent = arr[1]
		}
	}

	for _, arr := range linkScene {
		sourceCfg, isok0 := cfg.mapBasicConfig[arr[0]]
		if isok0 {
			sourceCfg.TagScenes = append(sourceCfg.TagScenes, arr[0])
		}

		targetCfg := cfg.mapBasicConfig[arr[1]]
		if targetCfg != nil {
			targetCfg.TargetScene = arr[0]
		}

		triggerCfg := mapTriggerID[arr[1]]
		if triggerCfg != nil {
			triggerCfg.TargetScene = arr[0]
		}
	}

	for _, arr := range linkOtherScene {
		sourceCfg := cfg.mapBasicConfig[arr[0]]
		if sourceCfg != nil {
			sourceCfg.TagOtherScenes = append(sourceCfg.TagOtherScenes, arr[0])
		}

		targetCfg := cfg.mapBasicConfig[arr[1]]
		if targetCfg != nil {
			targetCfg.TargetOtherScene = arr[0]
		}
	}

	for _, basicWinsCfg := range lstBasicWins {
		for _, k := range basicWinsCfg.BeforMainTriggerName {
			cfg, isok := mapTrigger[k]
			if !isok {
				goutils.Error("loadCells:BeforMain",
					zap.String("label", k),
					zap.Error(ErrIvalidTriggerLabel))

				return ErrIvalidTriggerLabel
			}

			basicWinsCfg.BeforMain = append(basicWinsCfg.BeforMain, cfg)
		}

		for _, k := range basicWinsCfg.AfterMainTriggerName {
			cfg, isok := mapTrigger[k]
			if !isok {
				goutils.Error("loadCells:AfterMain",
					zap.String("label", k),
					zap.Error(ErrIvalidTriggerLabel))

				return ErrIvalidTriggerLabel
			}

			basicWinsCfg.AfterMain = append(basicWinsCfg.AfterMain, cfg)
		}
	}

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

func NewGame3(fn string, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
	if strings.Contains(fn, ".json") {
		return NewGame2(fn, funcNewPlugin)
	}

	return NewGameEx(fn, funcNewPlugin)
}

func NewGame2WithData(data []byte, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
	game := &Game{
		BasicGame:    sgc7game.NewBasicGame(funcNewPlugin),
		MgrComponent: NewComponentMgr(),
	}

	cfg := &Config{
		Paytables:       make(map[string]string),
		MapPaytables:    make(map[string]*sgc7game.PayTables),
		Linedata:        make(map[string]string),
		MapLinedate:     make(map[string]*sgc7game.LineData),
		Reels:           make(map[string]string),
		MapReels:        make(map[string]*sgc7game.ReelsData),
		mapConfig:       make(map[string]any),
		StartComponents: make(map[int]string),
		mapBasicConfig:  make(map[string]*BasicComponentConfig),
		mapValWeights:   make(map[string]*sgc7game.ValWeights2),
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
