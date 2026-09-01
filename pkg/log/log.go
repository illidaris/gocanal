package log

import (
	"context"
	"fmt"

	"github.com/illidaris/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log *zap.Logger
}

func NewRestLogger() *Logger {
	return &Logger{
		log: zap.L().WithOptions(zap.AddCallerSkip(1)),
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.log.Debug(fmt.Sprintf(msg, args...), WithTrace(ctx)...)
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.log.Info(fmt.Sprintf(msg, args...), WithTrace(ctx)...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.log.Warn(fmt.Sprintf(msg, args...), WithTrace(ctx)...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.log.Error(fmt.Sprintf(msg, args...), WithTrace(ctx)...)
}

func WithTrace(ctx context.Context) []zapcore.Field {
	traceID := core.TraceID.GetString(ctx)
	sessionID := core.SessionID.GetString(ctx)
	return []zapcore.Field{
		zap.String(core.TraceID.String(), traceID),
		zap.String(core.SessionID.String(), sessionID),
	}
}
