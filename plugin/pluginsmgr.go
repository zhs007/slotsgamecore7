package sgc7plugin

import "sync"

// FuncNewPlugin - new a IPlugin
type FuncNewPlugin func() IPlugin

// PluginsMgr - plugins manager
type PluginsMgr struct {
	mutex         sync.Mutex
	plugins       []IPlugin
	funcNewPlugin FuncNewPlugin
}

// NewPluginsMgr - new a PluginsMgr
func NewPluginsMgr(funcNewPlugin FuncNewPlugin) *PluginsMgr {
	return &PluginsMgr{
		funcNewPlugin: funcNewPlugin,
	}
}

// NewPlugin - new a Plugin
func (mgr *PluginsMgr) NewPlugin() IPlugin {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if len(mgr.plugins) > 0 {
		plugin := mgr.plugins[0]
		plugin.ClearUsedRngs()

		mgr.plugins = mgr.plugins[1:]

		return plugin
	}

	plugin := mgr.funcNewPlugin()
	plugin.ClearUsedRngs()

	return plugin
}

// FreePlugin - free a Plugin
func (mgr *PluginsMgr) FreePlugin(plugin IPlugin) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mgr.plugins = append(mgr.plugins, plugin)
}
