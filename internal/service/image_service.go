package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"pano-redact/internal/model"
	"pano-redact/internal/repository"
	"path/filepath"
	"time"
)

type ImageService struct {
	repo      repository.ImageRepository
	uploadDir string
}

func NewImageService(repo repository.ImageRepository, uploadDir string) *ImageService {
	return &ImageService{
		repo:      repo,
		uploadDir: uploadDir,
	}
}

func (s *ImageService) generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *ImageService) UploadImage(file multipart.File, header *multipart.FileHeader) (*model.Image, error) {
	// Создаём директорию для загрузок, если её нет
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("не удалось создать директорию: %w", err)
	}

	// Генерируем уникальный ID
	id := s.generateID()
	ext := filepath.Ext(header.Filename)
	filename := id + ext

	// Сохраняем файл
	dstPath := filepath.Join(s.uploadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("не удалось скопировать файл: %w", err)
	}

	// Создаём модель изображения
	image := &model.Image{
		ID:        id,
		Filename:  header.Filename,
		URL:       "/uploads/" + filename,
		DrawURL:   "/draw/" + id,
		CreatedAt: time.Now(),
	}

	// Сохраняем в репозиторий
	if err := s.repo.Save(image); err != nil {
		return nil, fmt.Errorf("не удалось сохранить в репозиторий: %w", err)
	}

	return image, nil
}

func (s *ImageService) GetImageByID(id string) (*model.Image, error) {
	return s.repo.FindByID(id)
}

func (s *ImageService) GetAllImages() []*model.Image {
	return s.repo.FindAll()
}

func (s *ImageService) DeleteImage(id string) error {
	return s.repo.Delete(id)
}
