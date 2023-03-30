package sgc7pbutils

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func PB2Json(msg proto.Message) ([]byte, error) {
	result, err := protojson.Marshal(msg)
	if err != nil {
		goutils.Error("PB2Json:Marshal",
			zap.Error(err))

		return nil, err
	}

	return result, nil
}

// JSON - It's like zap.String(name, str)
func JSON(name string, msg proto.Message) zap.Field {
	result, err := protojson.Marshal(msg)
	if err != nil {
		goutils.Error("JSON:Marshal",
			zap.Error(err))

		return zap.String(name, err.Error())
	}

	return zap.String(name, string(result))
}
