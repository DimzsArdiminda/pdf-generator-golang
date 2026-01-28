package routes

import (
	"github.com/gorilla/mux"
	"pdf-api/handlers"
)

func RegisterRoutes(r *mux.Router) {
	api := r.PathPrefix("/api/pdf").Subrouter()

	api.HandleFunc("/generate", handlers.GeneratePDF).Methods("POST")
	api.HandleFunc("/upload", handlers.UploadPDF).Methods("POST")
	api.HandleFunc("/list", handlers.ListPDF).Methods("GET")
	api.HandleFunc("/{id}", handlers.DeletePDF).Methods("DELETE")
}
