package sgc7pbutils

import (
	"log/slog"

	"github.com/zhs007/goutils"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func PB2Json(msg proto.Message) ([]byte, error) {
	result, err := protojson.Marshal(msg)
	if err != nil {
		goutils.Error("PB2Json:Marshal",
			goutils.Err(err))

		return nil, err
	}

	return result, nil
}

// JSON - It's like slog.String(name, str)
func JSON(name string, msg proto.Message) slog.Attr {
	result, err := protojson.Marshal(msg)
	if err != nil {
		goutils.Error("JSON:Marshal",
			goutils.Err(err))

		return slog.String(name, err.Error())
	}

	return slog.String(name, string(result))
}
