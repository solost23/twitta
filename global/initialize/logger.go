package initialize

import (
	"errors"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"twitta/global"
)

func InitLogger(filePath string) {
	var logger *zap.Logger
	switch global.ServerConfig.DebugMode {
	case "debug":
		logger, _ = zap.NewDevelopment()
	case "release":
		logger, _ = zap.NewProduction(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(getFile(filePath)), zapcore.DebugLevel)
		}))
	default:
		panic(errors.New("失败"))
	}
	zap.ReplaceGlobals(logger)
	defer func() { _ = logger.Sync() }()
}

func getFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	filePath = path.Join(filePath, time.Now().Format(time.DateOnly))
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		_, err = os.Create(filePath)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return file
}
