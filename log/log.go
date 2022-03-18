package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

type Level = zapcore.Level

const (
	InfoLevel   Level = zap.InfoLevel   // 0, default level
	WarnLevel   Level = zap.WarnLevel   // 1
	ErrorLevel  Level = zap.ErrorLevel  // 2
	DPanicLevel Level = zap.DPanicLevel // 3, used in development log
	PanicLevel  Level = zap.PanicLevel  // 4
	FatalLevel  Level = zap.FatalLevel  // 5
	DebugLevel  Level = zap.DebugLevel  // -1
)

type Field = zap.Field

func (l *Logger) Debug(ctx context.Context, msg string, fields ...Field) {
	uuid := ctx.Value("uuid")
	if uuid != nil {
		fields = append(fields, zap.String("uuid", ctx.Value("uuid").(string)))
	}
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...Field) {
	uuid := ctx.Value("uuid")
	if uuid != nil {
		fields = append(fields, zap.String("uuid", uuid.(string)))
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...Field) {
	uuid := ctx.Value("uuid")
	if uuid != nil {
		fields = append(fields, zap.String("uuid", uuid.(string)))
	}
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...Field) {
	uuid := ctx.Value("uuid")
	if uuid != nil {
		fields = append(fields, zap.String("uuid", uuid.(string)))
	}
	l.l.Error(msg, fields...)
}
func (l *Logger) DPanic(ctx context.Context, msg string, fields ...Field) {
	uuid := ctx.Value("uuid")
	if uuid != nil {
		fields = append(fields, zap.String("uuid", uuid.(string)))
	}
	l.l.DPanic(msg, fields...)
}
func (l *Logger) Panic(ctx context.Context, msg string, fields ...Field) {
	uuid := ctx.Value("uuid")
	if uuid != nil {
		fields = append(fields, zap.String("uuid", uuid.(string)))
	}
	l.l.Panic(msg, fields...)
}
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...Field) {
	uuid := ctx.Value("uuid")
	if uuid != nil {
		fields = append(fields, zap.String("uuid", uuid.(string)))
	}
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Sugar() *SugarLogger {
	return &SugarLogger{sl: l.l.Sugar()}
}

// function variables for all field types
// in github.com/uber-go/zap/field.go
var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	Int8p       = zap.Int8p
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Any         = zap.Any
	ErrorType   = zap.NamedError

	Info       = std.Info
	Warn       = std.Warn
	Error      = std.Error
	DPanic     = std.DPanic
	Panic      = std.Panic
	Fatal      = std.Fatal
	Debug      = std.Debug
	ToSugarlog = zap.SugaredLogger{}
)

// not safe for concurrent use
func ResetDefault(l *Logger) {
	std = l
	Info = std.Info
	Warn = std.Warn
	Error = std.Error
	DPanic = std.DPanic
	Panic = std.Panic
	Fatal = std.Fatal
	Debug = std.Debug
}

type Logger struct {
	l     *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
	level Level
}

var std = New(os.Stderr, InfoLevel)

func Default() *Logger {
	return std
}

type Option = zap.Option

var (
	WithCaller = zap.WithCaller
	//AddStacktrace = zap.AddStacktrace
	AddCallerSkip = zap.AddCallerSkip
)

// New create a new logger.
func New(writer io.Writer, level Level, opts ...Option) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000Z0700"))
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(writer),
		zapcore.Level(level),
	)
	logger := &Logger{
		l:     zap.New(core, opts...),
		level: level,
	}
	return logger
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}
	return nil
}

type SugarLogger struct {
	sl *zap.SugaredLogger
}

func (s *SugarLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	uuid := ctx.Value("uuid")
	if uuid != nil {
		keysAndValues = append(keysAndValues, zap.String("uuid", uuid.(string)))
	}
	s.sl.Infow(msg, keysAndValues...)
}

// Errorw
// @Description: create error log
// @receiver s
// @param ctx
// @param msg
// @param keysAndValues
func (s *SugarLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	uuid := ctx.Value("uuid")
	if uuid != nil {
		keysAndValues = append(keysAndValues, zap.String("uuid", uuid.(string)))
	}
	s.sl.Errorw(msg, keysAndValues...)
}
