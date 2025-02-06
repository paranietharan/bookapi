package handlers

import (
	"bookapi/pkg/models"
	"bookapi/pkg/store"
	"bookapi/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	books := store.GetAllBooks()
	utils.WriteResponse(w, http.StatusOK, books)
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	book := store.GetBookById(id)
	if book == nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	utils.WriteResponse(w, http.StatusOK, book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	store.CreateNewBook(book)
	utils.WriteResponse(w, http.StatusCreated, book)

	fmt.Println("Book added : " + strconv.Itoa(book.ID) + "Book name : " + book.Name)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var updatedBook models.Book
	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedBook.ID = id
	success := store.UpdateBookById(id, updatedBook)
	if !success {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	utils.WriteResponse(w, http.StatusOK, updatedBook)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	success := store.DeleteBookById(id)
	if !success {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	utils.WriteResponse(w, http.StatusNoContent, nil)
}
