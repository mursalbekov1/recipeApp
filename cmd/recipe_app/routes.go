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
		//v1.GET("/getRecipeList", app.getRecipeList)
		v1.GET("/author/:id", app.getAuthor)
		v1.GET("/recipe/:id", app.getRecipe)
		//v1.GET("/check", app.healthcheckHandler)
		v1.POST("/recipe", app.addRecipe)
		v1.PUT("/recipe/:id", app.updateRecipe)
	}

	return router
}
