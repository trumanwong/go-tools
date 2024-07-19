package log

import (
	"context"
	"github.com/google/uuid"
	"testing"
)

func TestLogger_Debug(t *testing.T) {
	logger := NewLogger(&Options{
		TraceKey:  nil,
		Formatter: nil,
		Output:    nil,
	})
	logger.WithTraceId(context.Background(), uuid.New().String())
	logger.Info("hello")
}
