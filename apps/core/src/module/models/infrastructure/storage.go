package infrastructure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/authuser-dev/karacron/apps/core/src/module/models/domain"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

func (s *LocalStorage) CreateDirectoryIfNotExists() error {
	return os.MkdirAll(s.basePath, 0o755)
}

func (s *LocalStorage) BasePath() string {
	return s.basePath
}

func (s *LocalStorage) ModelDirectory(modelName string) string {
	return filepath.Join(s.basePath, sanitizeModelName(modelName))
}

func (s *LocalStorage) ModelFilePath(m domain.Model) string {
	fileName := fileNameFromURL(m.URL)
	if strings.TrimSpace(fileName) == "" {
		fileName = sanitizeModelName(m.Name) + ".bin"
	}
	return filepath.Join(s.ModelDirectory(m.Name), fileName)
}

func (s *LocalStorage) CompleteMarkerPath(m domain.Model) string {
	return filepath.Join(s.ModelDirectory(m.Name), ".complete")
}

func (s *LocalStorage) IncompleteMarkerPath(m domain.Model) string {
	return filepath.Join(s.ModelDirectory(m.Name), ".incomplete")
}

func (s *LocalStorage) EnsureModelDirectory(m domain.Model) error {
	return os.MkdirAll(s.ModelDirectory(m.Name), 0o755)
}

func (s *LocalStorage) CheckIfComplete(m domain.Model) (bool, error) {
	if _, err := os.Stat(s.CompleteMarkerPath(m)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	filePath := s.ModelFilePath(m)
	info, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	if m.SizeBytes > 0 && info.Size() != m.SizeBytes {
		return false, nil
	}

	return true, nil
}

func (s *LocalStorage) MarkComplete(m domain.Model) error {
	if err := os.Remove(s.IncompleteMarkerPath(m)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	marker := fmt.Sprintf("status=complete\nmodel=%s\n", m.Name)
	return os.WriteFile(s.CompleteMarkerPath(m), []byte(marker), 0o644)
}

func (s *LocalStorage) MarkIncomplete(m domain.Model) error {
	if err := os.Remove(s.CompleteMarkerPath(m)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	marker := fmt.Sprintf("status=incomplete\nmodel=%s\n", m.Name)
	return os.WriteFile(s.IncompleteMarkerPath(m), []byte(marker), 0o644)
}

func (s *LocalStorage) HasIncomplete(m domain.Model) (bool, error) {
	_, err := os.Stat(s.IncompleteMarkerPath(m))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func sanitizeModelName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "model"
	}

	replacer := strings.NewReplacer(" ", "-", "/", "-", "\\", "-", ":", "-", "*", "-", "?", "-", "\"", "-", "<", "-", ">", "-", "|", "-")
	name = replacer.Replace(name)
	return strings.ToLower(name)
}

func fileNameFromURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	if idx := strings.Index(trimmed, "?"); idx >= 0 {
		trimmed = trimmed[:idx]
	}
	if idx := strings.Index(trimmed, "#"); idx >= 0 {
		trimmed = trimmed[:idx]
	}

	base := filepath.Base(trimmed)
	if base == "." || base == "/" || base == "\\" {
		return ""
	}

	return base
}
