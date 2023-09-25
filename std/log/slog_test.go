package log

import (
	"context"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"
)

func TestSlog(t *testing.T) {
	// slog 是一个结构化日志库，与传统的格式化日志库有区别，格式化日志库没有key的概念。
	// 日志格式以 key=val 方式输出，key 这个字段一定是字符串。
	val1, val2 := "val1", "val2"
	slog.Info("info log", "key1", val1, "key2", val2)
	// 2023/09/23 22:18:02 INFO msg  key=val
	stu := struct {
		a string
		b int
	}{
		"abc", 100,
	}
	// 结构体字段也有打印。
	slog.Info("slog any log", slog.Any("struct", stu))
	slog.Info("slog log", "struct", stu)

	// slog 没有 Fatal 日志级别，最高级别为 ERROR。
	// 默认 logger 也是输出到 os.Stdout。
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// 设置日志来源。
		AddSource: true,
		// 设置日志级别。
		Level: slog.LevelDebug,
		// ReplaceAttr 可以去掉默认的一些前缀。
	}))
	slog.SetDefault(logger)

	// 日志格式 time=xx level=xxx msg=xxx key=val
	slog.Debug("debug msg", "for", "bar", "now", time.Now())
	slog.Info("info msg", "3+5", 8)
	slog.Warn("warn msg")
	slog.Error("error msg", "status", 500, "err", net.ErrClosed)
	// LogAttrs 针对对应的数据优化，效率高。
	slog.LogAttrs(context.Background(), slog.LevelError, "attrs error log",
		slog.Int64("status", 500), slog.Any("err", net.ErrClosed), slog.Time("now", time.Now()))

}

func TestSlogJsonFormat(t *testing.T) {
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	jsonLogger.Debug("debug msg", "for", "bar", "now", time.Now())
	jsonLogger.Info("info msg", "3+5", 8)
	jsonLogger.Warn("warn msg")
}

func TestSlogContext(t *testing.T) {
	// ctx 怎么添加 reqid，再自定义一个 handler ？
	ctx := context.WithValue(context.Background(), "reqid", "abcd")
	slog.InfoContext(ctx, "context log", "a", 10)
}
