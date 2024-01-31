package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	router.NoRoute(app.notFoundResponse)
	router.NoMethod(app.methodNotAllowedResponse)

	v1 := router.Group("/v1")
	{
		v1.Use(app.recoverPanic())
		v1.Use(app.enableCORS())
		v1.Use(app.rateLimit())
		v1.Use(app.authenticate())

		v1.GET("/recipe/:id", app.requirePermission("recipe:read", app.getRecipe))
		v1.POST("/recipe", app.requirePermission("recipe:write", app.addRecipe))
		v1.PATCH("/recipe/:id", app.requirePermission("recipe:write", app.updateRecipe))
		v1.DELETE("/recipe/:id", app.requirePermission("recipe:write", app.deleteRecipe))
		v1.GET("/recipe", app.requirePermission("recipe:read", app.getRecipeList))

		//v1.GET("/healthcheck", app.healthcheckHandler)

		v1.POST("/users", app.registerUserHandler)
		v1.PUT("/users/activated", app.activateUserHandler)

		v1.POST("/tokens/authentication", app.createAuthenticationTokenHandler)
	}

	return router
}
