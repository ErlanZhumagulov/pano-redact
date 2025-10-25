package handler

import (
	"encoding/json"
	"net/http"
	"pano-redact/internal/service"
	"strings"
)

type APIHandler struct {
	imageService *service.ImageService
}

func NewAPIHandler(imageService *service.ImageService) *APIHandler {
	return &APIHandler{
		imageService: imageService,
	}
}

func (h *APIHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Не удалось получить файл"})
		return
	}
	defer file.Close()

	image, err := h.imageService.UploadImage(file, header)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}

func (h *APIHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	images := h.imageService.GetAllImages()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

func (h *APIHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из URL (формат: /api/images/{id})
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	id := parts[3]

	if err := h.imageService.DeleteImage(id); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Image deleted successfully"})
}
