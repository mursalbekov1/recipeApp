package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_recipe/internal/data"
	"go_recipe/internal/validator"
	"golang.org/x/time/rate"
	"net"
	"strings"
	"sync"
	"time"
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
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(c, err)
				return
			}
			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}
			clients[ip].lastSeen = time.Now()
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(c)
				c.Abort()
				return
			}
			mu.Unlock()
		}
		c.Next()
	}
}

func (app *application) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Vary", "Authorization")

		authorizationHeader := c.GetHeader("Authorization")

		if authorizationHeader == "" {
			app.contextSetUser(c, data.AnonymousUser)
			c.Next()
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(c)
			c.Abort()
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(c)
			c.Abort()
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(c)
			default:
				app.serverErrorResponse(c, err)
			}
			c.Abort()
			return
		}

		app.contextSetUser(c, user)
		c.Next()
	}
}

func (app *application) requireAuthenticatedUser(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := app.contextGetUser(c)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(c)
			c.Abort()
			return
		}
		next(c)
	}
}

func (app *application) requireActivatedUser(next gin.HandlerFunc) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		user := app.contextGetUser(c)

		if !user.Activated {
			app.inactiveAccountResponse(c)
			c.Abort()
			return
		}

		next(c)
	}

	return app.requireAuthenticatedUser(fn)
}

func (app *application) requirePermission(code string, next gin.HandlerFunc) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		user := app.contextGetUser(c)

		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(c, err)
			return
		}

		if !permissions.Include(code) {
			app.notPermittedResponse(c)
			return
		}

		next(c)
	}

	return app.requireActivatedUser(fn)
}
