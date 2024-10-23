package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/trumanwong/go-tools/log"
)

type tracing struct {
	Logger *log.Logger
	Key    *string
}

func NewTracing(key *string, logger *log.Logger) Middleware {
	return &tracing{
		Logger: logger,
		Key:    key,
	}
}

func (p *tracing) getKey() string {
	if p.Key != nil {
		return *p.Key
	}
	return "X-Trace-Id"
}

func (p *tracing) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.GetHeader(p.getKey())
		if traceId == "" {
			traceId = uuid.New().String()
		}
		p.Logger.WithTraceId(ctx, traceId)
		ctx.Header(p.getKey(), traceId)
		ctx.Next()
	}
}
