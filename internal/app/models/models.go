package models

import "time"

type Recipe struct {
	ID          int64     `json:"id"`
	Time        time.Time `json:"Time"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Ingredients []string  `json:"Ingredients"`
	Steps       []string  `json:"Steps"`
	Author      int64     `json:"Author"`
}

type Author struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Recipes  []int64 `json:"recipes"`
}
