package models

type BookDto struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Category    string `json:"category"`
	Description string `json:"description"`
}
