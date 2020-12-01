package sgc7game

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

type lineInfo5 struct {
	R1   int `json:"R1"`
	R2   int `json:"R2"`
	R3   int `json:"R3"`
	R4   int `json:"R4"`
	R5   int `json:"R5"`
	Line int `json:"line"`
}

// LineData - line data
type LineData struct {
	Lines [][]int `json:"lines"`
}

// isValidLI5 - is it valid lineInfo5
func isValidLI5(li5s []lineInfo5) bool {
	if len(li5s) <= 0 {
		return false
	}

	// alllinezero := true
	for _, v := range li5s {
		if v.Line > 0 {
			// alllinezero = false

			return true
		}
	}

	return false
}

// LoadLine5JSON - load json file
func LoadLine5JSON(fn string) (*LineData, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []lineInfo5
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if !isValidLI5(li) {
		return nil, nil
	}

	d := &LineData{}
	for _, v := range li {
		cl := []int{v.R1, v.R2, v.R3, v.R4, v.R5}
		d.Lines = append(d.Lines, cl)
	}

	return d, nil
}

// LoadLine3JSON - load json file
func LoadLine3JSON(fn string) (*LineData, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []lineInfo5
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if !isValidLI5(li) {
		return nil, nil
	}

	d := &LineData{}
	for _, v := range li {
		cl := []int{v.R1, v.R2, v.R3}
		d.Lines = append(d.Lines, cl)
	}

	return d, nil
}
