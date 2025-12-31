package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type tracing struct {
	Key *string
}

// NewTracing 创建一个追踪中间件，将 traceId 存储到 gin.Context 中
// 使用时可以通过 ctx.GetString(key) 获取 traceId
func NewTracing(key *string) Middleware {
	return &tracing{
		Key: key,
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
		// 将 traceId 存储到 gin.Context 中，而不是修改全局 Logger
		ctx.Set(p.getKey(), traceId)
		ctx.Header(p.getKey(), traceId)
		ctx.Next()
	}
}
