package recipe

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRecipesList(ctx *gin.Context) {

	ctx.String(http.StatusOK, "List of recipes")
}
