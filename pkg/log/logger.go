package log

import (
	"context"
	"fmt"

	"github.com/illidaris/core"
	"github.com/illidaris/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var l *zap.Logger

func NewLogger() *zap.Logger {
	l = zap.L().WithOptions(zap.AddCallerSkip(2))
	return l
}

func Debug(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, fmt.Sprintf(msg, args...), zapcore.DebugLevel)
}

func Info(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, fmt.Sprintf(msg, args...), zapcore.InfoLevel)
}

func Warn(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, fmt.Sprintf(msg, args...), zapcore.WarnLevel)
}

func Error(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, fmt.Sprintf(msg, args...), zapcore.ErrorLevel)
}

func Log(ctx context.Context, msg string, lvl zapcore.Level, fields ...zapcore.Field) {
	base := FieldsFromCtx(ctx)
	if len(fields) > 0 {
		base = append(base, fields...)
	}
	l.Log(lvl, msg, base...)
}

func FieldsFromCtx(ctx context.Context) []zap.Field {
	return []zap.Field{
		buildZapField(ctx, core.TraceID),
		buildZapField(ctx, core.SessionID),
		buildZapField(ctx, core.Action),
		buildZapField(ctx, core.Step),
	}
}

// buildZapField use core meta build field
func buildZapField(ctx context.Context, key core.MetaData) zap.Field {
	return zap.String(key.String(), key.GetString(ctx))
}

func NewCtx(ctx context.Context, sessionID, traceID string) context.Context {
	ctx = core.SessionID.SetString(ctx, sessionID)
	ctx = core.TraceID.SetString(ctx, traceID)
	ctx = logger.NewContext(ctx,
		zap.String(core.SessionID.String(), sessionID),
		zap.String(core.TraceID.String(), traceID))
	return ctx
}
