package handler

import (
	"html/template"
	"net/http"
	"pano-redact/internal/service"
	"path/filepath"
)

type DrawHandler struct {
	imageService *service.ImageService
	templatePath string
}

func NewDrawHandler(imageService *service.ImageService, templatePath string) *DrawHandler {
	return &DrawHandler{
		imageService: imageService,
		templatePath: templatePath,
	}
}

func (h *DrawHandler) ShowDrawPage(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL (формат: /draw/{id})
	id := r.URL.Path[len("/draw/"):]

	if id == "" {
		http.Error(w, "ID изображения не указан", http.StatusBadRequest)
		return
	}

	image, err := h.imageService.GetImageByID(id)
	if err != nil {
		http.Error(w, "Ошибка получения изображения", http.StatusInternalServerError)
		return
	}

	if image == nil {
		http.Error(w, "Изображение не найдено", http.StatusNotFound)
		return
	}

	tmplPath := filepath.Join(h.templatePath, "draw.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Не удалось загрузить шаблон", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, image)
}
