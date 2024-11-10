package main

import (
	"encoding/json"
	"net/http"
	"photogallery/auth"
	"photogallery/database"
	"photogallery/logic"
	"photogallery/types"
	"strconv"

	"github.com/rs/cors"
)

func StartServer() {
	router := http.NewServeMux()

	// frontend
	distDir := http.Dir("../dist")
	fileServer := http.FileServer(distDir)
	router.Handle("/{path...}", http.StripPrefix("/", fileServer))

	//auth
	router.HandleFunc("POST /api/login", auth.LoginHandler)
	router.HandleFunc("GET /api/logout", auth.LogoutHandler)
	router.HandleFunc("GET /api/check-session", auth.CheckSessionHandler)

	// standard routes
	router.HandleFunc("GET /api/slugs", handleGetSlugs)
	router.HandleFunc("GET /api/slugs/random", handleGetRandomSlugs)
	router.HandleFunc("GET /api/metadata/{slug}", handleGetMetadataBySlug)
	router.HandleFunc("GET /api/thumbnail/{slug}", handleGetThumbnailBySlug)
	router.HandleFunc("GET /api/optimised/{slug}", handleGetOptimisedBySlug)
	router.HandleFunc("GET /api/original/{slug}", handleGetOriginalImageBlobBySlug)
	router.HandleFunc("GET /api/albums/{albumSlug}", handleGetAlbum)
	router.HandleFunc("GET /api/albums", handleGetAllAlbums)
	router.HandleFunc("GET /api/links/album/{albumSlug}", handleGetAlbumLinks)
	router.HandleFunc("GET /api/links/image/{imageSlug}", handleGetImageLinks)

	// authenticated routes
	router.Handle("PATCH /api/metadata/{slug}", auth.AuthMiddleware(http.HandlerFunc(handlePatchMetadataBySlug)))
	router.Handle("POST /api/albums", auth.AuthMiddleware(http.HandlerFunc(handlePostAlbumRow)))
	router.Handle("DELETE /api/albums/{albumSlug}", auth.AuthMiddleware(http.HandlerFunc(handleDeleteAlbumRow)))

	handler := cors.AllowAll().Handler(router)

	var serverAddress string
	if logic.IsLocalDevEnv() {
		serverAddress = "localhost:8080"
	} else {
		serverAddress = ":8080"
	}

	http.ListenAndServe(serverAddress, handler)
}

func handleGetSlugs(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("offset")
	if v == "" {
		v = "0"
	}
	offset, _ := strconv.Atoi(v)
	v = r.URL.Query().Get("limit")
	if v == "" {
		v = "1000"
	}
	limit, _ := strconv.Atoi(v)
	slugs, _ := database.GetSlugsOrderedByDateTaken(offset, limit)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slugs); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetRandomSlugs(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("limit")
	if v == "" {
		v = "1000"
	}
	limit, _ := strconv.Atoi(v)
	slugs, _ := database.GetSlugsOrderedRandom(limit)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slugs); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetMetadataBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	metadata, _ := database.GetMetadataBySlug(slug)
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
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(thumbnail)
}

func handleGetAlbum(w http.ResponseWriter, r *http.Request) {
	albumSlug := r.PathValue("albumSlug")
	album, err := database.GetAlbum(albumSlug)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(album); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetAllAlbums(w http.ResponseWriter, r *http.Request) {
	albums := database.GetAllAlbums()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetOptimisedBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	optimised, err := GetOptimisedBySlug(slug)
	if err != nil {
		http.Error(w, "Optimised not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(optimised)
}

func handleGetOriginalImageBlobBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	imageBlob, err := database.GetOriginalImageBlobBySlug(slug)

	if err != nil {
		http.Error(w, "Original image not found", http.StatusNotFound)
		return
	}
	mimeType := http.DetectContentType(imageBlob)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(imageBlob)
}

func handlePatchMetadataBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := database.UpdateMetadataBySlug(slug, updates); err != nil {
		http.Error(w, "Failed to update metadata", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metadata updated successfully"))
}

func handlePostAlbumRow(w http.ResponseWriter, r *http.Request) {
	var updates types.Album
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := database.InsertAlbumRow(updates); err != nil {
		http.Error(w, "Failed to update metadata", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metadata updated successfully"))
}

func handleDeleteAlbumRow(w http.ResponseWriter, r *http.Request) {
	albumSlug := r.PathValue("albumSlug")
	if err := database.DeleteAlbumRow(albumSlug); err != nil {
		http.Error(w, "Failed to delete album", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Album deleted successfully"))
}

func handleGetAlbumLinks(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("albumSlug")
	links, err := database.GetAlbumLinks(slug)

	if err != nil {
		http.Error(w, "Failed to retrieve album links", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(links); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetImageLinks(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("imageSlug")
	links, err := database.GetImageLinks(slug)

	if err != nil {
		http.Error(w, "Failed to retrieve image links", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(links); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
