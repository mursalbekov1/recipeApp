package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"recipeApp/internal/app/recipe"
	"time"
)

const PORT = ":8080"

func main() {
	router := gin.Default()
	router.GET("/recipeList", recipe.GetRecipesList)
	router.GET("/recipe/:id", recipe.GetRecipe)

	serv := &http.Server{
		Addr:        PORT,
		Handler:     router,
		ReadTimeout: 3 * time.Second,
	}

	err := serv.ListenAndServe()
	if err != nil {
		return
	}
}
