package relaxutils

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type Config struct {
	XMLName xml.Name `xml:"config"`
	General string   `xml:",comment"`
	SD      float32  `xml:"sd"`
}

func SaveConfig(fn string, cfg *Config) error {
	output, err := xml.MarshalIndent(cfg, "  ", "    ")
	if err != nil {
		goutils.Error("SaveConfig:MarshalIndent",
			zap.Error(err))

		return err
	}

	xmlhead := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	buf := []byte(xmlhead)
	buf = append(buf, output...)

	err = ioutil.WriteFile(fn, buf, 0644)
	if err != nil {
		goutils.Error("SaveConfig:WriteFile",
			zap.Error(err))

		return err
	}

	return nil
}
