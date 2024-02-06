package lowcode

const (
	CCVReelSet            string = "reelSet"            // 可以修改配置项里的 reelSet
	CCVWeightVal          string = "weightVal"          // 可以修改配置项里的 weightVal
	CCVMapChgWeight       string = "mapChgWeight"       // 可以修改配置项里的 mapChgWeight，这里因为是个map，所以要当成 mapChgWeight:S 这样传递
	CCVTriggerWeight      string = "triggerWeight"      // 可以修改配置项里的 triggerWeight
	CCVRetriggerRespinNum string = "retriggerRespinNum" // 可以修改配置项里的 retriggerRespinNum
	CCVWinMulti           string = "winMulti"           // 可以修改配置项里的 winMulti
	CCVSavedMoney         string = "savedMoney"         // 可以修改配置项里的 savedMoney
)

const (
	CVSymbolNum string = "symbolNum" // 触发后，中奖的符号数量
	CVWildNum   string = "wildNum"   // 触发后，中奖符号里的wild数量
	CVRespinNum string = "respinNum" // 触发后，如果有产生respin的逻辑，这就是最终respin的次数
	CVWins      string = "wins"      // 中奖的数值，线注的倍数
)
