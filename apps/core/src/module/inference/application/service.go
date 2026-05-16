package application

import (
	"context"
	"sync"

	"github.com/authuser-dev/karacron/apps/core/src/module/inference/domain"
)

type Service struct {
	engine *Engine
	mu     sync.RWMutex
	ready  bool
}

func NewService(engine *Engine) *Service {
	return &Service{engine: engine}
}

func (s *Service) Generate(ctx context.Context, req domain.InferenceRequest) (domain.InferenceResponse, error) {
	return s.engine.Generate(ctx, req)
}

func (s *Service) MarkReady(ready bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ready = ready
}

func (s *Service) IsReady() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ready
}

func (s *Service) Shutdown() error {
	return s.engine.Shutdown()
}
