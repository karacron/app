package inference

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/authuser-dev/karacron/apps/core/src/module/inference/application"
	"github.com/authuser-dev/karacron/apps/core/src/module/inference/domain"
	"github.com/authuser-dev/karacron/apps/core/src/module/inference/infrastructure"
	"github.com/authuser-dev/karacron/apps/core/src/server/config"

	"go.uber.org/zap"
)

type Module struct {
	Downloader *infrastructure.BinaryDownloader
	Service    *application.Service
}

var (
	defaultService   *application.Service
	defaultServiceMu sync.RWMutex
)

func New(basePath string, metadata domain.BinaryMetadata, cfg *config.Config, logger *zap.Logger) *Module {
	storage := infrastructure.NewBinaryStorage(filepath.Clean(basePath))
	timeout := application.ParseDurationSeconds(cfg.InferenceTimeout, 120*time.Second)
	maxTokens := application.ParseInt(cfg.InferenceMaxTokens, 256)
	temp := application.ParseFloat(cfg.InferenceTemperature, 0.7)
	serverPort := application.ParseInt(cfg.InferenceServerPort, 32111)
	serverPortMin := application.ParseInt(cfg.InferenceServerPortMin, 32111)
	serverPortMax := application.ParseInt(cfg.InferenceServerPortMax, 32130)

	executor := infrastructure.NewLocalExecutor(
		cfg.InferenceModel,
		maxTokens,
		temp,
		timeout,
		cfg.InferenceServerHost,
		serverPort,
		serverPortMin,
		serverPortMax,
		logger,
	)
	engine := application.NewEngine(executor, storage, cfg.InferenceModel, maxTokens, temp, logger)
	service := application.NewService(engine)
	downloader := infrastructure.NewBinaryDownloader(metadata, storage, logger)

	setDefaultService(service)

	return &Module{
		Downloader: downloader,
		Service:    service,
	}
}

func DefaultService() *application.Service {
	defaultServiceMu.RLock()
	defer defaultServiceMu.RUnlock()
	return defaultService
}

func setDefaultService(service *application.Service) {
	defaultServiceMu.Lock()
	defer defaultServiceMu.Unlock()
	defaultService = service
}
