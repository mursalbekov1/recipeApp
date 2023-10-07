package recipe

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRecipesList(c *gin.Context) {

	jsonData, err := json.Marshal(jsonR)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при кодировании JSON"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)

}
