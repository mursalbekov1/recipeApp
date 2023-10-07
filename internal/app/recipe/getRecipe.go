package recipe

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"recipeApp/internal/app/JSON"
	"recipeApp/internal/app/models"
	"time"
)

func GetRecipe(c *gin.Context) {
	recipe := models.Recipe{
		ID:          1,
		Time:        time.Now(),
		Title:       "",
		Description: "",
		Ingredients: []string{"200 г пасты", "2 помидора", "Свежий базилик"},
		Steps:       []string{"Сварите пасту по инструкции.", "Нарежьте помидоры и базилик.", "Смешайте готовую пасту с помидорами и базиликом."},
		Author:      1,
	}

	err := JSON.WriteJSON(c.Writer, http.StatusOK, JSON.Envelope{"recipe": recipe}, nil)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "The server encountered a problem and could not process your request"})
	}

}