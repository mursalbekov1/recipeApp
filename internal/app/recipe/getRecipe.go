package recipe

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"recipeApp/internal/app/JSON"
	"recipeApp/internal/app/models"
	"strconv"
	"time"
)

func GetRecipe(c *gin.Context) {
	id := c.Param("id")

	recipeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	recipe := models.Recipe{
		ID:          int64(recipeID),
		Time:        time.Now(),
		Title:       "Паста с помидорами",
		Description: "Простой рецепт пасты с помидорами и базиликом.",
		Ingredients: []string{"200 г пасты", "2 помидора", "Свежий базилик"},
		Steps:       []string{"Сварите пасту по инструкции.", "Нарежьте помидоры и базилик.", "Смешайте готовую пасту с помидорами и базиликом."},
		Author:      1,
	}

	err = JSON.WriteJSON(c.Writer, http.StatusOK, JSON.Envelope{"recipe": recipe}, nil)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "The server encountered a problem and could not process your request"})
	}

}
