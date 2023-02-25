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

func (mgr *ComponentMgr) NewComponent(component string) IComponent {
	funcNew, isok := mgr.MapComponent[component]
	if isok {
		return funcNew()
	}

	goutils.Error("ComponentMgr.NewComponent",
		zap.String("component", component),
		zap.Error(ErrInvalidComponent))

	return nil
}

func NewComponentMgr() *ComponentMgr {
	mgr := &ComponentMgr{
		MapComponent: make(map[string]FuncNewComponent),
	}

	mgr.Reg("basicReels", NewBasicReels)

	return mgr
}
