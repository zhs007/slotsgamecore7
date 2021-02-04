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
}

// PayTables - pay tables
type PayTables struct {
	MapPay map[int][]int `json:"paytables"`
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
		MapPay: make(map[int][]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3, v.X4, v.X5}
		p.MapPay[v.Code] = cl
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
		MapPay: make(map[int][]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3}
		p.MapPay[v.Code] = cl
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
		MapPay: make(map[int][]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3, v.X4, v.X5, v.X6}
		p.MapPay[v.Code] = cl
	}

	return p, nil
}
