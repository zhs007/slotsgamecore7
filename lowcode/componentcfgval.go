package lowcode

const (
	CCVReelSet                    string = "reelset"                    // 可以修改配置项里的 reelSet
	CCVWeightVal                  string = "weightval"                  // 可以修改配置项里的 weightVal
	CCVMapChgWeight               string = "mapchgweight"               // 可以修改配置项里的 mapChgWeight，这里因为是个map，所以要当成 mapChgWeight:S 这样传递
	CCVTriggerWeight              string = "triggerweight"              // 可以修改配置项里的 triggerWeight
	CCVRetriggerRespinNum         string = "retriggerrespinnum"         // 可以修改配置项里的 retriggerRespinNum
	CCVWinMulti                   string = "winmulti"                   // 可以修改配置项里的 winMulti
	CCVSavedMoney                 string = "savedmoney"                 // 可以修改配置项里的 savedMoney
	CCVSymbolNum                  string = "symbolnum"                  // 可以修改配置项里的 symbolNum
	CCVInputVal                   string = "inputval"                   // 可以修改配置项里的 inputVal
	CCVValueNum                   string = "valuenum"                   // 可以修改配置项里的 valueNum
	CCVQueue                      string = "queue"                      // 可以修改配置项里的 queue
	CCVLastRespinNum              string = "lastrespinnum"              // 可以修改配置项里的 lastRespinNum
	CCVRetriggerAddRespinNum      string = "retriggeraddrespinnum"      // 可以修改配置项里的 retriggerAddRespinNum
	CCVLastTriggerNum             string = "lasttriggernum"             // 可以修改配置项里的 lastTriggerNum
	CCVWeight                     string = "weight"                     // 可以修改配置项里的 weight
	CCVForceBranch                string = "forcebranch"                // 可以修改配置项里的 forceBranch
	CCVWins                       string = "wins"                       // 可以修改配置项里的 wins
	CCVMulti                      string = "multi"                      // 可以修改配置项里的 multi
	CCVNumber                     string = "number"                     // 可以修改配置项里的 number
	CCVHeight                     string = "height"                     // 可以修改配置项里的 height
	CCVReelSetWeight              string = "reelsetweight"              // 可以修改配置项里的 reelSetWeight
	CCVForceVal                   string = "forceval"                   // 可以修改配置项里的 forceVal
	CCVForceValNow                string = "forcevalnow"                // 可以修改配置项里的 forceValNow
	CCVClearForceTriggerOnceCache string = "clearforcetriggeroncecache" // 可以修改配置项里的 clearForceTriggerOnceCache
	CCVClear                      string = "clear"                      // 可以修改配置项里的 clear
	CCVLineData                   string = "linedata"                   // 可以修改配置项里的 linedata
	CCVClearNow                   string = "clearnow"                   // 可以修改配置项里的 clearNow
	CCVPickNum                    string = "pickNum"                    // 可以修改配置项里的 pickNum
)

const (
	CVSymbolNum             string = "symbolnum"             // 触发后，符号数量
	CVWildNum               string = "wildnum"               // 触发后，中奖符号里的wild数量
	CVRespinNum             string = "respinnum"             // 触发后，如果有产生respin的逻辑，这就是最终respin的次数
	CVWins                  string = "wins"                  // 中奖的数值，线注的倍数
	CVCurRespinNum          string = "currespinnum"          // curRespinNum
	CVCurTriggerNum         string = "curtriggernum"         // curTriggerNum
	CVLastRespinNum         string = "lastrespinnum"         // lastRespinNum
	CVLastTriggerNum        string = "lasttriggernum"        // lastTriggerNum
	CVRetriggerAddRespinNum string = "retriggeraddrespinnum" // retriggerAddRespinNum
	CVAvgSymbolValMulti     string = "avgsymbolvalmulti"     // avgSymbolValMulti
	CVAvgHeight             string = "avgheight"             // avgHeight
	CVWinMulti              string = "winmulti"              // winMulti
	CVNumber                string = "number"                // number
	CVResultNum             string = "resultnum"             // 触发后，中奖的数量
	CVOutputInt             string = "outputint"             // outputInt
	CVSelectedIndex         string = "selectedindex"         // selectedIndex，一般用于各种权重组件里，最后被选中的index
	CVValue                 string = "value"                 // value
	CVSymbolVal             string = "symbolval"             // symbolVal
	CVWinResultNum          string = "winresultnum"          // winResultNum
	CCValueNum              string = "valuenum"              // valueNum
	CVHeight                string = "height"                // height
)

const (
	CSVValue string = "value" // 组件触发后的具体值
)
