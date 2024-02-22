package mathtoolset2

import (
	"bytes"
	"encoding/base64"
	"io"

	"github.com/bytedance/sonic"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type FileDataMap struct {
	MapFiles map[string]string `json:"files"`
}

func (mapfd *FileDataMap) GetReader(fn string) io.Reader {
	str, isok := mapfd.MapFiles[fn]
	if !isok {
		return nil
	}

	sDec, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		goutils.Error("FileDataMap.GetReader:DecodeString",
			zap.Error(err))

		return nil
	}

	return bytes.NewBuffer(sDec)
}

func (mapfd *FileDataMap) AddReader(fn string, r io.Reader) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	sEnc := base64.StdEncoding.EncodeToString(buf.Bytes())

	mapfd.MapFiles[fn] = sEnc

	return nil
}

func (mapfd *FileDataMap) AddBuffer(fn string, buf *bytes.Buffer) {
	sEnc := base64.StdEncoding.EncodeToString(buf.Bytes())

	mapfd.MapFiles[fn] = sEnc
}

func (mapfd *FileDataMap) ToJson() (string, error) {
	buf, err := sonic.Marshal(mapfd)
	if err != nil {
		goutils.Error("FileDataMap:ToJson:Marshal",
			zap.Error(err))

		return "", err
	}

	return string(buf), nil
}

func NewFileDataMap(fd string) (*FileDataMap, error) {
	if len(fd) == 0 {
		return &FileDataMap{
			MapFiles: make(map[string]string),
		}, nil
	}

	mapFD := &FileDataMap{}
	err := sonic.Unmarshal([]byte(fd), mapFD)
	if err != nil {
		goutils.Error("NewFileDataMap:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	return mapFD, nil
}
