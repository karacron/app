package domain

import "context"

type InferenceUseCase interface {
	Generate(ctx context.Context, req InferenceRequest) (InferenceResponse, error)
	IsReady() bool
}
