package models

import (
	"path/filepath"
	"sync"

	"github.com/authuser-dev/karacron/apps/core/src/module/models/application"
	"github.com/authuser-dev/karacron/apps/core/src/module/models/domain"
	"github.com/authuser-dev/karacron/apps/core/src/module/models/infrastructure"

	"go.uber.org/zap"
)

type Module struct {
	Loader  *application.Loader
	Service *application.ModelService
}

var (
	defaultService   *application.ModelService
	defaultServiceMu sync.RWMutex
)

func New(basePath string, metadata domain.ModelMetadata, logger *zap.Logger) *Module {
	modelsPath := filepath.Clean(basePath)
	storage := infrastructure.NewLocalStorage(modelsPath)
	service := application.NewModelService()
	downloader := infrastructure.NewHuggingFaceDownloader(storage, logger)
	loader := application.NewLoader(metadata, downloader, storage, service, logger)
	setDefaultService(service)

	return &Module{
		Loader:  loader,
		Service: service,
	}
}

func DefaultService() *application.ModelService {
	defaultServiceMu.RLock()
	defer defaultServiceMu.RUnlock()
	return defaultService
}

func setDefaultService(service *application.ModelService) {
	defaultServiceMu.Lock()
	defer defaultServiceMu.Unlock()
	defaultService = service
}