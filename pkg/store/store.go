package store

import (
	"bookapi/pkg/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var Books []models.Book = []models.Book{
	{
		ID:          uuid.New().String(),
		Name:        "The Catcher in the Rye",
		Author:      "J.D. Salinger",
		Category:    "Fiction",
		Description: "A story about a young man’s journey of self-discovery.",
	},
	{
		ID:          uuid.New().String(),
		Name:        "To Kill a Mockingbird",
		Author:      "Harper Lee",
		Category:    "Fiction",
		Description: "A novel about the moral growth of a young girl in the Deep South.",
	},
	{
		ID:          uuid.New().String(),
		Name:        "1984",
		Author:      "George Orwell",
		Category:    "Dystopian",
		Description: "A dystopian novel about a totalitarian regime and surveillance state.",
	},
}

func CreateNewBook(book models.Book) {
	book.ID = uuid.New().String()
	Books = append(Books, book)
}

func DeleteBookById(id string) bool {
	for i, book := range Books {
		if strings.EqualFold(book.ID, id) {
			Books = append(Books[:i], Books[i+1:]...)
			return true
		}
	}
	return false
}

func GetBookById(id string) *models.Book {
	for i := range Books {
		if Books[i].ID == id {
			return &Books[i]
		}
	}
	return nil
}

func UpdateBookById(id string, updatedBook models.Book) bool {
	for i := range Books {
		if Books[i].ID == id {
			updatedBook.ID = id
			Books[i] = updatedBook
			return true
		}
	}
	return false
}

// GetAllBooks - Returns all books
func GetAllBooks() *[]models.Book {
	return &Books
}

func GetBookDetailsByISBN(isbn string) bool {
	//fmt.Println("Isbn : " + isbn)
	url := "https://openlibrary.org/api/books?bibkeys=ISBN:" + isbn + "&format=json&jscmd=data"
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Response...............................")
	// fmt.Println(res)

	defer res.Body.Close()

	var bookData map[string]map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&bookData)
	if err != nil {
		log.Fatal("Error decoding the JSON:", err)
		return false
	}

	//fmt.Println("Book Data...........................")
	//fmt.Println(bookData)

	bookInfo, ok := bookData["ISBN:"+isbn]
	if !ok {
		log.Println("No data found for ISBN::::", isbn)
		return false
	}
	//fmt.Println(bookInfo)
	authorInfo, ok := bookInfo["authors"]
	if !ok {
		log.Println("No data found for BookInfo:", isbn)
		return false
	}
	fmt.Println(authorInfo)
	return true
}
