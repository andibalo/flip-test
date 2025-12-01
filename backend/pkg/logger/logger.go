package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/andibalo/flip-test/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LevelTrace string = "trace"
	LevelDebug string = "debug"
	LevelInfo  string = "info"
	LevelWarn  string = "warn"
	LevelError string = "error"
	LevelFatal string = "fatal"
	LevelPanic string = "panic"

	// FieldKeyTraceID is the field key for the trace ID.
	FieldKeyTraceID = "trace.id"

	// FieldKeyTransactionID is the field key for the transaction ID.
	FieldKeyTransactionID = "transaction.id"

	// FieldKeySpanID is the field key for the span ID.
	FieldKeySpanID = "span.id"
)

type Logger interface {
	TraceWithContext(ctx context.Context, msg string, v ...zapcore.Field)
	DebugWithContext(ctx context.Context, msg string, v ...zapcore.Field)
	InfoWithContext(ctx context.Context, msg string, v ...zapcore.Field)
	WarnWithContext(ctx context.Context, msg string, v ...zapcore.Field)
	ErrorWithContext(ctx context.Context, msg string, v ...zapcore.Field)
	FatalWithContext(ctx context.Context, msg string, v ...zapcore.Field)
	PanicWithContext(ctx context.Context, msg string, v ...zapcore.Field)
	Trace(msg string, v ...zapcore.Field)
	Debug(msg string, v ...zapcore.Field)
	Info(msg string, v ...zapcore.Field)
	Warn(msg string, v ...zapcore.Field)
	Error(msg string, v ...zapcore.Field)
	Fatal(msg string, v ...zapcore.Field)
	Panic(msg string, v ...zapcore.Field)
}

var (
	once      = sync.Once{}
	appLogger Logger
)

func GetLogger() Logger {
	return appLogger
}

func InitLogger(cfg config.Config) Logger {
	once.Do(func() {
		opts := Options{
			DefaultFields: map[string]string{
				"service.name":    cfg.AppAddress(),
				"service.version": cfg.AppVersion(),
				"service.env":     cfg.AppEnv(),
			},
			ContextFields: map[string]string{
				"path":        "path",
				"method":      "method",
				"status_code": "status_code",
				"status":      "status",
				"error":       "error",
				"user_id":     "x-user-id",
				"user_email":  "x-user-email",
				"client_ip":   "x-forwarded-for",
				"payload":     "payload",
				"x-client-id": "x-client-id",
				"topic":       "topic",
				"broker":      "broker",
				"trace.id":    "trace.id",
				"span.id":     "span.id",
				"request_id":  "request_id",
			},
			Level:     LevelInfo,
			HookLevel: LevelError,
		}

		appLogger = initLogger(opts)
	})

	return appLogger
}

func GetLoggerWithDefaultOptions() Logger {
	return initLogger(DefaultLoggerOption)
}

type logger struct {
	mu       *sync.RWMutex
	log      *zap.Logger
	logEntry *zapcore.Entry
	opt      Options
}

type Options struct {
	Output        string
	Formatter     string
	Level         string
	ContextFields map[string]string
	DefaultFields map[string]string
	CustomFields  map[string]interface{}
	CustomWriter  io.Writer
	Hook          string
	HookLevel     string
	//Hooks         []Hook
}

var (
	DefaultLoggerOption = Options{
		DefaultFields: map[string]string{
			"service.name":    "transaction-service",
			"service.version": "0.1",
			"service.env":     "dev",
		},
		ContextFields: map[string]string{
			"path":        "path",
			"method":      "method",
			"status_code": "status_code",
			"status":      "status",
			"error":       "error",
			"user_id":     "x-user-id",
			"user_email":  "x-user-email",
			"client_ip":   "x-forwarded-for",
			"payload":     "payload",
			"x-client-id": "x-client-id",
			"topic":       "topic",
			"broker":      "broker",
			"trace.id":    "trace.id",
			"span.id":     "span.id",
		},
		Level:     LevelInfo,
		HookLevel: LevelError,
	}
)

func InitZapLogger() *zap.Logger {

	var zapLogger *zap.Logger

	config := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths: []string{"stdout"}, // Log to file
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			MessageKey:     "msg",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
	}

	zapLogger, err := config.Build()

	if err != nil {
		panic(fmt.Sprintf("logger initialization failed %v", err))
	}

	if os.Getenv("APP_ENV") == "DEV" {
		zapLogger, err = zap.NewDevelopment()

		if err != nil {
			panic(fmt.Sprintf("logger initialization failed %v", err))
		}
	}

	zapLogger.Info("logger started")

	defer zapLogger.Sync()

	return zapLogger
}

