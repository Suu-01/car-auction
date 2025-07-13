package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore" // ← 추가
)

var Logger *zap.Logger

func Init() {
	// 개발 단계에선 Development config, 프로덕션에선 Production config 사용
	cfg := zap.NewProductionConfig()
	// 예: ISO8601 타임스탬프, JSON 출력
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Logger = logger
}
