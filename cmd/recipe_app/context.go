package main

import (
	"github.com/gin-gonic/gin"
	"go_recipe/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(c *gin.Context, user *data.User) {
	c.Set("user", user)
}

func (app *application) contextGetUser(c *gin.Context) *data.User {
	user, ok := c.Get("user")
	if !ok {
		panic("missing user value in request context")
	}
	return user.(*data.User)
}
