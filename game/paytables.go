package sgc7game

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

type payInfo struct {
	Code   int    `json:"Code"`
	Symbol string `json:"Symbol"`
	X1     int    `json:"X1"`
	X2     int    `json:"X2"`
	X3     int    `json:"X3"`
	X4     int    `json:"X4"`
	X5     int    `json:"X5"`
	X6     int    `json:"X6"`
	X7     int    `json:"X7"`
	X8     int    `json:"X8"`
	X9     int    `json:"X9"`
	X10    int    `json:"X10"`
	X11    int    `json:"X11"`
	X12    int    `json:"X12"`
	X13    int    `json:"X13"`
	X14    int    `json:"X14"`
	X15    int    `json:"X15"`
}

// PayTables - pay tables
type PayTables struct {
	MapPay     map[int][]int  `json:"paytables"`
	MapSymbols map[string]int `json:"symbols"`
}

// LoadPayTables5JSON - load json file
func LoadPayTables5JSON(fn string) (*PayTables, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []payInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if len(li) <= 0 {
		return nil, nil
	}

	p := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3, v.X4, v.X5}
		p.MapPay[v.Code] = cl

		p.MapSymbols[v.Symbol] = v.Code
	}

	return p, nil
}

// LoadPayTables3JSON - load json file
func LoadPayTables3JSON(fn string) (*PayTables, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []payInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if len(li) <= 0 {
		return nil, nil
	}

	p := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3}
		p.MapPay[v.Code] = cl

		p.MapSymbols[v.Symbol] = v.Code
	}

	return p, nil
}

// LoadPayTables6JSON - load json file
func LoadPayTables6JSON(fn string) (*PayTables, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []payInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if len(li) <= 0 {
		return nil, nil
	}

	p := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3, v.X4, v.X5, v.X6}
		p.MapPay[v.Code] = cl

		p.MapSymbols[v.Symbol] = v.Code
	}

	return p, nil
}

// LoadPayTables15JSON - load json file
func LoadPayTables15JSON(fn string) (*PayTables, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []payInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if len(li) <= 0 {
		return nil, nil
	}

	p := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3, v.X4, v.X5, v.X6, v.X7, v.X8, v.X9, v.X10, v.X11, v.X12, v.X13, v.X14, v.X15}
		p.MapPay[v.Code] = cl

		p.MapSymbols[v.Symbol] = v.Code
	}

	return p, nil
}
