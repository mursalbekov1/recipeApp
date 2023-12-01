package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_recipe/internal/data"
	"go_recipe/internal/validator"
	"net/http"
	"time"
)

func (app *application) createAuthenticationTokenHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(c, &input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(c)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	err = app.writeJSON(c.Writer, http.StatusCreated, Envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}
}
