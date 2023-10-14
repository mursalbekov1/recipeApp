package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	router.GET("/v1/addRecipe", addRecipe)
	router.POST("/v1/recipe/:id", getRecipe)
	//router.PUT("/v1/addRecipe", recipe.AddRecipe)
	//router.DELETE("/v1/deleteRecipe/:id", recipe.DeleteRecipe)
	//router.PUT("/v1/updateRecipe/:id", recipe.UpdateRecipe)

	return router
}
