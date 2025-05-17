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

	// Store book in Redis
	err = config.RedisClient.Set(ctx, book.ID, bookJSON, 0).Err()
	if err != nil {
		log.Printf("Failed to save book to Redis: %v", err)
		return
	}

	err = config.RedisClient.Set(ctx, "timestamp:"+book.ID, book.LastAccessTime.Unix(), 0).Err()
	if err != nil {
		log.Printf("Failed to store book timestamp: %v", err)
	}

	err = config.RedisClient.SAdd(ctx, "books:all", book.ID).Err()
	if err != nil {
		log.Printf("Failed to add book ID to Redis set: %v", err)
	}

	// Add book to slice
	books = append(books, book)

	fmt.Println("Book saved with ID:", book.ID)
}

// Delete book from Redis only
func DeleteBookByIdRedis(id string) bool {
	_, err := config.RedisClient.Del(ctx, id, "timestamp:"+id).Result()
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

func DeleteBookById(id string) bool {
	// delete in redis
	del := DeleteBookByIdRedis(id)
	if !del {
		fmt.Println("Deleting failed: some time it may not in redis")
	}

	// delete from the slice
	for i, v := range books {
		if v.ID == id {
			books = append(books[:i], books[i+1:]...)
			fmt.Println("Book deleted successfully")
		}
	}

	return true
}

// Retrieve book from Redis & update last access time
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

	config.RedisClient.Set(ctx, "timestamp:"+book.ID, book.LastAccessTime.Unix(), 0)

	return &book
}

// Update book details
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

	config.RedisClient.Set(ctx, "timestamp:"+id, updatedBook.LastAccessTime.Unix(), 0)

	for i, book := range books {
		if book.ID == id {
			books[i] = updatedBook
			break
		}
	}

	log.Println("Book updated successfully")
	return true
}

func GetAllBooks() []models.BookDto {
	var bookDtos []models.BookDto
	for _, v := range books {
		bookDtos = append(bookDtos, v.ToBookDto())

		// change the last access time in redis
		config.RedisClient.Set(ctx, "timestamp:"+v.ID, v.LastAccessTime.Unix(), 0)
	}

	return bookDtos
}

// Fetch book details using ISBN
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
		log.Fatal("Error decoding JSON:", err)
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
			continue
		}
		if authorName, exists := authorMap["name"].(string); exists {
			if i > 0 {
				authorsName += ", "
			}
			authorsName += authorName
		}
	}

	newBook := models.Book{
		ID:             uuid.New().String(),
		Name:           title,
		Author:         authorsName,
		Category:       "",
		Description:    "",
		LastAccessTime: time.Now(),
	}

	books = append(books, newBook)
	CreateNewBook(newBook)

	fmt.Println("Book added successfully:", newBook)
	return true
}

// Cleanup Goroutine - Remove books not accessed for 10 minutes
func StartBookCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			cleanupOldBooks()
			fmt.Println("Celen Up completed")
		}
	}()
}

// Function to delete books not accessed in 10 mins
func cleanupOldBooks() {
	bookIDs, err := config.RedisClient.SMembers(ctx, "books:all").Result()
	if err != nil {
		log.Println("Failed to retrieve book IDs:", err)
		return
	}

	currentTime := time.Now().Unix()
	for _, bookID := range bookIDs {
		timestamp, err := config.RedisClient.Get(ctx, "timestamp:"+bookID).Int64()
		if err != nil {
			log.Printf("Error retrieving timestamp for book %s: %v", bookID, err)
			continue
		}

		// del when 20s or more access time items
		if currentTime-timestamp > 20 {
			success := DeleteBookByIdRedis(bookID)
			if success {
				log.Printf("Removed inactive book from Redis: %s\n", bookID)
			}
		}
	}
}
