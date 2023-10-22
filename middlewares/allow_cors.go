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
	allowHeaders string
	allowMethods string
	allowOrigins []string
}

func NewAllowCors(mode, allowHeaders, allowMethods string, allowOrigins []string) *AllowCors {
	return &AllowCors{mode: mode, allowHeaders: allowHeaders, allowMethods: allowMethods, allowOrigins: allowOrigins}
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
					if originUrl.Host == v {
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
		ctx.Header("Access-Control-Allow-Methods", cors.allowMethods)
		ctx.Header("Access-Control-Allow-Headers", cors.allowHeaders)
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}

		ctx.Next()
	}
}
