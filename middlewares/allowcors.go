package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// AllowCors is a struct that represents the CORS (Cross-Origin Resource Sharing) configuration.
// It contains the mode (debug or release), the allowed headers, the allowed methods, and the allowed origins.
type allowCors struct {
	mode         string
	allowHeaders string
	allowMethods string
	allowOrigins *[]string
}

// NewAllowCors is a function that creates a new AllowCors object.
// It takes the mode, the allowed headers, the allowed methods, and the allowed origins as parameters,
// and returns a pointer to the created AllowCors object.
func NewAllowCors(mode, allowHeaders, allowMethods string, allowOrigins *[]string) Middleware {
	return &allowCors{mode: mode, allowHeaders: allowHeaders, allowMethods: allowMethods, allowOrigins: allowOrigins}
}

// Handle is a method of AllowCors that returns a gin.HandlerFunc for handling CORS.
// The returned gin.HandlerFunc sets the Access-Control-Allow-Origin, Access-Control-Allow-Methods,
// Access-Control-Allow-Headers, and Access-Control-Allow-Credentials headers according to the AllowCors configuration.
// If the request method is OPTIONS, the gin.HandlerFunc aborts the request with the HTTP status code 204 (No Content).
// Otherwise, it calls the next handler in the chain.
func (cors *allowCors) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if cors.mode != gin.ReleaseMode {
			ctx.Header("Access-Control-Allow-Origin", origin)
		} else {
			originUrl, err := url.Parse(origin)
			if err == nil && cors.allowOrigins != nil {
				for _, v := range *cors.allowOrigins {
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
