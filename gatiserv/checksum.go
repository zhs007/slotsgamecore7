package gatiserv

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// Checksum - it's like sha1sum fn
func Checksum(fn string) (string, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	h.Write(data)

	return hex.EncodeToString(h.Sum(nil)), nil
}

// GenChecksum -
func GenChecksum(lst []*GATICriticalComponent) (*GATICriticalComponents, error) {
	ccs := &GATICriticalComponents{
		Components: make(map[int]*GATICriticalComponent),
	}

	for _, v := range lst {
		hash, err := Checksum(v.Filename)
		if err != nil {
			sgc7utils.Error("GenChecksum:Checksum",
				zap.String("filename", v.Filename),
				zap.Error(err))

			return nil, err
		}

		v.Checksum = hash

		ccs.Components[v.ID] = v
	}

	return ccs, nil
}

// LoadGATICriticalComponents - load
func LoadGATICriticalComponents(fn string) (*GATICriticalComponents, error) {
	if fn == "" {
		return &GATICriticalComponents{
			Components: make(map[int]*GATICriticalComponent),
		}, nil
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	ccs := &GATICriticalComponents{}
	err = json.Unmarshal(data, ccs)
	if err != nil {
		sgc7utils.Warn("gatiserv.LoadGATICriticalComponents",
			zap.Error(err))

		return nil, err
	}

	return ccs, nil
}

// SaveGATICriticalComponents - save
func SaveGATICriticalComponents(ccs *GATICriticalComponents, fn string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(ccs)
	if err != nil {
		sgc7utils.Warn("gatiserv.SaveGATICriticalComponents",
			zap.Error(err))

		return err
	}

	ioutil.WriteFile(fn, b, 0640)

	return nil
}
