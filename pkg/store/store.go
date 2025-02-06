package store

import "bookapi/pkg/models"

var Books []models.Book = []models.Book{
	{
		ID:          1,
		Name:        "The Catcher in the Rye",
		Author:      "J.D. Salinger",
		Category:    "Fiction",
		Description: "A story about a young man’s journey of self-discovery.",
	},
	{
		ID:          2,
		Name:        "To Kill a Mockingbird",
		Author:      "Harper Lee",
		Category:    "Fiction",
		Description: "A novel about the moral growth of a young girl in the Deep South.",
	},
	{
		ID:          3,
		Name:        "1984",
		Author:      "George Orwell",
		Category:    "Dystopian",
		Description: "A dystopian novel about a totalitarian regime and surveillance state.",
	},
}

func CreateNewBook(book models.Book) {
	Books = append(Books, book)
}

func DeleteBookById(id int) bool {
	for i, book := range Books {
		if book.ID == id {
			Books = append(Books[:i], Books[i+1:]...)
			return true
		}
	}

	return false
}

func GetBookById(id int) *models.Book {
	for _, book := range Books {
		if book.ID == id {
			return &book
		}
	}
	return nil
}

func UpdateBookById(id int, updatedBook models.Book) bool {
	for _, book := range Books {
		if book.ID == id {
			Books[id] = updatedBook
			return true
		}
	}
	return false
}

func GetAllBooks() *[]models.Book {
	return &Books
}
