package gatiserv

import (
	"crypto/sha1"
	"encoding/hex"
	"os"

	jsoniter "github.com/json-iterator/go"
	goutils "github.com/zhs007/goutils"
	"go.uber.org/zap"
)

// Checksum - it's like sha1sum fn
func Checksum(fn string) (string, error) {
	data, err := os.ReadFile(fn)
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
			goutils.Error("GenChecksum:Checksum",
				zap.String("filename", v.Filename),
				zap.Error(err))

			return nil, err
		}

		v.Checksum = hash

		ccs.Components[v.ID] = v
	}

	return ccs, nil
}

// // LoadGATICriticalComponents - load
// func LoadGATICriticalComponents(fn string) (*GATICriticalComponents, error) {
// 	if fn == "" {
// 		return &GATICriticalComponents{
// 			Components: make(map[int]*GATICriticalComponent),
// 		}, nil
// 	}

// 	json := jsoniter.ConfigCompatibleWithStandardLibrary

// 	data, err := os.ReadFile(fn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ccs := &GATICriticalComponents{}
// 	err = json.Unmarshal(data, ccs)
// 	if err != nil {
// 		goutils.Warn("gatiserv.LoadGATICriticalComponents",
// 			zap.Error(err))

// 		return nil, err
// 	}

// 	return ccs, nil
// }

// // SaveGATICriticalComponents - save
// func SaveGATICriticalComponents(ccs *GATICriticalComponents, fn string) error {
// 	json := jsoniter.ConfigCompatibleWithStandardLibrary

// 	b, err := json.Marshal(ccs)
// 	if err != nil {
// 		goutils.Warn("gatiserv.SaveGATICriticalComponents",
// 			zap.Error(err))

// 		return err
// 	}

// 	os.WriteFile(fn, b, 0640)

// 	return nil
// }

// LoadGATIGameInfo - load
func LoadGATIGameInfo(fn string) (*GATIGameInfo, error) {
	if fn == "" {
		return &GATIGameInfo{
			Components: make(map[int]*GATICriticalComponent),
		}, nil
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	ccs := &GATIGameInfo{}
	err = json.Unmarshal(data, ccs)
	if err != nil {
		goutils.Warn("gatiserv.LoadGATIGameInfo",
			zap.Error(err))

		return nil, err
	}

	return ccs, nil
}

// SaveGATIGameInfo - save
func SaveGATIGameInfo(gi *GATIGameInfo, fn string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(gi)
	if err != nil {
		goutils.Warn("gatiserv.SaveGATIGameInfo",
			zap.Error(err))

		return err
	}

	os.WriteFile(fn, b, 0640)

	return nil
}
