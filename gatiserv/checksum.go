package gatiserv

import (
	"crypto/sha1"
	"encoding/hex"
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	goutils "github.com/zhs007/goutils"
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
				slog.String("filename", v.Filename),
				goutils.Err(err))

			return nil, err
		}

		v.Checksum = hash

		ccs.Components[v.ID] = v
	}

	return ccs, nil
}

// LoadGATIGameInfo - load
func LoadGATIGameInfo(fn string) (*GATIGameInfo, error) {
	if fn == "" {
		return &GATIGameInfo{
			Components: make(map[int]*GATICriticalComponent),
		}, nil
	}

	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	ccs := &GATIGameInfo{}
	err = sonic.Unmarshal(data, ccs)
	if err != nil {
		goutils.Warn("gatiserv.LoadGATIGameInfo",
			goutils.Err(err))

		return nil, err
	}

	return ccs, nil
}

// SaveGATIGameInfo - save
func SaveGATIGameInfo(gi *GATIGameInfo, fn string) error {
	b, err := sonic.Marshal(gi)
	if err != nil {
		goutils.Warn("gatiserv.SaveGATIGameInfo",
			goutils.Err(err))

		return err
	}

	os.WriteFile(fn, b, 0640)

	return nil
}
