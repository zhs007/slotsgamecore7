package mathtoolset2

import (
	"fmt"

	"github.com/valyala/fastrand"
)

type SymbolData struct {
	Symbol string
	Num    int
	weight int
}

func (sd *SymbolData) String() string {
	if sd.Num == 1 {
		return sd.Symbol
	}

	return fmt.Sprintf("%v_%v", sd.Symbol, sd.Num)
}

type SymbolsPool struct {
	Pool []*SymbolData
}

func (pool *SymbolsPool) Clone() *SymbolsPool {
	newPool := &SymbolsPool{}

	for _, sd := range pool.Pool {
		newPool.Pool = append(newPool.Pool, &SymbolData{
			Symbol: sd.Symbol,
			Num:    sd.Num,
		})
	}

	return newPool
}

func (pool *SymbolsPool) Push(symbol string, num int) {
	pool.Pool = append(pool.Pool, &SymbolData{Symbol: symbol, Num: num})
}

func (pool *SymbolsPool) PushEx(symbol string, symbolNum int, num int) {
	for range num {
		pool.Pool = append(pool.Pool, &SymbolData{Symbol: symbol, Num: symbolNum})
	}
}

func (pool *SymbolsPool) Remove(symbol string, symbolNum int) {
	for i := 0; i < len(pool.Pool); i++ {
		if pool.Pool[i].Symbol == symbol && pool.Pool[i].Num == symbolNum {
			pool.Pool = append(pool.Pool[:i], pool.Pool[i+1:]...)
			return
		}
	}
}

func (pool *SymbolsPool) CountAllSymbolNumber() int {
	count := 0

	for _, sd := range pool.Pool {
		count += sd.Num
	}

	return count
}

func (pool *SymbolsPool) CountSymbolDataNumber(symbol string) int {
	count := 0

	for _, sd := range pool.Pool {
		if sd.Symbol == symbol {
			count++
		}
	}

	return count
}

func (pool *SymbolsPool) getList() []*SymbolData {
	lst := []*SymbolData{}
	mapSymbols := make(map[string]map[int]*SymbolData)

	for _, sd := range pool.Pool {
		misd, exists := mapSymbols[sd.Symbol]
		if !exists {
			mapSymbols[sd.Symbol] = make(map[int]*SymbolData)
			mapSymbols[sd.Symbol][sd.Num] = sd
		} else {
			if _, exists := misd[sd.Num]; !exists {
				misd[sd.Num] = sd
			}
		}
	}

	for _, misd := range mapSymbols {
		if len(misd) == 1 {
			for _, sd := range misd {
				lst = append(lst, sd)

				sd.weight = pool.CountSymbolDataNumber(sd.Symbol)
			}
		} else {
			maxv := -1
			for n := range misd {
				if n > maxv {
					maxv = n
				}
			}

			lst = append(lst, misd[maxv])

			misd[maxv].weight = pool.CountSymbolDataNumber(misd[maxv].Symbol)
		}
	}

	return lst
}

func randomSymbolData(lst []*SymbolData) *SymbolData {
	totalWeight := 0
	for _, sd := range lst {
		totalWeight += sd.weight
	}

	v := fastrand.Uint32n(uint32(totalWeight))

	for _, sd := range lst {
		if v < uint32(sd.weight) {
			return sd
		}
		v -= uint32(sd.weight)
	}

	return lst[len(lst)-1] // Fallback to the last symbol if something goes wrong
}

func rmSymbolData(lst []*SymbolData, sd *SymbolData) []*SymbolData {
	for i, v := range lst {
		if v.Symbol == sd.Symbol && v.Num == sd.Num {
			return append(lst[:i], lst[i+1:]...)
		}
	}

	return lst // Return the original list if not found
}
