package models

type BookDto struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

func (b Book) ToBookDto() BookDto {
	return BookDto{
		Name:        b.Name,
		Author:      b.Author,
		Category:    b.Category,
		Description: b.Description,
	}
}
