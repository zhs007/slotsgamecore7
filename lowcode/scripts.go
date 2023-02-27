package lowcode

type OnInitScriptFunc func(scriptFunc *ScriptFunc, script string) error
type OnRunScriptFunc func(gameProp *GameProperty, scriptFunc *ScriptFunc) error

type ScriptFunc struct {
	IntParams []int           `json:"intParams"`
	StrParams []string        `json:"strParams"`
	OnRun     OnRunScriptFunc `json:"-"`
}

type ScriptFuncMgr struct {
	MapFuncs map[string]*ScriptFunc
}
