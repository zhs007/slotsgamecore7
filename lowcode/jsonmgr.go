package lowcode

import (
	"log/slog"
	"strings"

	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
)

type FuncLoadComponentInJson func(gamecfg *BetConfig, cell *ast.Node) (string, error)

type JsonMgr struct {
	mapLoadComponent map[string]FuncLoadComponentInJson
}

func (mgr *JsonMgr) RegLoadComponent(typename string, loader FuncLoadComponentInJson) {
	mgr.mapLoadComponent[typename] = loader
}

func (mgr *JsonMgr) LoadComponent(componentType string, gamecfg *BetConfig, cell *ast.Node) (string, error) {
	loader, isok := mgr.mapLoadComponent[componentType]
	if isok {
		componentName, err := loader(gamecfg, cell)
		if err != nil {
			goutils.Error("JsonMgr.LoadComponent:loader",
				slog.String("componentType", componentType),
				goutils.Err(err))

			return "", err
		}

		return componentName, nil
	}

	goutils.Error("JsonMgr.LoadComponent:ErrUnsupportedComponentType",
		slog.String("componentType", componentType),
		goutils.Err(ErrUnsupportedComponentType))

	return "", ErrUnsupportedComponentType
}

var gJsonMgr *JsonMgr

func init() {
	gJsonMgr = &JsonMgr{
		mapLoadComponent: make(map[string]FuncLoadComponentInJson),
	}

	gJsonMgr.RegLoadComponent("weightreels", parseWeightReels)
	gJsonMgr.RegLoadComponent("basicreels", parseBasicReels)
	gJsonMgr.RegLoadComponent("scattertrigger", parseScatterTrigger)
	gJsonMgr.RegLoadComponent("linestrigger", parseLinesTrigger)
	gJsonMgr.RegLoadComponent("waystrigger", parseWaysTrigger)
	gJsonMgr.RegLoadComponent("movesymbol", parseMoveSymbol)
	gJsonMgr.RegLoadComponent("respin", parseRespin)
	gJsonMgr.RegLoadComponent("symbolcollection", parseSymbolCollection2)
	gJsonMgr.RegLoadComponent("removesymbols", parseRemoveSymbols)
	gJsonMgr.RegLoadComponent("dropdownsymbols", parseDropDownSymbols)
	gJsonMgr.RegLoadComponent("refillsymbols", parseRefillSymbols)
	gJsonMgr.RegLoadComponent("collector", parseCollector)
	gJsonMgr.RegLoadComponent("queuebranch", parseQueueBranch)
	gJsonMgr.RegLoadComponent("delayqueue", parseQueueBranch)
	gJsonMgr.RegLoadComponent("replacesymbolgroup", parseReplaceSymbolGroup)
	gJsonMgr.RegLoadComponent("rollsymbol", parseRollSymbol)
	gJsonMgr.RegLoadComponent("mask", parseMask)
	gJsonMgr.RegLoadComponent("replacereelwithmask", parseReplaceReelWithMask)
	gJsonMgr.RegLoadComponent("piggybank", parsePiggyBank)
	gJsonMgr.RegLoadComponent("addsymbols", parseAddSymbols)
	gJsonMgr.RegLoadComponent("intvalmapping", parseIntValMapping)
	gJsonMgr.RegLoadComponent("weightbranch", parseWeightBranch)
	gJsonMgr.RegLoadComponent("clustertrigger", parseClusterTrigger)
	gJsonMgr.RegLoadComponent("gengigasymbol", parseGenGigaSymbol)
	gJsonMgr.RegLoadComponent("winresultcache", parseWinResultCache)
	gJsonMgr.RegLoadComponent("gensymbolvalswithwinresult", parseGenSymbolValsWithPos)
	gJsonMgr.RegLoadComponent("checksymbolvals", parseCheckSymbolVals)
	gJsonMgr.RegLoadComponent("positioncollection", parsePositionCollection)
	gJsonMgr.RegLoadComponent("chgsymbolvals", parseChgSymbolVals)
	gJsonMgr.RegLoadComponent("chgsymbols", parseChgSymbols)
	gJsonMgr.RegLoadComponent("gensymbolvalswithsymbol", parseGenSymbolValsWithSymbol)
	gJsonMgr.RegLoadComponent("symbolvalswins", parseSymbolValWins)
	gJsonMgr.RegLoadComponent("rebuildreelindex", parseRebuildReelIndex)
	gJsonMgr.RegLoadComponent("gensymbolvals", parseGenSymbolVals)
	gJsonMgr.RegLoadComponent("rebuildsymbols", parseRebuildSymbols)
	gJsonMgr.RegLoadComponent("rollnumber", parseRollNumber)
	gJsonMgr.RegLoadComponent("controllerworker", parseControllerWorker)
	gJsonMgr.RegLoadComponent("catchsymbols", parseCatchSymbols)
	gJsonMgr.RegLoadComponent(strings.ToLower(BurstSymbolsTypeName), parseBurstSymbols)
	gJsonMgr.RegLoadComponent(strings.ToLower(WinResultModifierTypeName), parseWinResultModifier)
	gJsonMgr.RegLoadComponent(strings.ToLower(ReelTriggerTypeName), parseReelTrigger)
	gJsonMgr.RegLoadComponent(strings.ToLower(JackpotTypeName), parseJackpot)
	gJsonMgr.RegLoadComponent(strings.ToLower(CheckValTypeName), parseCheckVal)
	gJsonMgr.RegLoadComponent(strings.ToLower(AdjacentPayTriggerTypeName), parseAdjacentPayTrigger)
	gJsonMgr.RegLoadComponent(strings.ToLower(WinResultMultiTypeName), parseWinResultMulti)
	gJsonMgr.RegLoadComponent(strings.ToLower(WeightReels2TypeName), parseWeightReels2)
	gJsonMgr.RegLoadComponent(strings.ToLower(WinResultModifierExTypeName), parseWinResultModifierEx)
	gJsonMgr.RegLoadComponent(strings.ToLower(RandomMoveSymbolsTypeName), parseRandomMoveSymbols)
	gJsonMgr.RegLoadComponent(strings.ToLower(GenPositionCollectionTypeName), parseGenPositionCollection)
	gJsonMgr.RegLoadComponent(strings.ToLower(FeatureBarTypeName), parseFeatureBar)
	gJsonMgr.RegLoadComponent(strings.ToLower(BombTypeName), parseBomb)
	gJsonMgr.RegLoadComponent(strings.ToLower(SumSymbolValsTypeName), parseSumSymbolVals)
}
