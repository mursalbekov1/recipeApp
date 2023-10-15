package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	router.NoRoute(app.notFoundResponse)
	router.NoMethod(app.methodNotAllowedResponse)

	// Define your routes and associate them with handlers
	v1 := router.Group("/v1")
	{
		v1.GET("/addRecipe", app.addRecipe)
		v1.GET("/recipe/:id", app.getRecipe)
		v1.GET("/check", app.healthcheckHandler)
		//v1.PUT("/v1/addRecipe", app.AddRecipe)
		//v1.DELETE("/v1/deleteRecipe/:id", app.DeleteRecipe)
		//v1.PUT("/v1/updateRecipe/:id", app.UpdateRecipe)
	}

	return router
}
