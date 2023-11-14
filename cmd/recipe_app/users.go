package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_recipe/internal/data"
	"go_recipe/internal/validator"
	"net/http"
)

func (app *application) registerUserHandler(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"gmail"`
		Password string `json:"password"`
	}

	err := app.readJSON(c, &input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}
	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(c, v.Errors)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	err = app.writeJSON(c.Writer, http.StatusCreated, Envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(c, err)
	}
}
