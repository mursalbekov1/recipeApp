package main

import (
	"github.com/gin-gonic/gin"
	"go_recipe/internal/data"
	"net/http"
)

func (app *application) getAuthor(c *gin.Context) {
	authorId, err := app.readIDParam(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Автор не найден!"})
		return
	}

	author := data.Author{
		ID:             authorId,
		Name:           "Merei",
		Email:          "m_mursalbekov@kbtu.kz",
		Password:       "",
		Recipes:        []int64{1},
		RecipeAccesses: []int64{1, 2, 3, 4, 5},
	}

	err = app.writeJSON(c.Writer, http.StatusOK, Envelope{"author": author}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}
}
