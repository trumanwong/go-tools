package middlewares

import "github.com/gin-gonic/gin"

// Middleware is an interface that represents a middleware in the Gin framework.
// It contains a single method, Handle, which returns a gin.HandlerFunc.
// The gin.HandlerFunc returned by Handle is a function that handles an HTTP request and can be registered with a gin.Engine or a gin.RouterGroup.
type Middleware interface {
	// Handle is a method that returns a gin.HandlerFunc.
	// The returned gin.HandlerFunc is a function that takes a gin.Context as a parameter,
	// performs some operations (such as logging, authentication, etc.), and then calls the next handler in the chain.
	// If the gin.HandlerFunc decides to terminate the request handling, it can call gin.Context.Abort or gin.Context.AbortWithStatus.
	Handle() gin.HandlerFunc
}
