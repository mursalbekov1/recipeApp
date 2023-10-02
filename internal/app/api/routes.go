package app

import (
	"github.com/gin-gonic/gin"
	"recipeApp/internal/app/recipe"
)

func Routes() *gin.Engine {
	router := gin.Default()

	router.GET("/v1/recipeList", recipe.GetRecipesList)
	router.GET("/v1/recipe/:id", recipe.GetRecipe)
	router.PUT("/v1/addRecipe", recipe.AddRecipe)
	router.DELETE("/v1/deleteRecipe/:id", recipe.DeleteRecipe)
	router.PUT("/v1/updateRecipe/:id", recipe.UpdateRecipe)

	return router
}
