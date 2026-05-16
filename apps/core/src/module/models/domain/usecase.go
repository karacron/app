package domain

import "context"

// ModelLoaderUseCase define los contratos de carga, validacion y acceso a modelos.
type ModelLoaderUseCase interface {
	LoadAllModels(ctx context.Context)
	ValidateAndRetryIncomplete(ctx context.Context)
	GetModelPath(name string) (string, bool)
	IsModelReady(name string) bool
}
