package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func StartServer() {
	router := http.NewServeMux()
	router.HandleFunc("GET /api/slugs", handleGetSlugs)
	router.HandleFunc("GET /api/metadata/{slug}", handleGetMetadataBySlug)
	router.HandleFunc("GET /api/thumbnail/{slug}", handleGetThumbnailBySlug)
	router.HandleFunc("GET /api/optimised/{slug}", handleGetOptimisedBySlug)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server listening on port :8080")
	server.ListenAndServe()
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
	if ImageFormat == "webp" {
		w.Header().Set("Content-Type", "image/webp")
	} else if ImageFormat == "jpeg" || ImageFormat == "jpg" {
		w.Header().Set("Content-Type", "image/jpeg")
	} else {
		http.Error(w, "Unsupported image format", http.StatusInternalServerError)
		return
	}
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
	if ImageFormat == "webp" {
		w.Header().Set("Content-Type", "image/webp")
	} else if ImageFormat == "jpeg" || ImageFormat == "jpg" {
		w.Header().Set("Content-Type", "image/jpeg")
	} else {
		http.Error(w, "Unsupported image format", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(optimised)
}
