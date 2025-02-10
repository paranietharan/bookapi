package store

import (
	"bookapi/pkg/config"
	"bookapi/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

var ctx = context.Background()

func CreateNewBook(book models.Book) {
	bookJSON, err := json.Marshal(book)
	if err != nil {
		fmt.Println(err)
	}

	err = config.RedisClient.Set(ctx, book.ID, bookJSON, 0).Err()
	if err != nil {
		log.Printf("Failed to save book to Redis: %v", err)
	}

	fmt.Printf("Book saved ID: %s", book.ID)
}

func DeleteBookById(id string) bool {
	_, err := config.RedisClient.Del(ctx, id).Result()

	if err != nil {
		fmt.Println("Delete error")
		fmt.Println(err)
		return false
	}

	return true
}

func GetBookById(id string) *models.Book {
	data, err := config.RedisClient.Get(ctx, id).Bytes()
	if err != nil {
		log.Printf("Book not found in Redis: %v", err)
		return nil
	}

	var book models.Book
	err = json.Unmarshal(data, &book)
	if err != nil {
		return nil
	}
	return &book
}

func UpdateBookById(id string, updatedBook models.Book) bool {
	book := GetBookById(id)

	if book == nil {
		fmt.Printf("book not found")
	}

	CreateNewBook(updatedBook)
	return true
}

func GetAllBooks() []models.Book {
	ids, err := config.RedisClient.SMembers(ctx, "books:all").Result()
	if err != nil {
		log.Println("Failed to retrieve book IDs from Redis:", err)
		return nil
	}

	var books []models.Book
	for _, id := range ids {
		book := GetBookById(id)
		books = append(books, *book)
	}

	return books
}

func GetBookDetailsByISBN(isbn string) bool {
	url := "https://openlibrary.org/api/books?bibkeys=ISBN:" + isbn + "&format=json&jscmd=data"
	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
		return false
	}

	defer response.Body.Close()

	var bookData map[string]map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&bookData)
	if err != nil {
		log.Fatal("Error decoding the JSON:", err)
		return false
	}

	bookInfo, ok := bookData["ISBN:"+isbn]
	if !ok {
		log.Println("No data found for ISBN:", isbn)
		return false
	}

	title, ok := bookInfo["title"].(string)
	if !ok {
		log.Println("Book title not found")
		return false
	}

	authorInfo, ok := bookInfo["authors"]
	if !ok {
		log.Println("No authors found for this book")
		return false
	}

	authors, ok := authorInfo.([]interface{})
	if !ok {
		log.Println("Unexpected format for authors data")
		return false
	}

	var authorsName string
	for i, author := range authors {
		authorMap, ok := author.(map[string]interface{})
		if !ok {
			log.Println("Unexpected format for an author entry")
			continue
		}
		if authorName, exists := authorMap["name"].(string); exists {
			if i > 0 {
				authorsName += ", "
			}
			authorsName += authorName
		} else {
			log.Println("Author name not found")
		}
	}

	newBook := models.Book{
		ID:          uuid.New().String(),
		Name:        title,
		Author:      authorsName,
		Category:    "",
		Description: "",
	}

	Books = append(Books, newBook)

	fmt.Println("Book added successfully:", newBook)
	return true
}
