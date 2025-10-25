package repository

import (
	"encoding/json"
	"os"
	"pano-redact/internal/model"
	"path/filepath"
	"sync"
)

type ImageRepository interface {
	Save(image *model.Image) error
	FindByID(id string) (*model.Image, error)
	FindAll() []*model.Image
	Delete(id string) error
	LoadFromDisk(uploadDir string) error
}

type InMemoryImageRepository struct {
	images   map[string]*model.Image
	mu       sync.RWMutex
	metaFile string
}

func NewInMemoryImageRepository(metaFile string) *InMemoryImageRepository {
	repo := &InMemoryImageRepository{
		images:   make(map[string]*model.Image),
		metaFile: metaFile,
	}
	repo.loadMetadata()
	return repo
}

func (r *InMemoryImageRepository) Save(image *model.Image) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.images[image.ID] = image
	return r.saveMetadata()
}

func (r *InMemoryImageRepository) FindByID(id string) (*model.Image, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if image, exists := r.images[id]; exists {
		return image, nil
	}
	return nil, nil
}

func (r *InMemoryImageRepository) FindAll() []*model.Image {
	r.mu.RLock()
	defer r.mu.RUnlock()

	images := make([]*model.Image, 0, len(r.images))
	for _, img := range r.images {
		images = append(images, img)
	}
	return images
}

func (r *InMemoryImageRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if image, exists := r.images[id]; exists {
		// Удаляем файл изображения
		filePath := "static" + image.URL
		os.Remove(filePath)

		// Удаляем из памяти
		delete(r.images, id)

		// Сохраняем метаданные
		return r.saveMetadata()
	}

	return nil
}

func (r *InMemoryImageRepository) LoadFromDisk(uploadDir string) error {
	return r.loadMetadata()
}

func (r *InMemoryImageRepository) saveMetadata() error {
	data, err := json.MarshalIndent(r.images, "", "  ")
	if err != nil {
		return err
	}

	// Создаём директорию, если не существует
	dir := filepath.Dir(r.metaFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(r.metaFile, data, 0644)
}

func (r *InMemoryImageRepository) loadMetadata() error {
	data, err := os.ReadFile(r.metaFile)
	if err != nil {
		// Файл не существует - это нормально при первом запуске
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var images map[string]*model.Image
	if err := json.Unmarshal(data, &images); err != nil {
		return err
	}

	r.images = images
	return nil
}
