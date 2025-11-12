package handlers

import (
	"encoding/json"
	"gallery/core/database"
	"gallery/core/image"
	"gallery/core/net"
	"gallery/core/optimised"
	"gallery/core/thumbnails"
	"gallery/core/types"
	"log"
	"net/http"
)

func HandleDeleteImageBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	log.Printf("deleting slug %s", slug)
	w.Header().Set("Content-Type", "application/json")

	err := optimised.DeleteOptimisedBySlug(slug)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	err = thumbnails.DeleteThumbnailBySlug(slug)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	filename, err := database.DeleteImageBySlug(slug)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	err = image.DeleteOriginalImage(filename)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func HandleGetSlugs(w http.ResponseWriter, r *http.Request) {
	slugs, _ := database.GetSlugsOrderedByDateTaken()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slugs); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func HandleGetSlugsWithDimensions(w http.ResponseWriter, r *http.Request) {
	slugs, _ := database.GetSlugsWithDimensions()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slugs); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func HandleGetRandomSlugs(w http.ResponseWriter, r *http.Request) {
	slugs, _ := database.GetSlugsOrderedRandom()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slugs); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func HandleGetMetadataBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	metadata, _ := database.GetMetadataBySlug(slug)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metadata); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func HandleGetThumbnailBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	thumbnail, err := thumbnails.GetThumbnailBySlug(slug)
	if err != nil {
		http.Error(w, "Thumbnail not found", http.StatusNotFound)
		return
	}
	net.EnableCdnCaching(w)
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(thumbnail)
}

func HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
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

func HandleGetAllAlbums(w http.ResponseWriter, r *http.Request) {
	albums := database.GetAllAlbums()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func HandleGetOptimisedBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	optimised, err := optimised.GetOptimisedBySlug(slug)
	if err != nil {
		http.Error(w, "Optimised not found", http.StatusNotFound)
		return
	}
	net.EnableCdnCaching(w)
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(optimised)
}

func HandleGetOriginalImageBlobBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	imageBlob, err := database.GetOriginalImageBlobBySlug(slug)

	if err != nil {
		http.Error(w, "Original image not found", http.StatusNotFound)
		return
	}
	mimeType := http.DetectContentType(imageBlob)
	net.EnableCdnCaching(w)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(imageBlob)
}

func HandlePatchMetadataBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := database.UpdateMetadataBySlug(slug, updates); err != nil {
		log.Printf("Failed to update metadata: %s", err)
		http.Error(w, "Failed to update metadata", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metadata updated successfully"))
}

func HandlePostAlbumRow(w http.ResponseWriter, r *http.Request) {
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

func HandlePostLinkRow(w http.ResponseWriter, r *http.Request) {
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

func HandleDeleteAlbumLinkRow(w http.ResponseWriter, r *http.Request) {
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

func HandlePostLinkRows(w http.ResponseWriter, r *http.Request) {
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

func HandlePatchAlbumCover(w http.ResponseWriter, r *http.Request) {
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

func HandlePatchAlbumName(w http.ResponseWriter, r *http.Request) {
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

func HandleDeleteAlbumRow(w http.ResponseWriter, r *http.Request) {
	albumSlug := r.PathValue("albumSlug")
	if err := database.DeleteAlbumRow(albumSlug); err != nil {
		http.Error(w, "Failed to delete album", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Album deleted successfully"))
}

func HandleGetAlbumLinks(w http.ResponseWriter, r *http.Request) {
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

func HandleGetImageLinks(w http.ResponseWriter, r *http.Request) {
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

func HandlePostNewImage(w http.ResponseWriter, r *http.Request) {
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

func HandleGetTags(w http.ResponseWriter, r *http.Request) {
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

func HandleGetTagsBySlug(w http.ResponseWriter, r *http.Request) {
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

func HandleGetSlugsByTag(w http.ResponseWriter, r *http.Request) {
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

func HandlePostNewTags(w http.ResponseWriter, r *http.Request) {
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

func HandleDeleteTagRow(w http.ResponseWriter, r *http.Request) {
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

func HandleGetDimensionsBySlug(w http.ResponseWriter, r *http.Request) {
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
