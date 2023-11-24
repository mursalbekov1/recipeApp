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

type Author struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Password       string  `json:"-"`
	Recipes        []int64 `json:"recipes"`
	RecipeAccesses []int64 `json:"access_recipes"`
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

	Author interface {
		Insert(author *Author) error
		Get(id int64) (*Author, error)
		Update(author *Author) error
		Delete(id int64) error
		GetAll(name string, email string) ([]*Author, Metadata, error)
	}
	Users  UserModel
	Tokens TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db},
		Recipe: RecipeModel{DB: db},
	}
}
