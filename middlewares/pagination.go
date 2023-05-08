package middlewares

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Pagination struct{}

func (ac *Pagination) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, err := strconv.ParseUint(ctx.Query("page"), 10, 64)
		if err != nil || page == 0 {
			page = 1
		}
		pageSize, err := strconv.ParseUint(ctx.Query("page_size"), 10, 64)
		if err != nil || pageSize == 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}
		ctx.Set("page", page)
		ctx.Set("page_size", pageSize)
		ctx.Next()
	}
}
