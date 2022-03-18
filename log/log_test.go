package log

import (
	"context"
	"os"
	"testing"

	hctx "git.zhwenxue.com/zhgo/gocontrib/context"
)

func TestNewDebug(t *testing.T) {
	logg := New(os.Stdout, DebugLevel, WithCaller(false), AddCallerSkip(1))
	ctx := hctx.GetContext(context.Background(), "")
	logg.Debug(ctx, "Debug", String("Debug", "ok"))
	logg.Info(ctx, "Info", String("Info", "ok"))
	logg.Error(ctx, "Error", String("Error", "ok"))
}

func TestNewInfo(t *testing.T) {
	logg := New(os.Stdout, InfoLevel, WithCaller(false), AddCallerSkip(1))
	ctx := hctx.GetContext(context.Background(), "")
	logg.Debug(ctx, "Debug", String("Debug", "ok")) //不输出
	logg.Info(ctx, "Info", String("Info", "ok"))
	logg.Error(ctx, "Error", String("Error", "ok"))
}

func TestNewERROR(t *testing.T) {
	logg := New(os.Stdout, ErrorLevel, WithCaller(false), AddCallerSkip(1))
	ctx := hctx.GetContext(context.Background(), "")
	logg.Debug(ctx, "Debug", String("Debug", "ok")) //不输出
	logg.Info(ctx, "Info", String("Info", "ok"))    //不输出
	logg.Error(ctx, "Error", String("Error", "ok"))
}
