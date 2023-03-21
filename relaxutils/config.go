package relaxutils

import (
	"encoding/xml"
	"os"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type FloatList struct {
	Vals []float32 `xml:"item"`
}

type StringList struct {
	Vals []string `xml:"item"`
}

type IntList struct {
	Vals []int `xml:"item"`
}

type Symbol struct {
	XMLName xml.Name `xml:"symbol"`
	Name    string   `xml:"name,attr"`
	Val     int      `xml:"value,attr"`
}

type Symbols struct {
	XMLName xml.Name  `xml:"symbols"`
	Symbols []*Symbol `xml:"symbol"`
}

type PayoutWin struct {
	XMLName xml.Name `xml:"win"`
	Count   int      `xml:"count,attr"`
	Payout  int      `xml:"payout,attr"`
}

type Payout struct {
	XMLName xml.Name `xml:"symbol"`
	Name    string   `xml:"name,attr"`
	Wins    []*PayoutWin
}

type Payouts struct {
	XMLName xml.Name `xml:"payouts"`
	Payouts []*Payout
}

type Table struct {
	XMLName      xml.Name      `xml:"table"`
	TableComment string        `xml:",comment"`
	Reel         []*StringList `xml:"reel"`
}

type Reels struct {
	Tables []*Table
}

type Weights struct {
	Entries []*IntList `xml:"entry"`
}

type WeightsArr struct {
	Weights []*Weights `xml:"weights"`
}

type Int2DArray struct {
	Rows []*IntList `xml:"row"`
}

type Int3DArray struct {
	Tables []*Int2DArray `xml:"tbl"`
}

type Config struct {
	XMLName        xml.Name    `xml:"config"`
	GeneralComment string      `xml:",comment"`
	SD             float32     `xml:"sd"`
	RTP            float32     `xml:"rtp"`
	ID             int         `xml:"id"`
	Name           string      `xml:"name"`
	ConfigVersion  string      `xml:"configVersion"`
	RTPsComment    string      `xml:",comment"`
	RTPs           *FloatList  `xml:"rtps"`
	ConfigComment  string      `xml:",comment"`
	Symbols        *Symbols    `xml:"symbols"`
	Wilds          *StringList `xml:"wilds"`
	Payouts        *Payouts    `xml:"payouts"`
	PayingSymbols  *StringList `xml:"payingSymbols"`
}

func NewConfig() *Config {
	return &Config{
		GeneralComment: " General game information ",
		RTPsComment:    " RTP infomation ",
		ConfigComment:  " Configuration ",
	}
}

type FuncOnSelfCloseTags func(string) string

func SaveConfig(fn string, cfg interface{}, procSelfCloseTags FuncOnSelfCloseTags) error {
	output, err := xml.MarshalIndent(cfg, "  ", "    ")
	if err != nil {
		goutils.Error("SaveConfig:MarshalIndent",
			zap.Error(err))

		return err
	}

	xmlhead := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	buf := []byte(xmlhead)
	buf = append(buf, output...)

	if procSelfCloseTags != nil {
		buf = []byte(procSelfCloseTags(string(buf)))
	}

	err = os.WriteFile(fn, buf, 0644)
	if err != nil {
		goutils.Error("SaveConfig:WriteFile",
			zap.Error(err))

		return err
	}

	return nil
}
