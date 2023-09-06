package lowcode

import (
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
)

func loadBasicInfo(cfg *Config, buf []byte) error {
	gameName, err := sonic.Get(buf, "gameName")
	if err != nil {
		goutils.Error("loadBasicInfo:Get",
			zap.String("key", "gameName"),
			zap.Error(err))

		return err
	}

	cfg.Name, _ = gameName.String()

	lstParam, err := sonic.Get(buf, "parameter")
	if err != nil {
		goutils.Error("loadBasicInfo:Get",
			zap.String("key", "parameter"),
			zap.Error(err))

		return err
	}

	lst, err := lstParam.ArrayUseNode()
	if err != nil {
		goutils.Error("loadBasicInfo:ArrayUseNode",
			zap.Error(err))

		return err
	}

	for i, v := range lst {
		str, err := v.Get("name").String()
		if err != nil {
			goutils.Error("loadBasicInfo:name",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}

		if str == "Width" {
			w, err := v.Get("value").Int64()
			if err != nil {
				goutils.Error("loadBasicInfo:value",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}

			cfg.Width = int(w)
		} else if str == "Height" {
			h, err := v.Get("value").Int64()
			if err != nil {
				goutils.Error("loadBasicInfo:value",
					zap.Int("i", i),
					zap.Error(err))

				return err
			}

			cfg.Height = int(h)
		}
	}

	return nil
}

func parse2IntSlice(n *ast.Node) ([]int, error) {
	arr, err := n.ArrayUseNode()
	if err != nil {
		goutils.Error("parse2IntSlice:Array",
			zap.Error(err))

		return nil, err
	}

	iarr := []int{}

	for i, v := range arr {
		iv, err := v.Int64()
		if err != nil {
			goutils.Error("parse2IntSlice:Int64",
				zap.Int("i", i),
				zap.Error(err))

			return nil, err
		}

		iarr = append(iarr, int(iv))
	}

	return iarr, nil
}

func parsePaytables(n *ast.Node) (*sgc7game.PayTables, error) {
	paytables := &sgc7game.PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	syms, err := n.ArrayUseNode()
	if err != nil {
		goutils.Error("parsePaytables:ArrayUseNode",
			zap.Error(err))

		return nil, err
	}

	for j, sym := range syms {
		c, err := sym.Get("Code").Int64()
		if err != nil {
			goutils.Error("parsePaytables:syms:Code",
				zap.Int("j", j),
				zap.Error(err))

			return nil, err
		}

		s, err := sym.Get("Symbol").String()
		if err != nil {
			goutils.Error("parsePaytables:syms:Symbol",
				zap.Int("j", j),
				zap.Error(err))

			return nil, err
		}

		arr, err := parse2IntSlice(sym.Get("data"))
		if err != nil {
			goutils.Error("parsePaytables:syms:data",
				zap.Int("j", j),
				zap.Error(err))

			return nil, err
		}

		paytables.MapSymbols[s] = int(c)
		paytables.MapPay[int(c)] = arr
	}

	return paytables, nil
}

func loadPaytables(cfg *Config, lstPaytables *ast.Node) error {
	lst, err := lstPaytables.ArrayUseNode()
	if err != nil {
		goutils.Error("loadPaytables:ArrayUseNode",
			zap.Error(err))

		return err
	}

	for i, v := range lst {
		k, err := v.Get("fileName").String()
		if err != nil {
			goutils.Error("loadPaytables:fileName",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}

		paytables, err := parsePaytables(v.Get("fileJson"))
		if err != nil {
			goutils.Error("loadPaytables:parsePaytables",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}

		cfg.Paytables[k] = k
		cfg.MapPaytables[k] = paytables
	}

	return nil
}

func NewGame2(fn string, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
	game := &Game{
		BasicGame:    sgc7game.NewBasicGame(funcNewPlugin),
		MgrComponent: NewComponentMgr(),
	}

	cfg := &Config{
		Paytables:    make(map[string]string),
		MapPaytables: make(map[string]*sgc7game.PayTables),
	}

	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("NewGame2:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	err = loadBasicInfo(cfg, data)
	if err != nil {
		goutils.Error("NewGame2:loadBasicInfo",
			zap.Error(err))

		return nil, err
	}

	lstPaytables, err := sonic.Get(data, "repository", "paytableList")
	if err != nil {
		goutils.Error("NewGame2:Get",
			zap.String("key", "repository.paytableList"),
			zap.Error(err))

		return nil, err
	}

	err = loadPaytables(cfg, &lstPaytables)
	if err != nil {
		goutils.Error("NewGame2:loadPaytables",
			zap.Error(err))

		return nil, err
	}

	return game, nil
}
