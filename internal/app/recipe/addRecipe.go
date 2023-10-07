package recipe

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddRecipe(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Recipe added")
}
