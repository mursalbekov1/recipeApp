package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_recipe/internal/data"
	"go_recipe/internal/validator"
	"net/http"
)

const version = "1.0.0"

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

//func (app *application) getRecipeList(c *gin.Context) {
//	//jsonData, err := json.Marshal(jsonR)
//	//if err != nil {
//	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при кодировании JSON"})
//	//	return
//	//}
//	//
//	//c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
//}

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

	err = app.writeJSON(c.Writer, http.StatusCreated, Envelope{"recipe": recipe}, headers)

	//c.JSON(http.StatusOK, gin.H{"data": input})

}

func (app *application) healthcheckHandler(c *gin.Context) {
	env := Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(c.Writer, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}

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

	var input struct {
		Title         string   `json:"title"`
		Description   string   `json:"description"`
		Ingredients   []string `json:"ingredients"`
		Steps         []string `json:"steps"`
		Collaborators []int64  `json:"collaborators"`
	}

	err = app.readJSON(c, &input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	recipe.Title = input.Title
	recipe.Description = input.Description
	recipe.Ingredients = input.Ingredients
	recipe.Steps = input.Steps
	recipe.Collaborators = input.Collaborators

	v := validator.New()

	if data.ValidateRecipe(v, recipe); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	err = app.models.Recipe.Update(recipe)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	err = app.writeJSON(c.Writer, http.StatusOK, Envelope{"Recipe": recipe}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}
}
