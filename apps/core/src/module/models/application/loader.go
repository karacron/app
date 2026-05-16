package application

import (
	"context"
	"sync"

	"github.com/authuser-dev/karacron/apps/core/src/module/models/domain"
	"github.com/authuser-dev/karacron/apps/core/src/module/models/infrastructure"

	"go.uber.org/zap"
)

type Loader struct {
	metadata   domain.ModelMetadata
	downloader *infrastructure.HuggingFaceDownloader
	storage    *infrastructure.LocalStorage
	service    *ModelService
	log        *zap.Logger
}

func NewLoader(
	metadata domain.ModelMetadata,
	downloader *infrastructure.HuggingFaceDownloader,
	storage *infrastructure.LocalStorage,
	service *ModelService,
	logger *zap.Logger,
) *Loader {
	return &Loader{
		metadata:   metadata,
		downloader: downloader,
		storage:    storage,
		service:    service,
		log:        logger,
	}
}

func (l *Loader) LoadAllModels(ctx context.Context) {
	if err := l.storage.CreateDirectoryIfNotExists(); err != nil {
		l.log.Error("no se pudo crear carpeta base de modelos", zap.Error(err))
		return
	}

	if len(l.metadata.Models) == 0 {
		l.log.Info("no hay modelos locales configurados")
		return
	}

	l.ValidateAndRetryIncomplete(ctx)

	var wg sync.WaitGroup
	for _, model := range l.metadata.Models {
		m := model
		wg.Add(1)
		go func() {
			defer wg.Done()

			if ok, err := l.storage.CheckIfComplete(m); err == nil && ok {
				path := l.storage.ModelFilePath(m)
				l.service.SetPath(m.Name, path)
				l.log.Info("modelo ya disponible", zap.String("name", m.Name), zap.String("path", path))
				return
			} else if err != nil {
				l.log.Warn("error validando estado de modelo", zap.String("name", m.Name), zap.Error(err))
			}

			l.service.SetStatus(m.Name, domain.DownloadDownloading)
			path, err := l.downloader.Download(ctx, m)
			if err != nil {
				l.service.SetStatus(m.Name, domain.DownloadFailed)
				l.log.Error("fallo descargando modelo", zap.String("name", m.Name), zap.Error(err))
				return
			}

			l.service.SetPath(m.Name, path)
		}()
	}
	wg.Wait()
}

func (l *Loader) ValidateAndRetryIncomplete(ctx context.Context) {
	for _, model := range l.metadata.Models {
		m := model
		hasIncomplete, err := l.storage.HasIncomplete(m)
		if err != nil {
			l.log.Warn("no se pudo validar marcador incompleto", zap.String("name", m.Name), zap.Error(err))
			continue
		}
		if !hasIncomplete {
			continue
		}

		l.log.Warn("descarga incompleta detectada, reintentando", zap.String("name", m.Name))
		l.service.SetStatus(m.Name, domain.DownloadDownloading)
		path, err := l.downloader.Download(ctx, m)
		if err != nil {
			l.service.SetStatus(m.Name, domain.DownloadFailed)
			l.log.Error("fallo reintento de descarga", zap.String("name", m.Name), zap.Error(err))
			continue
		}
		l.service.SetPath(m.Name, path)
	}
}

func (l *Loader) GetModelPath(name string) (string, bool) {
	return l.service.GetModelPath(name)
}

func (l *Loader) IsModelReady(name string) bool {
	return l.service.IsModelReady(name)
}
