package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/trumanwong/go-tools/log"
	"time"
)

type Logger struct {
	Logger *log.Logger
}

func NewLogger(logger *log.Logger) *Logger {
	return &Logger{Logger: logger}
}

func (l *Logger) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		l.Logger.WithFields(logrus.Fields{
			// 请求方式
			"method": ctx.Request.Method,
			// 请求路由
			"uri": ctx.Request.RequestURI,
			// 请求ip
			"client_ip": ctx.ClientIP(),
			// 请求头
			"header": ctx.Request.Header,
			// 返回code
			"status_code": ctx.Writer.Status(),
			// 执行时间
			"execute_time": endTime.Sub(startTime) / time.Millisecond,
			"created_at":   time.Now(),
		})
	}
}
