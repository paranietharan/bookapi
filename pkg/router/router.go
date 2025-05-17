package router

import (
	"bookapi/pkg/config"
	handlers "bookapi/pkg/handler"
	"bookapi/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"sync"

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
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start loging function
	logger.StartLogListener(ctx, &wg)

	router := InitializeRoutes()
	http.Handle("/", router)

	// create goroutines &
	// Start the server
	server := &http.Server{Addr: ":8080"}

	fmt.Println("Server started on port 8080...")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Server stopped:", err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(cancel, &wg)
}

func waitForShutdown(cancel context.CancelFunc, wg *sync.WaitGroup) {
	fmt.Println("Press Enter to shut down...")
	fmt.Scanln()
	cancel()
	wg.Wait()
	fmt.Println("Shutdown complete.")
}
