package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_recipe/internal/data"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(c *gin.Context) *data.User {
	user, _ := c.Request.Context().Value(userContextKey).(*data.User)
	return user
}
