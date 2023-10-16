package data

import (
	"time"
)

type Recipe struct {
	ID            int64     `json:"id"`
	Time          time.Time `json:"time"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Runtime       Runtime   `json:"runtime,omitempty"`
	Ingredients   []string  `json:"ingredients"`
	Steps         []string  `json:"steps"`
	Author        int64     `json:"author"`
	Collaborators int64     `json:"collaborators"`
}

type Author struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Password       string  `json:"-"`
	Recipes        []int64 `json:"recipes"`
	RecipeAccesses []int64 `json:"access_recipes"`
}