func initLogger(opt Options) Logger {

	zapLog := InitZapLogger()

	lg := &logger{
		mu:  &sync.RWMutex{},
		log: zapLog,
		opt: opt,
	}

	lg.setDefaultOptions()

	return lg
}

func (l *logger) setDefaultOptions() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.opt.Output == "" {
		//never put default to discard, error will not be displayed!
		l.opt.Output = DefaultLoggerOption.Output
	}
	if l.opt.Formatter == "" {
		l.opt.Formatter = DefaultLoggerOption.Formatter
	}
	if l.opt.Level == "" {
		l.opt.Level = DefaultLoggerOption.Level
	}
	if l.opt.ContextFields == nil {
		l.opt.ContextFields = DefaultLoggerOption.ContextFields
	}
	if l.opt.DefaultFields == nil {
		l.opt.DefaultFields = DefaultLoggerOption.DefaultFields
	}
	if l.opt.Hook == "" {
		l.opt.Hook = DefaultLoggerOption.Hook
	}
	if l.opt.HookLevel == "" {
		l.opt.HookLevel = DefaultLoggerOption.HookLevel
	}
}

func (l *logger) parseContextFields(ctx context.Context) *zap.Logger {
	doLog := l.log
	if ctx != nil {
		for k, v := range l.opt.ContextFields {
			if val := ctx.Value(v); val != nil {
				doLog = doLog.With(zap.Any(k, val))
			}
		}

		if len(l.opt.CustomFields) != 0 {
			for k, v := range l.opt.CustomFields {
				if val := ctx.Value(v); val != nil {
					doLog = doLog.With(zap.Any(k, val))
				}
			}
		}

		//switch l.opt.Hook {
		//case APM:
		//	doLog = doLog.WithFields(apmlogrus.TraceContext(ctx))
		//case OTEL:
		//	span := trace.SpanFromContext(ctx)
		//	traceID := span.SpanContext().TraceID().String()
		//	spanID := span.SpanContext().SpanID().String()
		//
		//	fields := lr.Fields{
		//		FieldKeyTraceID: traceID,
		//		FieldKeySpanID:  spanID,
		//	}
		//
		//	doLog = doLog.WithFields(fields)
		//}
	}

	return doLog
}

func (l *logger) TraceWithContext(ctx context.Context, msg string, v ...zapcore.Field) {
	//l.parseContextFields(ctx).Trace(msg,v...)
}

func (l *logger) Trace(msg string, v ...zapcore.Field) {
	l.TraceWithContext(context.TODO(), msg, v...)
}

func (l *logger) DebugWithContext(ctx context.Context, msg string, v ...zapcore.Field) {
	l.parseContextFields(ctx).Debug(msg, v...)
}

func (l *logger) Debug(msg string, v ...zapcore.Field) {
	l.DebugWithContext(context.TODO(), msg, v...)
}

func (l *logger) InfoWithContext(ctx context.Context, msg string, v ...zapcore.Field) {
	l.parseContextFields(ctx).Info(msg, v...)
}

func (l *logger) Info(msg string, v ...zapcore.Field) {
	l.InfoWithContext(context.TODO(), msg, v...)
}

func (l *logger) WarnWithContext(ctx context.Context, msg string, v ...zapcore.Field) {
	l.parseContextFields(ctx).Warn(msg, v...)
}

func (l *logger) Warn(msg string, v ...zapcore.Field) {
	l.WarnWithContext(context.TODO(), msg, v...)
}

func (l *logger) ErrorWithContext(ctx context.Context, msg string, v ...zapcore.Field) {
	l.parseContextFields(ctx).Error(msg, v...)
}

func (l *logger) Error(msg string, v ...zapcore.Field) {
	l.ErrorWithContext(context.TODO(), msg, v...)
}

func (l *logger) FatalWithContext(ctx context.Context, msg string, v ...zapcore.Field) {
	l.parseContextFields(ctx).Fatal(msg, v...)
}

func (l *logger) Fatal(msg string, v ...zapcore.Field) {
	l.FatalWithContext(context.TODO(), msg, v...)
}

func (l *logger) PanicWithContext(ctx context.Context, msg string, v ...zapcore.Field) {
	l.parseContextFields(ctx).Panic(msg, v...)
}

func (l *logger) Panic(msg string, v ...zapcore.Field) {
	l.PanicWithContext(context.TODO(), msg, v...)
}
