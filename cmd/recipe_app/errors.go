package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) logError(err error) {
	app.logger.Println(err)
}

func (app *application) errorResponse(c *gin.Context, status int, message interface{}) {
	env := gin.H{"error": message}
	c.JSON(status, env)
}

func (app *application) serverErrorResponse(c *gin.Context, err error) {
	app.logError(err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(c, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(c *gin.Context) {
	message := "the requested resource could not be found"
	app.errorResponse(c, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(c *gin.Context) {
	message := fmt.Sprintf("the %s method is not supported for this resource", c.Request.Method)
	app.errorResponse(c, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(c *gin.Context, err error) {
	app.errorResponse(c, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(c *gin.Context, errors map[string]string) {
	app.errorResponse(c, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(c *gin.Context) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(c, http.StatusOK, message)
}
