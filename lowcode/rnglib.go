package lowcode

import (
	"log/slog"
	"os"
	"sync"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"gopkg.in/yaml.v2"
)

type RngData struct {
	Script string         `yaml:"script"`
	FO2    *ForceOutcome2 `yaml:"-"`
}

type RngLib struct {
	mutex   sync.Mutex          `yaml:"-"`
	MapRNGs map[string]*RngData `yaml:"rngs"`
}

func (rngLib *RngLib) onResults(results []*sgc7game.PlayResult) string {
	rngLib.mutex.Lock()
	defer rngLib.mutex.Unlock()

	for k, v := range rngLib.MapRNGs {
		if v.FO2.IsValid(results) {
			delete(rngLib.MapRNGs, k)

			return k
		}
	}

	return ""
}

func LoadRngLib(fn string) (*RngLib, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("LoadRngLib:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	rngLib := &RngLib{}
	err = yaml.Unmarshal(data, rngLib)
	if err != nil {
		goutils.Error("LoadRngLib:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	for _, v := range rngLib.MapRNGs {
		fo2, err := NewForceOutcome2(v.Script)
		if err != nil {
			goutils.Error("LoadRngLib:NewForceOutcome2",
				goutils.Err(err))

			return nil, err
		}

		v.FO2 = fo2
	}

	return rngLib, nil
}
