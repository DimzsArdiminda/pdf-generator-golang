package main

import (
	"log"
	"net/http"
	"os"

	"pdf-api/config"
	"pdf-api/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main(){
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	config.ConnectDB()

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running at :%s", port)
	http.ListenAndServe(":"+port, r)
}