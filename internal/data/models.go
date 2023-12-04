package data

import (
	"database/sql"
	"errors"
	"time"
)

type Recipe struct {
	ID            int64     `json:"id"`
	Time          time.Time `json:"time"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Ingredients   []string  `json:"ingredients"`
	Steps         []string  `json:"steps"`
	Author        int64     `json:"author"`
	Collaborators []int64   `json:"collaborators"`
	Version       int64     `json:"version"`
}

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Recipe interface {
		Insert(recipe *Recipe) error
		Get(id int64) (*Recipe, error)
		Update(recipe *Recipe) error
		Delete(id int64) error
		GetAll(title string, ingredients []string, filters Filters) ([]*Recipe, Metadata, error)
	}

	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Recipe:      RecipeModel{DB: db},
		Permissions: PermissionModel{DB: db},
	}
}
