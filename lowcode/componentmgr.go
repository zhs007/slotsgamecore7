package lowcode

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type ComponentMgr struct {
	MapComponent map[string]FuncNewComponent
}

func (mgr *ComponentMgr) Reg(component string, funcNew FuncNewComponent) {
	mgr.MapComponent[component] = funcNew
}

func (mgr *ComponentMgr) NewComponent(cfgComponent *ComponentConfig) IComponent {
	funcNew, isok := mgr.MapComponent[cfgComponent.Type]
	if isok {
		return funcNew(cfgComponent.Name)
	}

	goutils.Error("ComponentMgr.NewComponent",
		zap.String("component", cfgComponent.Type),
		zap.Error(ErrInvalidComponent))

	return nil
}

func NewComponentMgr() *ComponentMgr {
	mgr := &ComponentMgr{
		MapComponent: make(map[string]FuncNewComponent),
	}

	mgr.Reg("basicReels", NewBasicReels)
	mgr.Reg("mystery", NewMystery)
	mgr.Reg("basicWins", NewBasicWins)
	mgr.Reg("lightning", NewLightning)
	mgr.Reg("multiLevelReels", NewMultiLevelReels)
	mgr.Reg("collector", NewCollector)
	mgr.Reg("multiLevelMystery", NewMultiLevelMystery)
	mgr.Reg("bookOf", NewBookOf)
	mgr.Reg("symbolMulti", NewSymbolMulti)
	mgr.Reg("symbolVal", NewSymbolVal)
	mgr.Reg("symbolValWins", NewSymbolValWins)
	mgr.Reg("symbolVal2", NewSymbolVal2)
	mgr.Reg("overlaySymbol", NewOverlaySymbol)
	mgr.Reg("multiFG", NewMultiFG)
	mgr.Reg("reelSetMystery", NewReelSetMystery)
	mgr.Reg("weightTrigger", NewWeightTrigger)

	return mgr
}
