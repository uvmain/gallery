package main

import (
	"embed"
	"gallery/core/auth"
	"gallery/core/handlers"
	"gallery/core/logic"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/go-swiss/compress"
	"github.com/rs/cors"
)

//go:embed all:frontend/dist
var dist embed.FS
var distSubFS fs.FS
var err error

func StartServer() {
	router := http.NewServeMux()

	distSubFS, err = fs.Sub(dist, "frontend/dist")
	if err != nil {
		log.Fatal("Failed to create sub filesystem:", err)
	}

	router.HandleFunc("/", HandleFrontend)

	//auth
	router.HandleFunc("POST /api/login", auth.LoginHandler)
	router.HandleFunc("GET /api/logout", auth.LogoutHandler)
	router.HandleFunc("GET /api/check-session", auth.CheckSessionHandler)

	// standard routes
	router.HandleFunc("GET /api/slugs", handlers.HandleGetSlugs)
	router.HandleFunc("GET /api/slugs/random", handlers.HandleGetRandomSlugs)
	router.HandleFunc("GET /api/slugs/with-dimensions", handlers.HandleGetSlugsWithDimensions)
	router.HandleFunc("GET /api/metadata/{slug}", handlers.HandleGetMetadataBySlug)
	router.HandleFunc("GET /api/thumbnail/{slug}", handlers.HandleGetThumbnailBySlug)
	router.HandleFunc("GET /api/optimised/{slug}", handlers.HandleGetOptimisedBySlug)
	router.HandleFunc("GET /api/original/{slug}", handlers.HandleGetOriginalImageBlobBySlug)
	router.HandleFunc("GET /api/albums/{albumSlug}", handlers.HandleGetAlbum)
	router.HandleFunc("GET /api/albums", handlers.HandleGetAllAlbums)
	router.HandleFunc("GET /api/links/album/{albumSlug}", handlers.HandleGetAlbumLinks)
	router.HandleFunc("GET /api/links/image/{imageSlug}", handlers.HandleGetImageLinks)
	router.HandleFunc("GET /api/tags", handlers.HandleGetTags)
	router.HandleFunc("GET /api/tags/{imageSlug}", handlers.HandleGetTagsBySlug)
	router.HandleFunc("GET /api/slugs/tag/{tag}", handlers.HandleGetSlugsByTag)
	router.HandleFunc("GET /api/dimensions/{imageSlug}", handlers.HandleGetDimensionsBySlug)

	// authenticated routes
	router.Handle("DELETE /api/slugs/{slug}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteImageBySlug)))
	router.Handle("PATCH /api/metadata/{slug}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePatchMetadataBySlug)))
	router.Handle("PATCH /api/albums/cover", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePatchAlbumCover)))
	router.Handle("PATCH /api/albums/name", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePatchAlbumName)))
	router.Handle("POST /api/albums", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePostAlbumRow)))
	router.Handle("DELETE /api/albums/{albumSlug}", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteAlbumRow)))
	router.Handle("POST /api/link", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePostLinkRow)))
	router.Handle("DELETE /api/link", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteAlbumLinkRow)))
	router.Handle("POST /api/links", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePostLinkRows)))
	router.Handle("POST /api/upload", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePostNewImage)))
	router.Handle("POST /api/tags", auth.AuthMiddleware(http.HandlerFunc(handlers.HandlePostNewTags)))
	router.Handle("DELETE /api/tags", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleDeleteTagRow)))

	handler := cors.AllowAll().Handler(
		compress.Middleware(router),
	)

	var serverAddress = ":8080"
	log.Println("Application running at http://localhost:8080")
	http.ListenAndServe(serverAddress, handler)
}

func HandleFrontend(w http.ResponseWriter, r *http.Request) {
	bootTime := logic.GetBootTime().Truncate(time.Second).UTC()

	cleanPath := path.Clean(r.URL.Path)
	if cleanPath == "/" {
		cleanPath = "/index.html"
	} else {
		cleanPath = strings.TrimPrefix(cleanPath, "/")
	}

	// static content
	file, err := distSubFS.Open(cleanPath)
	if err == nil {
		defer file.Close()

		http.ServeContent(w, r, cleanPath, bootTime, file.(io.ReadSeeker))
		return
	}

	// serve index.html for vue-router content
	indexFile, err := distSubFS.Open("index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusNotFound)
		return
	}
	defer indexFile.Close()

	http.ServeContent(w, r, "index.html", bootTime, indexFile.(io.ReadSeeker))
}
