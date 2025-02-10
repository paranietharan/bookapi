package store

import (
	"bookapi/pkg/config"
	"bookapi/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

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
		}
	}

	return books
}
