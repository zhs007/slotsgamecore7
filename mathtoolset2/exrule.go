package mathtoolset2

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/zhs007/goutils"
)

type ExRule struct {
	Code    string
	Params  []int
	Symbols []string
}

func (rule *ExRule) procWeight(sd []*SymbolData) {
	if rule.Code == "EXC" {
		totalWeight := 0
		for _, s := range sd {
			if slices.Contains(rule.Symbols, s.Symbol) {
				totalWeight += s.weight
			}
		}

		totalWeight *= len(rule.Symbols)

		for _, s := range sd {
			if slices.Contains(rule.Symbols, s.Symbol) {
				s.weight = totalWeight
			}
		}
	}
}

func (rule *ExRule) isOKForHead(rd []string, sd *SymbolData, lastNum int) bool {
	// separation
	if rule.Code == "SEP" {
		end := lastNum - sd.Num
		if end <= rule.Params[0] {
			for i := range rule.Params[0] - end {
				if rd[i] == sd.Symbol {
					return false
				}
			}
		}
	} else if rule.Code == "EXC" {
		end := lastNum - sd.Num
		if end <= rule.Params[0] {
			for i := range rule.Params[0] - end {
				if slices.Contains(rule.Symbols, rd[i]) {
					return false
				}
			}
		}
	}

	return true
}

func (rule *ExRule) IsOK(rd []string, sd *SymbolData, lastNum int) bool {
	// separation
	if rule.Code == "SEP" {
		if !slices.Contains(rule.Symbols, sd.Symbol) {
			return true
		}

		if len(rd) <= rule.Params[0] {
			if slices.Contains(rd, sd.Symbol) {
				return false
			}
		} else {
			for i := range rule.Params[0] {
				if rd[len(rd)-1-i] == sd.Symbol {
					return false
				}
			}
		}
	} else if rule.Code == "EXC" {
		if !slices.Contains(rule.Symbols, sd.Symbol) {
			return true
		}

		if len(rd) <= rule.Params[0] {
			for _, s := range rule.Symbols {
				if slices.Contains(rd, s) {
					return false
				}
			}
		} else {
			for i := range rule.Params[0] {
				if slices.Contains(rule.Symbols, rd[len(rd)-1-i]) {
					return false
				}
			}
		}
	}

	return rule.isOKForHead(rd, sd, lastNum)
}

// ParseExRule - code is like "SEP_3,SC,WL"
func ParseExRule(code string) (*ExRule, error) {
	if code == "" {
		return nil, ErrInvalidCode
	}

	arr := strings.Split(code, ",")
	if len(arr) < 2 {
		return nil, ErrInvalidCode
	}

	rule := &ExRule{}

	codes := strings.Split(arr[0], "_")
	if len(codes) > 1 {
		for i, v := range codes {
			if i == 0 {
				rule.Code = v
			} else {
				param, err := goutils.String2Int64(v)
				if err != nil {
					goutils.Error("ParseExRule:Rule:Params",
						slog.Any("arr", arr),
						goutils.Err(err))

					return nil, err
				}

				rule.Params = append(rule.Params, int(param))
			}
		}
	} else {
		rule.Code = arr[0]
	}

	rule.Symbols = arr[1:]

	return rule, nil
}

// ParseExRules - code is like "OFF_3,SC,WL;OFF_5,H1,H2;"
func ParseExRules(code string) ([]*ExRule, error) {
	rules := []*ExRule{}
	arr := strings.Split(code, ";")
	for _, v := range arr {
		if v != "" {
			rule, err := ParseExRule(v)
			if err != nil {
				goutils.Error("ParseExRules:ParseExRule",
					slog.Any("code", code),
					goutils.Err(err))

				return nil, err
			}

			rules = append(rules, rule)
		}
	}

	return rules, nil
}

func BuildCurSymbols(rd []string, rules []*ExRule, pool *SymbolsPool) []*SymbolData {
	lst := []*SymbolData{}
	lastNum := pool.CountAllSymbolNumber()
	lstSymbols := pool.getList()

	for _, sd := range lstSymbols {
		isok := true
		for _, rule := range rules {
			if !rule.IsOK(rd, sd, lastNum) {
				isok = false
				break
			}
		}

		if isok {
			lst = append(lst, sd)
		}
	}

	for _, rule := range rules {
		rule.procWeight(lst)
	}

	return lst
}
