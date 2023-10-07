package middleWare

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	"strconv"
)

func ReadIDParam(r *gin.Context) (int64, error) {
	params := httprouter.ParamsFromContext(r)
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}
