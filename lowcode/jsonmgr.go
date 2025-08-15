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

	gJsonMgr.RegLoadComponent(strings.ToLower(WeightReelsTypeName), parseWeightReels)
	gJsonMgr.RegLoadComponent(strings.ToLower(BasicReelsTypeName), parseBasicReels)
	gJsonMgr.RegLoadComponent(strings.ToLower(ScatterTriggerTypeName), parseScatterTrigger)
	gJsonMgr.RegLoadComponent(strings.ToLower(LinesTriggerTypeName), parseLinesTrigger)
	gJsonMgr.RegLoadComponent(strings.ToLower(WaysTriggerTypeName), parseWaysTrigger)
	gJsonMgr.RegLoadComponent(strings.ToLower(MoveSymbolTypeName), parseMoveSymbol)
	gJsonMgr.RegLoadComponent(strings.ToLower(RespinTypeName), parseRespin)
	gJsonMgr.RegLoadComponent(strings.ToLower(SymbolCollection2TypeName), parseSymbolCollection2)
	gJsonMgr.RegLoadComponent(strings.ToLower(RemoveSymbolsTypeName), parseRemoveSymbols)
	gJsonMgr.RegLoadComponent(strings.ToLower(DropDownSymbolsTypeName), parseDropDownSymbols)
	gJsonMgr.RegLoadComponent(strings.ToLower(RefillSymbolsTypeName), parseRefillSymbols)
	gJsonMgr.RegLoadComponent(strings.ToLower(CollectorTypeName), parseCollector)
	gJsonMgr.RegLoadComponent(strings.ToLower(QueueBranchTypeName), parseQueueBranch)
	gJsonMgr.RegLoadComponent("delayqueue", parseQueueBranch)
	gJsonMgr.RegLoadComponent(strings.ToLower(ReplaceSymbolGroupTypeName), parseReplaceSymbolGroup)
	gJsonMgr.RegLoadComponent(strings.ToLower(RollSymbolTypeName), parseRollSymbol)
	gJsonMgr.RegLoadComponent(strings.ToLower(MaskTypeName), parseMask)
	gJsonMgr.RegLoadComponent(strings.ToLower(ReplaceReelWithMaskTypeName), parseReplaceReelWithMask)
	gJsonMgr.RegLoadComponent(strings.ToLower(SymbolExplanderTypeName), parseSymbolExplander)
	gJsonMgr.RegLoadComponent(strings.ToLower(PiggyBankTypeName), parsePiggyBank)
	gJsonMgr.RegLoadComponent(strings.ToLower(AddSymbolsTypeName), parseAddSymbols)
	gJsonMgr.RegLoadComponent(strings.ToLower(IntValMappingTypeName), parseIntValMapping)
	gJsonMgr.RegLoadComponent(strings.ToLower(WeightBranchTypeName), parseWeightBranch)
	gJsonMgr.RegLoadComponent(strings.ToLower(ClusterTriggerTypeName), parseClusterTrigger)
	gJsonMgr.RegLoadComponent(strings.ToLower(GenGigaSymbolTypeName), parseGenGigaSymbol)
	gJsonMgr.RegLoadComponent(strings.ToLower(WinResultCacheTypeName), parseWinResultCache)
	gJsonMgr.RegLoadComponent(strings.ToLower(GenSymbolValsWithPosTypeName), parseGenSymbolValsWithPos)
	gJsonMgr.RegLoadComponent(strings.ToLower(CheckSymbolValsTypeName), parseCheckSymbolVals)
	gJsonMgr.RegLoadComponent(strings.ToLower(PositionCollectionTypeName), parsePositionCollection)
	gJsonMgr.RegLoadComponent(strings.ToLower(ChgSymbolValsTypeName), parseChgSymbolVals)
	gJsonMgr.RegLoadComponent(strings.ToLower(ChgSymbolsTypeName), parseChgSymbols)
	gJsonMgr.RegLoadComponent(strings.ToLower(GenSymbolValsWithSymbolTypeName), parseGenSymbolValsWithSymbol)
	gJsonMgr.RegLoadComponent(strings.ToLower(SymbolValWinsTypeName), parseSymbolValWins)
	gJsonMgr.RegLoadComponent("symbolvalswins", parseSymbolValWins)
	gJsonMgr.RegLoadComponent(strings.ToLower(RebuildReelIndexTypeName), parseRebuildReelIndex)
	gJsonMgr.RegLoadComponent(strings.ToLower(GenSymbolValsTypeName), parseGenSymbolVals)
	gJsonMgr.RegLoadComponent(strings.ToLower(RebuildSymbolsTypeName), parseRebuildSymbols)
	gJsonMgr.RegLoadComponent(strings.ToLower(RollNumberTypeName), parseRollNumber)
	gJsonMgr.RegLoadComponent(strings.ToLower(ControllerWorkerTypeName), parseControllerWorker)
	gJsonMgr.RegLoadComponent(strings.ToLower(CatchSymbolsTypeName), parseCatchSymbols)
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
	gJsonMgr.RegLoadComponent(strings.ToLower(TreasureChestTypeName), parseTreasureChest)
	gJsonMgr.RegLoadComponent(strings.ToLower(MoveSymbols2TypeName), parseMoveSymbols2)
	gJsonMgr.RegLoadComponent(strings.ToLower(GenSymbolVals2TypeName), parseGenSymbolVals2)
	gJsonMgr.RegLoadComponent(strings.ToLower(MergePositionCollectionTypeName), parseMergePositionCollection)
	gJsonMgr.RegLoadComponent(strings.ToLower(FeatureBar2TypeName), parseFeatureBar2)
	gJsonMgr.RegLoadComponent(strings.ToLower(ChgSymbols2TypeName), parseChgSymbols2)
	gJsonMgr.RegLoadComponent(strings.ToLower(FeaturePickTypeName), parseFeaturePick)
	gJsonMgr.RegLoadComponent(strings.ToLower(ReelsCollectorTypeName), parseReelsCollector)
	gJsonMgr.RegLoadComponent(strings.ToLower(FlowDownSymbolsTypeName), parseFlowDownSymbols)
	gJsonMgr.RegLoadComponent(strings.ToLower(HoldAndWinTypeName), parseHoldAndWin)
	gJsonMgr.RegLoadComponent(strings.ToLower(ChgSymbolsInReelsTypeName), parseChgSymbolsInReels)
	gJsonMgr.RegLoadComponent(strings.ToLower(CascadingRegulatorTypeName), parseCascadingRegulator)
	gJsonMgr.RegLoadComponent(strings.ToLower(WinResultLimiterTypeName), parseWinResultLimiter)
	gJsonMgr.RegLoadComponent(strings.ToLower(SymbolValsSPTypeName), parseSymbolValsSP)
	gJsonMgr.RegLoadComponent(strings.ToLower(Collector2TypeName), parseCollector2)
	gJsonMgr.RegLoadComponent(strings.ToLower(DropDownSymbols2TypeName), parseDropDownSymbols2)
	gJsonMgr.RegLoadComponent(strings.ToLower(HoldAndRespinReelsTypeName), parseHoldAndRespinReels)
}
