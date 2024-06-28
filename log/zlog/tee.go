package zlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

type LevelEnablerFunc func(level Level) bool

type TeeOption struct {
	Out io.Writer
	LevelEnablerFunc
}

func NewTee(tees []TeeOption, opts ...Option) *Logger {
	var cores []zapcore.Core
	for _, tee := range tees {
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.RFC3339TimeEncoder
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg),
			zapcore.AddSync(tee.Out),
			zap.LevelEnablerFunc(tee.LevelEnablerFunc),
		)
		cores = append(cores, core)
	}
	return &Logger{
		l:  zap.New(zapcore.NewTee(cores...), opts...),
		al: nil,
	}
}
