package lowcode

type OnInitScriptFunc func(scriptFunc *ScriptFunc, script string) error
type OnRunScriptFunc func(scriptFunc *ScriptFunc) error

type ScriptFunc struct {
	IntParams []int
	StrParams []string
	OnRun     OnRunScriptFunc
}

type ScriptFuncMgr struct {
	MapFuncs map[string]*ScriptFunc
}
