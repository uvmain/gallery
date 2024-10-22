package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rs/cors"
)

func StartServer() {
	router := http.NewServeMux()
	router.HandleFunc("GET /api/slugs", handleGetSlugs)
	router.HandleFunc("GET /api/metadata/{slug}", handleGetMetadataBySlug)
	router.HandleFunc("GET /api/thumbnail/{slug}", handleGetThumbnailBySlug)
	router.HandleFunc("GET /api/optimised/{slug}", handleGetOptimisedBySlug)
	router.HandleFunc("GET /api/original/{slug}", handleGetOriginalImageBySlug)

	handler := cors.AllowAll().Handler(router)

	http.ListenAndServe("localhost:8080", handler)
}

func handleGetSlugs(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("offset")
	offset, _ := strconv.Atoi(v)
	v = r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(v)
	slugs, _ := GetSlugsOrderedByDateTaken(offset, limit)
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

func handleGetThumbnailBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	thumbnail, err := GetThumbnailBySlug(slug)
	if err != nil {
		http.Error(w, "Thumbnail not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/webp")
	w.WriteHeader(http.StatusOK)
	w.Write(thumbnail)
}

func handleGetOptimisedBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	optimised, err := GetOptimisedBySlug(slug)
	if err != nil {
		http.Error(w, "Optimised not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/webp")
	w.WriteHeader(http.StatusOK)
	w.Write(optimised)
}

func handleGetOriginalImageBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	imageBlob, err := GetOriginalImageBySlug(slug)

	if err != nil {
		http.Error(w, "Optimised not found", http.StatusNotFound)
		return
	}
	mimeType := http.DetectContentType(imageBlob)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(imageBlob)
}
