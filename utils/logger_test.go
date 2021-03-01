package sgc7utils

import (
	"testing"
	"time"

	"go.uber.org/zap/zapcore"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_string2LogLevel(t *testing.T) {

	in := []string{
		"debug",
		"info",
		"warn",
		"dpanic",
		"panic",
		"error",
		"fatal",
		"",
		"haha",
	}

	out := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.DPanicLevel,
		zapcore.PanicLevel,
		zapcore.ErrorLevel,
		zapcore.FatalLevel,
		zapcore.WarnLevel,
		zapcore.WarnLevel,
	}

	for i, v := range in {
		ret := string2LogLevel(v)
		if ret != out[i] {
			t.Fatalf("Test_string2LogLevel string2LogLevel \"%s\" error",
				in[i])
		}
	}

	t.Logf("Test_string2LogLevel OK")
}

func Test_buildLogFilename(t *testing.T) {

	type blf struct {
		logtype string
		subname string
	}

	in := []blf{
		{
			logtype: "debug",
			subname: "main",
		},
		{
			logtype: "error",
			subname: "gamename",
		},
	}

	out := []string{
		"main.debug.log",
		"gamename.error.log",
	}

	for i, v := range in {
		ret := buildLogFilename(v.logtype, v.subname)
		if ret != out[i] {
			t.Fatalf("Test_buildLogFilename buildLogFilename [\"%s\", \"%s\"] != \"%s\" (\"%s\") ",
				v.logtype, v.subname, ret, out[i])
		}
	}

	t.Logf("Test_buildLogFilename OK")
}

func Test_MockLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockITime(ctrl)
	m.EXPECT().Now().Return(time.Unix(1597647832, 0))
	m.EXPECT().Now().Return(time.Unix(1597648787, 0))
	m.EXPECT().Now().Return(time.Unix(1597647832, 0))
	m.EXPECT().Now().Return(time.Unix(1597647832, 0))

	gTime = m

	type blsf struct {
		appName string
		version string
	}

	in := []blsf{
		{
			appName: "game",
			version: "v1.0.100",
		},
		{
			appName: "game2",
			version: "v2.100.256",
		},
	}

	out := []string{
		"game.v1.0.100.2020-08-17_07:03:52",
		"game2.v2.100.256.2020-08-17_07:19:47",
	}

	for i, v := range in {
		ret := buildLogSubFilename(v.appName, v.version)
		if ret != out[i] {
			t.Fatalf("Test_MockLogger buildLogSubFilename [\"%s\", \"%s\"] != \"%s\" (\"%s\") ",
				v.appName, v.version, ret, out[i])
		}
	}

	log, err := initLogger("main", "v1.0.0", "debug", false, "./")
	assert.NoError(t, err)
	assert.NotNil(t, log)

	InitLogger("main", "v1.0.0", "debug", true, "./")

	// 这里配合mock确认了 InitLogger 不会调用2次 initLogger
	InitLogger("main", "v1.0.0", "debug", true, "./")

	Debug("debug", zap.String("value", "123"))

	Info("info", zap.Int("value", 123))

	Warn("warn", zap.Int("value", 123))

	Error("info", zap.Int("value", 123))

	SyncLogger()

	ClearLogs()

	GetLogger()

	logger = nil

	Debug("debug", zap.String("value", "123"))

	Info("info", zap.Int("value", 123))

	Warn("warn", zap.Int("value", 123))

	Error("info", zap.Int("value", 123))

	Error("info", JSON("value", []int{1, 2, 3}))

	t.Logf("Test_MockLogger OK")
}
