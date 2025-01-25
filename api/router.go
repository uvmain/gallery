package main

import (
	"encoding/json"
	"log"
	"net/http"
	"photogallery/auth"
	"photogallery/database"
	"photogallery/image"
	"photogallery/logic"
	"photogallery/optimised"
	"photogallery/thumbnails"
	"photogallery/types"
	"time"

	"github.com/rs/cors"
)

func enableCdnCaching(w http.ResponseWriter) {
	expiryDate := time.Now().AddDate(1, 0, 0)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("Expires", expiryDate.String())
}

func StartServer() {
	router := http.NewServeMux()

	// frontend
	distDir := http.Dir("../dist")
	fileServer := http.FileServer(distDir)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// serve static files
		if _, err := distDir.Open(r.URL.Path); err == nil {
			enableCdnCaching(w)
			fileServer.ServeHTTP(w, r)
			return
		}
		// serve index.html for non-static files
		http.ServeFile(w, r, "../dist/index.html")
	})

	//auth
	router.HandleFunc("POST /api/login", auth.LoginHandler)
	router.HandleFunc("GET /api/logout", auth.LogoutHandler)
	router.HandleFunc("GET /api/check-session", auth.CheckSessionHandler)

	// standard routes
	router.HandleFunc("GET /api/slugs", handleGetSlugs)
	router.HandleFunc("GET /api/slugs/random", handleGetRandomSlugs)
	router.HandleFunc("GET /api/slugs/with-dimensions", handleGetSlugsWithDimensions)
	router.HandleFunc("GET /api/metadata/{slug}", handleGetMetadataBySlug)
	router.HandleFunc("GET /api/thumbnail/{slug}", handleGetThumbnailBySlug)
	router.HandleFunc("GET /api/optimised/{slug}", handleGetOptimisedBySlug)
	router.HandleFunc("GET /api/original/{slug}", handleGetOriginalImageBlobBySlug)
	router.HandleFunc("GET /api/albums/{albumSlug}", handleGetAlbum)
	router.HandleFunc("GET /api/albums", handleGetAllAlbums)
	router.HandleFunc("GET /api/links/album/{albumSlug}", handleGetAlbumLinks)
	router.HandleFunc("GET /api/links/image/{imageSlug}", handleGetImageLinks)
	router.HandleFunc("GET /api/tags", handleGetTags)
	router.HandleFunc("GET /api/tags/{imageSlug}", handleGetTagsBySlug)
	router.HandleFunc("GET /api/slugs/tag/{tag}", handleGetSlugsByTag)
	router.HandleFunc("GET /api/dimensions/{imageSlug}", handleGetDimensionsBySlug)

	// authenticated routes
	router.Handle("PATCH /api/metadata/{slug}", auth.AuthMiddleware(http.HandlerFunc(handlePatchMetadataBySlug)))
	router.Handle("PATCH /api/albums/cover", auth.AuthMiddleware(http.HandlerFunc(handlePatchAlbumCover)))
	router.Handle("PATCH /api/albums/name", auth.AuthMiddleware(http.HandlerFunc(handlePatchAlbumName)))
	router.Handle("POST /api/albums", auth.AuthMiddleware(http.HandlerFunc(handlePostAlbumRow)))
	router.Handle("DELETE /api/albums/{albumSlug}", auth.AuthMiddleware(http.HandlerFunc(handleDeleteAlbumRow)))
	router.Handle("POST /api/link", auth.AuthMiddleware(http.HandlerFunc(handlePostLinkRow)))
	router.Handle("DELETE /api/link", auth.AuthMiddleware(http.HandlerFunc(handleDeleteAlbumLinkRow)))
	router.Handle("POST /api/links", auth.AuthMiddleware(http.HandlerFunc(handlePostLinkRows)))
	router.Handle("POST /api/upload", auth.AuthMiddleware(http.HandlerFunc(handlePostNewImage)))
	router.Handle("POST /api/tags", auth.AuthMiddleware(http.HandlerFunc(handlePostNewTags)))
	router.Handle("DELETE /api/tags", auth.AuthMiddleware(http.HandlerFunc(handleDeleteTagRow)))

	handler := cors.AllowAll().Handler(router)

	var serverAddress string
	if logic.IsLocalDevEnv() {
		serverAddress = "localhost:8080"
		log.Println("Application running at https://photogallery.localhost")
	} else {
		serverAddress = ":8080"
		log.Println("Application running at http://localhost:8080")
	}

	http.ListenAndServe(serverAddress, handler)
}

