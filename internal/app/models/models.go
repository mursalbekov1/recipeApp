package models

import "time"

type Recipe struct {
	ID          int64     `json:"ID"`
	Time        time.Time `json:"Time"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Ingredients []string  `json:"Ingredients"`
	Steps       []string  `json:"Steps"`
	Author      int64     `json:"Author"`
}
