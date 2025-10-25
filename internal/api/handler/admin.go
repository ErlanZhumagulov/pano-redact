package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type AdminHandler struct {
	templatePath string
}

func NewAdminHandler(templatePath string) *AdminHandler {
	return &AdminHandler{
		templatePath: templatePath,
	}
}

func (h *AdminHandler) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join(h.templatePath, "login.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Не удалось загрузить шаблон", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

func (h *AdminHandler) ShowAdminPage(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join(h.templatePath, "admin.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Не удалось загрузить шаблон", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}
