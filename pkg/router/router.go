package router

import (
	"bookapi/pkg/config"
	handlers "bookapi/pkg/handler"
	"bookapi/pkg/logger"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitializeRoutes() *mux.Router {
	config.InitRedis()

	router := mux.NewRouter()

	router.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", handlers.GetBookByID).Methods("GET")
	router.HandleFunc("/books", handlers.CreateBook).Methods("POST")
	router.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")

	// endoint that add books when user passes only the isbn number
	router.HandleFunc("/book/isbn/{id}", handlers.GetBookDetailsByISBN).Methods("GET")

	return router
}

func StartServer() {
	logger.StartLogListener()

	router := InitializeRoutes()
	http.Handle("/", router)
	fmt.Println("Book api server started.......")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
