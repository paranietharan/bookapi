package store

import (
	"bookapi/pkg/config"
	"bookapi/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var books []models.Book

var ctx = context.Background()

func CreateNewBook(book models.Book) {
	book.ID = uuid.New().String()
	bookJSON, err := json.Marshal(book)
	if err != nil {
		log.Println("Error marshalling book:", err)
		return
	}

	err = config.RedisClient.Set(ctx, book.ID, bookJSON, 0).Err()
	if err != nil {
		log.Printf("Failed to save book to Redis: %v", err)
		return
	}

	err = config.RedisClient.SAdd(ctx, "books:all", book.ID).Err()
	if err != nil {
		log.Printf("Failed to add book ID to Redis set: %v", err)
	}

	books = append(books, book) // slice
	book.LastAccessTime = time.Now()

	fmt.Println("Book saved with ID:", book.ID)
}

func DeleteBookById(id string) bool {
	_, err := config.RedisClient.Del(ctx, id).Result()
	if err != nil {
		log.Println("Error deleting book from Redis:", err)
		return false
	}

	_, err = config.RedisClient.SRem(ctx, "books:all", id).Result()
	if err != nil {
		log.Println("Error removing book ID from Redis set:", err)
		return false
	}

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			break
		}
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
		log.Println("Error unmarshalling book data:", err)
		return nil
	}
	book.LastAccessTime = time.Now()
	return &book
}

func UpdateBookById(id string, updatedBook models.Book) bool {
	existingBook := GetBookById(id)
	if existingBook == nil {
		log.Println("Book not found, update failed")
		return false
	}

	updatedBook.ID = id

	bookJSON, err := json.Marshal(updatedBook)
	if err != nil {
		log.Println("Error marshalling updated book:", err)
		return false
	}

	err = config.RedisClient.Set(ctx, id, bookJSON, 0).Err()
	if err != nil {
		log.Println("Failed to update book in Redis:", err)
		return false
	}
	log.Println("Book updated successfully in Redis")

	for i, book := range books {
		if book.ID == id {
			books[i] = updatedBook
			break
		}
	}
	updatedBook.LastAccessTime = time.Now()
	log.Println("Book updated successfully in slice")

	return true
}

func GetAllBooks() []models.BookDto {
	ids, err := config.RedisClient.SMembers(ctx, "books:all").Result()
	if err != nil {
		log.Println("Failed to retrieve book IDs from Redis:", err)
		return nil
	}

	var books []models.BookDto
	for _, id := range ids {
		book := GetBookById(id)
		if book != nil {
			books = append(books, book.ToBookDto())
			book.LastAccessTime = time.Now()
		}
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

	books = append(books, newBook)
	CreateNewBook(newBook)
	newBook.LastAccessTime = time.Now()

	fmt.Println("Book added successfully:", newBook)
	return true
}
