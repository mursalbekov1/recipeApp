package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_recipe/internal/data"
	"go_recipe/internal/validator"
	"net/http"
	"strconv"
	"time"
)

var version = "1.0.0"

func (app *application) getRecipe(c *gin.Context) {
	recipeID, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return
	}

	recipe, err := app.models.Recipe.Get(recipeID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	err = app.writeJSON(c.Writer, http.StatusOK, Envelope{"recipe": recipe}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}

}

func (app *application) addRecipe(c *gin.Context) {
	var input struct {
		Title         string   `json:"title"`
		Description   string   `json:"description"`
		Ingredients   []string `json:"ingredients"`
		Steps         []string `json:"steps"`
		Author        int64    `json:"author"`
		Collaborators []int64  `json:"collaborators"`
	}

	if err := app.readJSON(c, &input); err != nil {
		app.badRequestResponse(c, err)
		return
	}

	recipe := &data.Recipe{
		Title:         input.Title,
		Description:   input.Description,
		Ingredients:   input.Ingredients,
		Steps:         input.Steps,
		Author:        input.Author,
		Collaborators: input.Collaborators,
	}

	v := validator.New()

	if data.ValidateRecipe(v, recipe); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	err := app.models.Recipe.Insert(recipe)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/recipe/%d", recipe.ID))


	err = app.writeJSON(c.Writer, http.StatusCreated, Envelope{"Recipe": recipe}, headers)

	//c.JSON(http.StatusOK, gin.H{"data": input})

}

func (app *application) updateRecipe(c *gin.Context) {
	id, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return
	}

	recipe, err := app.models.Recipe.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	if c.GetHeader("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(recipe.Version), 32) != c.GetHeader("X-Expected-Version") {
			app.editConflictResponse(c)
			return
		}
	}

	var input struct {
		Title         *string  `json:"title"`
		Description   *string  `json:"description"`
		Ingredients   []string `json:"ingredients"`
		Steps         []string `json:"steps"`
		Collaborators []int64  `json:"collaborators"`
	}

	err = app.readJSON(c, &input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	if input.Title != nil {
		recipe.Title = *input.Title
	}

	if input.Description != nil {
		recipe.Description = *input.Description
	}

	if input.Ingredients != nil {
		recipe.Ingredients = input.Ingredients
	}

	if input.Steps != nil {
		recipe.Steps = input.Steps
	}

	if input.Collaborators != nil {
		recipe.Collaborators = input.Collaborators
	}

	v := validator.New()

	if data.ValidateRecipe(v, recipe); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	err = app.models.Recipe.Update(recipe)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	err = app.writeJSON(c.Writer, http.StatusOK, Envelope{"Recipe": recipe}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}
}

func (app *application) deleteRecipe(c *gin.Context) {
	id, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return
	}

	err = app.models.Recipe.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	err = app.writeJSON(c.Writer, http.StatusOK, Envelope{"Recipe": "Recipe was successful deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}
}

func (app *application) getRecipeList(c *gin.Context) {
	var input struct {
		Title         string
		Time          time.Time
		Description   string
		Ingredients   []string
		Steps         []string
		Author        int
		Collaborators []int
		data.Filters
	}

	v := validator.New()

	qs := c.Request.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Description = app.readString(qs, "description", "")
	input.Ingredients = app.readCSV(qs, "ingredients", []string{})
	input.Steps = app.readCSV(qs, "steps", []string{})
	input.Author = app.readInt(qs, "author", 1, v)
	input.Collaborators = app.read(qs, "collaborators", []int{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "author", "-id", "-title", "-author"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	recipes, metadata, err := app.models.Recipe.GetAll(input.Title, input.Ingredients, input.Filters)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	err = app.writeJSON(c.Writer, http.StatusOK, Envelope{"recipes": recipes, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}

}

//func (app *application) healthcheckHandler(c *gin.Context) {
//	env := Envelope{
//		"status": "available",
//		"system_info": map[string]string{
//			"environment": app.config.env,
//			"version":     version,
//		},
//	}
//	err := app.writeJSON(c.Writer, http.StatusOK, env, nil)
//	if err != nil {
//		app.serverErrorResponse(c, err)
//	}
//}
