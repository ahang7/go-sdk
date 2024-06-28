package main

import (
	"fmt"
	"github.com/ahang7/go-sdk/log/zlog"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

func main() {
	// 开箱即用
	{
		defer zlog.Sync()
		zlog.Info("hello zlog!", zlog.String("zzz", "xxx"))
		zlog.Warn("warn message", zlog.Int("int", 3))
		zlog.Error("Error message", zlog.Duration("backoff", time.Second))

		// 修改日志级别
		zlog.SetLevel(zlog.ErrorLevel)
		zlog.Info("Info message")
		zlog.Warn("warn message")
		zlog.Error("Error message")

		// 替换默认logger
		file, err := os.OpenFile("custom.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		logger := zlog.New(file, zlog.InfoLevel)
		zlog.ReplaceDefault(logger)
		zlog.Info("Info message in replace default logger")
	}

	// 选项
	{
		file, err := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		opts := []zlog.Option{
			zlog.WithCaller(true),
			zlog.AddCallerSkip(1),
			zlog.Hooks(func(entry zapcore.Entry) error {
				if entry.Level == zlog.WarnLevel {
					fmt.Printf("Warn Hook: msg = %s\n", entry.Message)
				}
				return nil
			}),
			zlog.WithFatalHook(Hook{}),
		}
		logger := zlog.New(io.MultiWriter(os.Stdout, file), zlog.InfoLevel, opts...)
		defer logger.Sync()

		logger.Info("Info message", zlog.String("val", "hello world"))
		logger.Warn("Warn message", zlog.Duration("backoff", time.Second))
		logger.Fatal("Fatal msg", zlog.Time("val", time.Now()))
	}

	// 不同级别日志输出到不同位置
	{
		file, err := os.OpenFile("test-warn.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		tees := []zlog.TeeOption{
			{
				Out: os.Stdout,
				LevelEnablerFunc: zlog.LevelEnablerFunc(func(level zlog.Level) bool {
					return level == zlog.InfoLevel
				}),
			},
			{
				Out: file,
				LevelEnablerFunc: zlog.LevelEnablerFunc(func(level zlog.Level) bool {
					return level == zlog.WarnLevel
				}),
			},
		}
		logger := zlog.NewTee(tees)
		defer logger.Sync()

		logger.Info("Info tee message")
		logger.Warn("Warn tee message")
		logger.Error("Error tee message")
	}

	// 日志轮转
	{
		tees := []zlog.TeeOption{
			{
				Out: zlog.NewProductionRotateBySize("rotate-by-size.log"),
				LevelEnablerFunc: zlog.LevelEnablerFunc(
					func(level zlog.Level) bool {
						return level < zlog.WarnLevel
					}),
			},
			{
				Out: zlog.NewProductionRotateByTime("rotate-by-time.log"),
				LevelEnablerFunc: zlog.LevelEnablerFunc(func(level zlog.Level) bool {
					return level >= zlog.WarnLevel
				}),
			},
		}
		lts := zlog.NewTee(tees)
		defer lts.Sync()

		lts.Debug("Debug TeeAndRotate")
		lts.Info("Info TeeAndRotate")
		lts.Warn("Warn TeeAndRotate")
		lts.Error("Error TeeAndRotate")
	}
}

type Hook struct {
}

func (h Hook) OnWrite(ce *zapcore.CheckedEntry, fields []zapcore.Field) {
	fmt.Printf("Fatal Hook: msg=%s, field=%+v\n", ce.Message, fields)
}
