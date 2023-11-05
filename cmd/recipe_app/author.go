package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_recipe/internal/data"
	"go_recipe/internal/validator"
	"net/http"
)

func (app *application) getAuthor(c *gin.Context) {
	authorID, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return
	}

	author, err := app.models.Author.Get(authorID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	err = app.writeJSON(c.Writer, http.StatusOK, Envelope{"author": author}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}
}

func (app *application) addAuthor(c *gin.Context) {
	var input struct {
		Name           string  `json:"name"`
		Email          string  `json:"email"`
		Password       string  `json:"password"`
		Recipes        []int64 `json:"recipes"`
		RecipeAccesses []int64 `json:"steps"`
	}

	if err := app.readJSON(c, &input); err != nil {
		app.badRequestResponse(c, err)
		return
	}

	author := &data.Author{
		Name:           input.Name,
		Email:          input.Email,
		Password:       input.Password,
		Recipes:        input.Recipes,
		RecipeAccesses: input.RecipeAccesses,
	}

	v := validator.New()

	if data.ValidateAuthor(v, author); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	err := app.models.Author.Insert(author)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/author/%d", author.ID))

	err = app.writeJSON(c.Writer, http.StatusCreated, Envelope{"author": author}, headers)

	//c.JSON(http.StatusOK, gin.H{"data": input})

}

func (app *application) updateAuthor(c *gin.Context) {
	id, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return
	}

	author, err := app.models.Author.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	var input struct {
		Name           *string `json:"name"`
		Email          *string `json:"email"`
		Password       *string `json:"password"`
		Recipes        []int64 `json:"recipes"`
		RecipeAccesses []int64 `json:"recipeaccesses"`
	}

	err = app.readJSON(c, &input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	if input.Name != nil {
		author.Name = *input.Name
	}

	if input.Email != nil {
		author.Email = *input.Email
	}

	if input.Password != nil {
		author.Password = *input.Password
	}

	if input.Recipes != nil {
		author.Recipes = input.Recipes
	}

	if input.RecipeAccesses != nil {
		author.RecipeAccesses = input.RecipeAccesses
	}

	v := validator.New()

	if data.ValidateAuthor(v, author); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	err = app.models.Author.Update(author)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	err = app.writeJSON(c.Writer, http.StatusOK, Envelope{"author": author}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}
}

func (app *application) deleteAuthor(c *gin.Context) {
	id, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return
	}

	err = app.models.Author.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	err = app.writeJSON(c.Writer, http.StatusOK, Envelope{"Author": "Author was successful deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}
}
