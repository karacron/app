package domain

import "go.uber.org/zap"

type HealthUseCase struct {
	Log *zap.Logger
}

func NewHealthUseCase(logger *zap.Logger) *HealthUseCase {
	return &HealthUseCase{Log: logger}
}

