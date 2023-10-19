package main

import (
	"github.com/gin-gonic/gin"
	"go_recipe/internal/data"
	"go_recipe/internal/validator"
	"net/http"
)

const version = "1.0.0"

func (app *application) getRecipe(c *gin.Context) {
	recipeID, err := app.readIDParam(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Рецепт не найден!"})
		return
	}

	//c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Отображение информации о рецепте %d", recipeID)})

	recipe := data.Recipe{
		ID:            recipeID,
		Runtime:       102,
		Title:         "Паста с помидорами",
		Description:   "Простой рецепт пасты с помидорами и базиликом.",
		Ingredients:   []string{"200 г пасты", "2 помидора", "Свежий базилик"},
		Steps:         []string{"Сварите пасту по инструкции.", "Нарежьте помидоры и базилик.", "Смешайте готовую пасту с помидорами и базиликом."},
		Author:        1,
		Collaborators: []int64{1, 2, 3, 4, 5},
	}

	err = app.writeJSON(c.Writer, http.StatusOK, Envelope{"recipe": recipe}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}

}

func (app *application) getRecipeList(c *gin.Context) {
	//jsonData, err := json.Marshal(jsonR)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при кодировании JSON"})
	//	return
	//}
	//
	//c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
}

func (app *application) addRecipe(c *gin.Context) {
	var input data.Recipe

	if err := app.readJSON(c, &input); err != nil {
		app.badRequestResponse(c, err)
		return
	}

	v := validator.New()

	if data.ValidateRecipe(v, input); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": input})

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

//func (app *application) updateRecipe(c *gin.Context) {
//	recipeId, err := app.readIDParam(c)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Рецепт не найден"})
//	}
//
//}
