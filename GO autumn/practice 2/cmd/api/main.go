package main

import (
	"log"
	"net/http"
	"practice_2/internal/handlers"
	"practice_2/internal/middleware"
)

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/user", handlers.UserHandler)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", middleware.AuthMiddleware(mux))
}