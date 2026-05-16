package application

import (
	"sync"

	"github.com/authuser-dev/karacron/apps/core/src/module/models/domain"
)

type ModelService struct {
	mu       sync.RWMutex
	paths    map[string]string
	statuses map[string]domain.DownloadStatus
}

func NewModelService() *ModelService {
	return &ModelService{
		paths:    make(map[string]string),
		statuses: make(map[string]domain.DownloadStatus),
	}
}

func (s *ModelService) SetStatus(name string, status domain.DownloadStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statuses[name] = status
}

func (s *ModelService) SetPath(name, path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paths[name] = path
	s.statuses[name] = domain.DownloadComplete
}

func (s *ModelService) GetModelPath(name string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	path, ok := s.paths[name]
	return path, ok
}

func (s *ModelService) IsModelReady(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.statuses[name] == domain.DownloadComplete
}
