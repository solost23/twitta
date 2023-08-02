package initialize

import (
	"errors"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(filePath string) {
	var logger *zap.Logger
	switch gin.Mode() {
	case gin.DebugMode:
		logger, _ = zap.NewDevelopment()
	case gin.ReleaseMode:
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
