package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

func (app *application) recoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Connection", "close")

				app.serverErrorResponse(c, fmt.Errorf("%s", err))
				c.Abort()
			}
		}()
		c.Next()
	}
}

func (app *application) rateLimit() gin.HandlerFunc {
	limiter := rate.NewLimiter(2, 4)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			app.rateLimitExceededResponse(c)
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