func handleGetSlugs(w http.ResponseWriter, r *http.Request) {
	slugs, _ := database.GetSlugsOrderedByDateTaken()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slugs); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetSlugsWithDimensions(w http.ResponseWriter, r *http.Request) {
	slugs, _ := database.GetSlugsWithDimensions()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slugs); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetRandomSlugs(w http.ResponseWriter, r *http.Request) {
	slugs, _ := database.GetSlugsOrderedRandom()
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
	thumbnail, err := thumbnails.GetThumbnailBySlug(slug)
	if err != nil {
		http.Error(w, "Thumbnail not found", http.StatusNotFound)
		return
	}
	enableCdnCaching(w)
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
	optimised, err := optimised.GetOptimisedBySlug(slug)
	if err != nil {
		http.Error(w, "Optimised not found", http.StatusNotFound)
		return
	}
	enableCdnCaching(w)
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
	enableCdnCaching(w)
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
		http.Error(w, "Failed to post album", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Album posted successfully"))
}

func handlePostLinkRow(w http.ResponseWriter, r *http.Request) {
	var updates types.Link
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := database.InsertAlbumLinkRow(updates); err != nil {
		http.Error(w, "Failed to insert link", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Link row inserted successfully"))
}

func handleDeleteAlbumLinkRow(w http.ResponseWriter, r *http.Request) {
	var updates types.Link
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := database.DeleteAlbumLinkRow(updates); err != nil {
		http.Error(w, "Failed to insert link", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Link row inserted successfully"))
}

func handlePostLinkRows(w http.ResponseWriter, r *http.Request) {
	var updates types.Links

	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	for _, imageSlug := range updates.ImageSlugs {
		update := types.Link{
			AlbumSlug: updates.AlbumSlug,
			ImageSlug: imageSlug,
		}
		if err := database.InsertAlbumLinkRow(update); err != nil {
			http.Error(w, "Failed to insert link: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Link rows inserted successfully"))
}

func handlePatchAlbumCover(w http.ResponseWriter, r *http.Request) {
	type CoverUpdate struct {
		AlbumSlug string
		CoverSlug string
	}
	var updates CoverUpdate
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	if err := database.UpdateAlbumCover(updates.AlbumSlug, updates.CoverSlug); err != nil {
		http.Error(w, "Failed to insert link: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cover rows updated successfully"))
}

func handlePatchAlbumName(w http.ResponseWriter, r *http.Request) {
	type AlbumNameUpdate struct {
		AlbumSlug string
		AlbumName string
	}
	var update AlbumNameUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	if err := database.UpdateAlbumName(update.AlbumSlug, update.AlbumName); err != nil {
		http.Error(w, "Failed to udpate Album Name: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Album Name updated successfully"))
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

func handlePostNewImage(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	slug := image.UploadImage(file, fileHeader)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slug); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := database.GetAllTags()

	if err != nil {
		http.Error(w, "Failed to retrieve tags", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetTagsBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("imageSlug")
	tags, err := database.GetTagsForSlug(slug)

	if err != nil {
		http.Error(w, "Failed to retrieve tags for slug", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetSlugsByTag(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	slugs, err := database.GetSlugsForTag(tag)

	w.Header().Set("Content-Type", "application/json")

	if err != nil || slugs == nil {
		slugs = []string{}
	}

	if err := json.NewEncoder(w).Encode(slugs); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handlePostNewTags(w http.ResponseWriter, r *http.Request) {
	var updates types.Tags

	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	for _, imageSlug := range updates.ImageSlugs {
		update := types.Tag{
			Tag:       updates.Tag,
			ImageSlug: imageSlug,
		}
		if err := database.InsertTagsRow(update); err != nil {
			http.Error(w, "Failed to insert tag: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("tag rows inserted successfully"))
}

func handleDeleteTagRow(w http.ResponseWriter, r *http.Request) {
	var updates types.Tag

	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := database.DeleteTagsRow(updates); err != nil {
		http.Error(w, "Failed to delete tag row", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tag deleted successfully"))
}

func handleGetDimensionsBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("imageSlug")
	dimensions, err := database.GetDimensionForSlug(slug)

	if err != nil {
		http.Error(w, "Failed to retrieve dimensions for slug", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dimensions); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
