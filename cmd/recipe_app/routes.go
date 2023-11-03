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
		v1.GET("/author/:id", app.getAuthor)
		v1.GET("/recipe/:id", app.getRecipe)
		v1.POST("/recipe", app.addRecipe)
		v1.PATCH("/recipe/:id", app.updateRecipe)
		v1.DELETE("/recipe/:id", app.deleteRecipe)
	}

	return router
}
