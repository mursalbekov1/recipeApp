package recipe

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRecipe(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Recipe by id")
}
