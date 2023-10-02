package recipe

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateRecipe(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Recipe added")
}
