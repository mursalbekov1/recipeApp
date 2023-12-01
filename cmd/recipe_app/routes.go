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
		v1.Use(app.rateLimit())
		v1.Use(app.authenticate())

		v1.GET("/author/:id", app.getAuthor)
		v1.POST("/author", app.addAuthor)
		v1.PATCH("/author/:id", app.updateAuthor)
		v1.DELETE("/author/:id", app.deleteAuthor)

		v1.GET("/recipe/:id", app.requireActivatedUser(), app.getRecipe)
		v1.POST("/recipe", app.requireActivatedUser(), app.addRecipe)
		v1.PATCH("/recipe/:id", app.requireActivatedUser(), app.updateRecipe)
		v1.DELETE("/recipe/:id", app.requireActivatedUser(), app.deleteRecipe)
		v1.GET("/recipe", app.requireActivatedUser(), app.getRecipeList)

		v1.GET("/healthcheck", app.healthcheckHandler)

		v1.POST("/users", app.registerUserHandler)
		v1.PUT("/users/activated", app.activateUserHandler)

		v1.POST("/tokens/authentication", app.createAuthenticationTokenHandler)
	}

	return router
}
