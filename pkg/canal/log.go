package canal

import (
	"context"
	"fmt"
)

type ILogger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

type defaultLogger struct{}

func (l defaultLogger) Debug(ctx context.Context, msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
}
func (l defaultLogger) Info(ctx context.Context, msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
}
func (l defaultLogger) Warn(ctx context.Context, msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
}
func (l defaultLogger) Error(ctx context.Context, msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
}
