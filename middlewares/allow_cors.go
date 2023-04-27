package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type AllowCors struct {
	mode         string
	allowOrigins []string
}

func NewAllowCors(mode string, allowOrigins []string) *AllowCors {
	return &AllowCors{mode: mode, allowOrigins: allowOrigins}
}

func (cors *AllowCors) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if cors.mode != gin.ReleaseMode {
			ctx.Header("Access-Control-Allow-Origin", origin)
		} else {
			originUrl, err := url.Parse(origin)
			if err == nil {
				for _, v := range cors.allowOrigins {
					if origin == v {
						ctx.Header("Access-Control-Allow-Origin", origin)
						break
					} else if v == "*" {
						ctx.Header("Access-Control-Allow-Origin", origin)
						break
					} else if strings.Contains(v, "*") {
						matched, err := regexp.MatchString(v, originUrl.Host)
						if err == nil && matched {
							ctx.Header("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}
			}
		}
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,Accept,Authorization,X-Requested-With,X-XSRF-TOKEN,x-csrf-token,Cache-Control,crm-token")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}

		ctx.Next()
	}
}
