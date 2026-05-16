package infrastructure

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/authuser-dev/karacron/apps/core/src/module/inference/domain"
)

type BinaryStorage struct {
	basePath string
}

func NewBinaryStorage(basePath string) *BinaryStorage {
	return &BinaryStorage{basePath: filepath.Clean(basePath)}
}

func (s *BinaryStorage) CreateBaseDirectoryIfNotExists() error {
	return os.MkdirAll(s.basePath, 0o755)
}

func (s *BinaryStorage) PlatformDirectory(asset domain.BinaryAsset) string {
	return filepath.Join(s.basePath, asset.Name, runtime.GOOS, runtime.GOARCH)
}

func (s *BinaryStorage) EnsurePlatformDirectory(asset domain.BinaryAsset) error {
	return os.MkdirAll(s.PlatformDirectory(asset), 0o755)
}

func (s *BinaryStorage) ExecutablePath(asset domain.BinaryAsset) string {
	return filepath.Join(s.PlatformDirectory(asset), asset.Executable)
}

func (s *BinaryStorage) CompleteMarkerPath(asset domain.BinaryAsset) string {
	return filepath.Join(s.PlatformDirectory(asset), ".complete")
}

func (s *BinaryStorage) IncompleteMarkerPath(asset domain.BinaryAsset) string {
	return filepath.Join(s.PlatformDirectory(asset), ".incomplete")
}

func (s *BinaryStorage) IsComplete(asset domain.BinaryAsset) (bool, error) {
	if _, err := os.Stat(s.CompleteMarkerPath(asset)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	if _, err := os.Stat(s.ExecutablePath(asset)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *BinaryStorage) MarkIncomplete(asset domain.BinaryAsset) error {
	if err := os.Remove(s.CompleteMarkerPath(asset)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.WriteFile(s.IncompleteMarkerPath(asset), []byte("status=incomplete\n"), 0o644)
}

func (s *BinaryStorage) MarkComplete(asset domain.BinaryAsset) error {
	if err := os.Remove(s.IncompleteMarkerPath(asset)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.WriteFile(s.CompleteMarkerPath(asset), []byte("status=complete\n"), 0o644)
}

func (s *BinaryStorage) PrepareExtractionDirectory(asset domain.BinaryAsset) (string, error) {
	dir := filepath.Join(s.PlatformDirectory(asset), "_extract")
	if err := os.RemoveAll(dir); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func (s *BinaryStorage) ArchivePath(asset domain.BinaryAsset) string {
	return filepath.Join(s.PlatformDirectory(asset), "download."+asset.Archive)
}
