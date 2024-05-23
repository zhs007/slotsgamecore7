package lowcode

const (
	CCVReelSet               string = "reelSet"               // 可以修改配置项里的 reelSet
	CCVWeightVal             string = "weightVal"             // 可以修改配置项里的 weightVal
	CCVMapChgWeight          string = "mapChgWeight"          // 可以修改配置项里的 mapChgWeight，这里因为是个map，所以要当成 mapChgWeight:S 这样传递
	CCVTriggerWeight         string = "triggerWeight"         // 可以修改配置项里的 triggerWeight
	CCVRetriggerRespinNum    string = "retriggerRespinNum"    // 可以修改配置项里的 retriggerRespinNum
	CCVWinMulti              string = "winMulti"              // 可以修改配置项里的 winMulti
	CCVSavedMoney            string = "savedMoney"            // 可以修改配置项里的 savedMoney
	CCVSymbolNum             string = "symbolNum"             // 可以修改配置项里的 symbolNum
	CCVInputVal              string = "inputVal"              // 可以修改配置项里的 inputVal
	CCVValueNum              string = "valueNum"              // 可以修改配置项里的 valueNum
	CCVQueue                 string = "queue"                 // 可以修改配置项里的 queue
	CCVLastRespinNum         string = "lastRespinNum"         // 可以修改配置项里的 lastRespinNum
	CCVRetriggerAddRespinNum string = "retriggerAddRespinNum" // 可以修改配置项里的 retriggerAddRespinNum
	CCVLastTriggerNum        string = "lastTriggerNum"        // 可以修改配置项里的 lastTriggerNum
)

const (
	CVSymbolNum             string = "symbolNum"             // 触发后，中奖的符号数量
	CVWildNum               string = "wildNum"               // 触发后，中奖符号里的wild数量
	CVRespinNum             string = "respinNum"             // 触发后，如果有产生respin的逻辑，这就是最终respin的次数
	CVWins                  string = "wins"                  // 中奖的数值，线注的倍数
	CVCurRespinNum          string = "curRespinNum"          // curRespinNum
	CVCurTriggerNum         string = "curTriggerNum"         // curTriggerNum
	CVLastRespinNum         string = "lastRespinNum"         // lastRespinNum
	CVLastTriggerNum        string = "lastTriggerNum"        // lastTriggerNum
	CVRetriggerAddRespinNum string = "retriggerAddRespinNum" // retriggerAddRespinNum
	CVAvgSymbolValMulti     string = "avgSymbolValMulti"     // avgSymbolValMulti
	CVAvgHeight             string = "avgHeight"             // avgHeight
	CVWinMulti              string = "winMulti"              // winMulti
	CVNumber                string = "number"                // number
)
