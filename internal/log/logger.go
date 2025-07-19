package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// Init はロガーを初期化します。
// 開発環境では DevelopmentConfig、プロダクション環境では ProductionConfig を使用します。
func Init() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Logger = logger
}
