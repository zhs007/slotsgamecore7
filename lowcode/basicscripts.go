package lowcode

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

// setval(type, val)
func runSetVal(gameProp *GameProperty, scriptFunc *ScriptFunc) error {
	err := gameProp.SetVal(scriptFunc.IntParams[0], scriptFunc.IntParams[1])
	if err != nil {
		goutils.Error("runSetVal:SetVal",
			goutils.JSON("scriptFunc", scriptFunc),
			zap.Error(err))

		return err
	}

	return nil
}

func initSetVal(scriptFunc *ScriptFunc, script string) error {
	// err := gameProp.SetVal(scriptFunc.IntParams[0], scriptFunc.IntParams[1])
	// if err != nil {
	// 	goutils.Error("runSetVal:SetVal",
	// 		goutils.JSON("scriptFunc", scriptFunc),
	// 		zap.Error(err))

	// 	return err
	// }

	return nil
}
