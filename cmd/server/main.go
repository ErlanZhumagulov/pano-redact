package main

import (
	"fmt"
	"pano-redact/internal/api/handler"
	"pano-redact/internal/config"
	"pano-redact/internal/middleware"
	"pano-redact/internal/repository"
	"pano-redact/internal/service"
	"strings"

	"log"
	"net/http"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализируем слои
	metaFile := cfg.UploadDir + "/metadata.json"
	imageRepo := repository.NewInMemoryImageRepository(metaFile)
	imageService := service.NewImageService(imageRepo, cfg.UploadDir)

	// Инициализируем middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.AdminUsername, cfg.AdminPassword)

	// Инициализируем handlers
	adminHandler := handler.NewAdminHandler(cfg.StaticDir)
	apiHandler := handler.NewAPIHandler(imageService)
	drawHandler := handler.NewDrawHandler(imageService, cfg.StaticDir)

	// Настраиваем роуты
	mux := http.NewServeMux()

	// Страница логина (публичная)
	mux.HandleFunc("/login", adminHandler.ShowLoginPage)

	// API логин
	mux.HandleFunc("/api/login", authMiddleware.Login)

	// Админ панель (требует авторизации)
	mux.HandleFunc("/admin", authMiddleware.RequireAuth(adminHandler.ShowAdminPage))

	// Главная страница редиректит на логин
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	// API endpoints (требуют авторизации)
	mux.HandleFunc("/api/upload", authMiddleware.RequireAuth(apiHandler.UploadImage))
	mux.HandleFunc("/api/images", authMiddleware.RequireAuth(apiHandler.ListImages))

	// Обработчик для удаления изображений
	mux.HandleFunc("/api/images/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/images/") && r.Method == http.MethodDelete {
			authMiddleware.RequireAuth(apiHandler.DeleteImage)(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// Страница рисования (без авторизации)
	mux.HandleFunc("/draw/", drawHandler.ShowDrawPage)

	// Статические файлы (загруженные изображения)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.UploadDir))))

	// Запускаем сервер
	fmt.Println("🚀 Сервер запущен на http://localhost" + cfg.ServerPort)
	fmt.Println("👤 Логин:", cfg.AdminUsername)
	fmt.Println("🔑 Пароль:", cfg.AdminPassword)
	fmt.Println("📝 Метаданные сохраняются в:", metaFile)
	fmt.Println()

	if err := http.ListenAndServe(cfg.ServerPort, mux); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
