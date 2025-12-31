package log

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestLogger_Debug(t *testing.T) {
	logger := NewLogger(&Options{
		TraceKey:  nil,
		Formatter: nil,
		Output:    nil,
	})
	// 使用新的 WithTraceId 方法（返回 Entry）
	logger.WithTraceId(uuid.New().String()).Info("hello with traceId")

	// 也可以使用 WithContext 方法
	ctx := context.WithValue(context.Background(), "X-Trace-Id", uuid.New().String())
	logger.WithContext(ctx).Info("hello with context 1")
	logger.WithContext(ctx).Info("hello with context 2")

	// 原有的无 traceId 日志方法仍然可用
	logger.Info("hello without traceId")
}
