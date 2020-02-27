package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	L *zap.Logger
	S *zap.SugaredLogger
)

func Init() {
	writeSyncer := zapcore.AddSync(os.Stdout)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	L = zap.New(core, zap.AddCaller())
	S = L.Sugar()

	zap.RedirectStdLog(L)
	zap.ReplaceGlobals(L)
}
