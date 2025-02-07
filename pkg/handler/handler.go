package handlers

import (
	"bookapi/pkg/models"
	"bookapi/pkg/store"
	"bookapi/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	books := store.GetAllBooks()
	utils.WriteResponse(w, http.StatusOK, books)
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	book := store.GetBookById(id)
	if book == nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	utils.WriteResponse(w, http.StatusOK, book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var bookDto models.BookDto

	err := json.NewDecoder(r.Body).Decode(&bookDto)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}

	book := models.Book{
		ID:          uuid.New().String(),
		Name:        bookDto.Name,
		Author:      bookDto.Author,
		Category:    bookDto.Category,
		Description: bookDto.Description,
	}

	store.CreateNewBook(book)
	utils.WriteResponse(w, http.StatusCreated, book)

	fmt.Println("Book added : " + book.ID + " | Book name : " + book.Name)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var bookDto models.BookDto
	err := json.NewDecoder(r.Body).Decode(&bookDto)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	existingBook := store.GetBookById(id)
	if existingBook == nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	updatedBook := models.Book{
		ID:          id,
		Name:        bookDto.Name,
		Author:      bookDto.Author,
		Category:    bookDto.Category,
		Description: bookDto.Description,
	}

	success := store.UpdateBookById(id, updatedBook)
	if !success {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	utils.WriteResponse(w, http.StatusOK, updatedBook)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	success := store.DeleteBookById(id)
	if !success {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	utils.WriteResponse(w, http.StatusNoContent, nil)
}

// functions for get from external api
func GetBookDetailsByISBN(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isbn := vars["id"]

	success := store.GetBookDetailsByISBN(isbn)
	if !success {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	utils.WriteResponse(w, http.StatusNoContent, nil)
}
