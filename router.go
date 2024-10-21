package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func StartServer() {
	router := http.NewServeMux()
	router.HandleFunc("GET /api/slugs", handleGetSlugs)
	router.HandleFunc("GET /api/metadata/{slug}", handleGetMetadataBySlug)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server listening on port :8080")
	server.ListenAndServe()
}

func handleGetSlugs(w http.ResponseWriter, r *http.Request) {
	slugs, _ := GetSlugsOrderedByDateTaken()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slugs); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetMetadataBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	metadata, _ := GetMetadataBySlug(slug)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metadata); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
