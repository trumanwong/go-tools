package middlewares

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// Pagination is a struct that represents the pagination configuration.
// It contains the default page size and the maximum page size.
type Pagination struct {
	// defaultPageSize is the default number of items per page when the "page_size" query parameter is not provided or is not a valid positive integer.
	defaultPageSize uint64
	// maxPageSize is the maximum number of items per page when the "page_size" query parameter is greater than this value.
	maxPageSize uint64
}

// NewPagination is a function that creates a new Pagination object.
// It takes the default page size and the maximum page size as parameters,
// and returns a pointer to the created Pagination object.
func NewPagination(defaultPageSize, maxPageSize uint64) *Pagination {
	return &Pagination{defaultPageSize: defaultPageSize, maxPageSize: maxPageSize}
}

// Handle is a method of Pagination that returns a gin.HandlerFunc for handling pagination.
// The returned gin.HandlerFunc retrieves the "page" and "page_size" query parameters from the request.
// If the "page" query parameter is not provided or is not a valid positive integer, it defaults to 1.
// If the "page_size" query parameter is not provided or is not a valid positive integer, it defaults to the default page size.
// If the "page_size" query parameter is greater than the maximum page size, it is set to the maximum page size.
// The gin.HandlerFunc then sets the "page" and "page_size" values in the gin.Context,
// and calls the next handler in the chain.
func (p *Pagination) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse the "page" query parameter.
		page, err := strconv.ParseUint(ctx.Query("page"), 10, 64)
		// If the "page" query parameter is not provided or is not a valid positive integer, default to 1.
		if err != nil || page == 0 {
			page = 1
		}
		// Parse the "page_size" query parameter.
		pageSize, err := strconv.ParseUint(ctx.Query("page_size"), 10, 64)
		// If the "page_size" query parameter is not provided or is not a valid positive integer, default to the default page size.
		if err != nil || pageSize == 0 {
			pageSize = p.defaultPageSize
		}
		// If the "page_size" query parameter is greater than the maximum page size, set it to the maximum page size.
		if pageSize > p.maxPageSize {
			pageSize = p.maxPageSize
		}
		// Set the "page" and "page_size" values in the gin.Context.
		ctx.Set("page", page)
		ctx.Set("page_size", pageSize)
		// Call the next handler in the chain.
		ctx.Next()
	}
}
