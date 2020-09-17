package sgc7utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var onceLogger sync.Once
var logPath string

var panicFile *os.File
var logSubName string

func initPanicFile() error {
	file, err := os.OpenFile(
		path.Join(logPath, buildLogFilename("panic", logSubName)),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Warn("initPanicFile:OpenFile",
			zap.Error(err))

		return err
	}

	panicFile = file

	if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		Warn("initPanicFile:Dup2",
			zap.Error(err))

		return err
	}

	return nil
}

// string2LogLevel - string => zapcore.Level
func string2LogLevel(str string) zapcore.Level {
	if str == "debug" {
		return zapcore.DebugLevel
	}

	if str == "info" {
		return zapcore.InfoLevel
	}

	if str == "warn" {
		return zapcore.WarnLevel
	}

	if str == "error" {
		return zapcore.ErrorLevel
	}

	if str == "dpanic" {
		return zapcore.DPanicLevel
	}

	if str == "panic" {
		return zapcore.PanicLevel
	}

	if str == "fatal" {
		return zapcore.FatalLevel
	}

	return zapcore.WarnLevel
}

// buildLogSubFilename -
func buildLogSubFilename(appName string, version string) string {
	return fmt.Sprintf("%v.%v.%v", appName, version, FormatNow(gTime))
}

// buildLogFilename -
func buildLogFilename(logtype string, subname string) string {
	return fmt.Sprintf("%v.%v.log", subname, logtype)
}

func initLogger(appName string, appVersion string, strLevel string, isConsole bool, logpath string) (*zap.Logger, error) {
	logSubName = buildLogSubFilename(appName, appVersion)
	logPath = logpath

	level := string2LogLevel(strLevel)

	loglevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})

	if isConsole {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		consoleDebugging := zapcore.Lock(os.Stdout)
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleDebugging, loglevel),
		)

		cl := zap.New(core)
		// defer cl.Sync()

		return cl, nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfg := &zap.Config{}

	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.OutputPaths = []string{"stdout", path.Join(pwd, logpath, buildLogFilename("output", logSubName))}
	cfg.ErrorOutputPaths = []string{"stderr", path.Join(pwd, logpath, buildLogFilename("error", logSubName))}
	cfg.Encoding = "json"
	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:     "T",
		LevelKey:    "L",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		MessageKey:  "msg",
	}

	cl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	err = initPanicFile()
	if err != nil {
		return nil, err
	}

	return cl, nil
}

// InitLogger - initializes a thread-safe singleton logger
func InitLogger(logname string, appVersion string, level string, isConsole bool, logpath string) {

	// once ensures the singleton is initialized only once
	onceLogger.Do(func() {
		cl, err := initLogger(logname, appVersion, level, isConsole, logpath)
		if err != nil {
			fmt.Printf("initLogger error! %v \n", err)

			os.Exit(-1)
		}

		logger = cl
	})

	return
}

// // Log a message at the given level with given fields
// func Log(level zap.Level, message string, fields ...zap.Field) {
// 	singleton.Log(level, message, fields...)
// }

// Debug logs a debug message with the given fields
func Debug(message string, fields ...zap.Field) {
	if logger == nil {
		return
	}

	logger.Debug(message, fields...)
}

// Info logs a debug message with the given fields
func Info(message string, fields ...zap.Field) {
	if logger == nil {
		return
	}

	logger.Info(message, fields...)
}

// Warn logs a debug message with the given fields
func Warn(message string, fields ...zap.Field) {
	if logger == nil {
		return
	}

	logger.Warn(message, fields...)
}

// Error logs a debug message with the given fields
func Error(message string, fields ...zap.Field) {
	if logger == nil {
		return
	}

	logger.Error(message, fields...)
}

// Fatal logs a message than calls os.Exit(1)
func Fatal(message string, fields ...zap.Field) {
	if logger == nil {
		return
	}

	logger.Fatal(message, fields...)
}

// SyncLogger - sync logger
func SyncLogger() {
	logger.Sync()
}

// ClearLogs - clear logs
func ClearLogs() error {
	if logPath != "" {
		fn := path.Join(logPath, "*.log")
		lst, err := filepath.Glob(fn)
		if err != nil {
			return err
		}

		panicfile := buildLogFilename("panic", logSubName)
		outputfile := buildLogFilename("output", logSubName)
		errorfile := buildLogFilename("error", logSubName)

		for _, v := range lst {
			cfn := filepath.Base(v)
			if cfn != panicfile && cfn != outputfile && cfn != errorfile {
				os.Remove(v)
			}
		}
	}

	return nil
}

// GetLogger - get zap.Logger
func GetLogger() *zap.Logger {
	return logger
}

// JSON - It's like zap.String(name, str)
func JSON(name string, jobj interface{}) zap.Field {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(jobj)
	if err != nil {
		Warn("sgc7utils.JSON",
			zap.Error(err))

		return zap.String(name, err.Error())
	}

	return zap.String(name, string(b))
}
