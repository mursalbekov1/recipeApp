package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (app *application) readIDParam(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("недопустимый параметр id")
	}
	return id, nil
}

type Envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
