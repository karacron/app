package application

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/authuser-dev/karacron/apps/core/src/module/inference/domain"
	"github.com/authuser-dev/karacron/apps/core/src/module/inference/infrastructure"
	"github.com/authuser-dev/karacron/apps/core/src/module/models"

	"go.uber.org/zap"
)

type Engine struct {
	executor        domain.Executor
	binaryStorage   *infrastructure.BinaryStorage
	defaultModel    string
	defaultMaxToken int
	defaultTemp     float64
	log             *zap.Logger
}

func NewEngine(executor domain.Executor, storage *infrastructure.BinaryStorage, defaultModel string, defaultMaxToken int, defaultTemp float64, logger *zap.Logger) *Engine {
	if strings.TrimSpace(defaultModel) == "" {
		defaultModel = "qwen3-1.7b-light"
	}
	return &Engine{
		executor:        executor,
		binaryStorage:   storage,
		defaultModel:    defaultModel,
		defaultMaxToken: defaultMaxToken,
		defaultTemp:     defaultTemp,
		log:             logger,
	}
}

func (e *Engine) Generate(ctx context.Context, req domain.InferenceRequest) (domain.InferenceResponse, error) {
	modelName := strings.TrimSpace(req.ModelName)
	if modelName == "" {
		modelName = e.defaultModel
	}

	modelSvc := models.DefaultService()
	if modelSvc == nil || !modelSvc.IsModelReady(modelName) {
		return domain.InferenceResponse{}, domain.ErrModelNotReady
	}

	modelPath, ok := modelSvc.GetModelPath(modelName)
	if !ok || strings.TrimSpace(modelPath) == "" {
		return domain.InferenceResponse{}, domain.ErrModelNotReady
	}

	executablePath, err := e.resolveExecutablePath()
	if err != nil {
		return domain.InferenceResponse{}, err
	}

	if req.MaxTokens <= 0 {
		req.MaxTokens = e.defaultMaxToken
	}
	if req.Temperature <= 0 {
		req.Temperature = e.defaultTemp
	}
	if strings.TrimSpace(req.ModelName) == "" {
		req.ModelName = modelName
	}

	res, err := e.executor.Run(ctx, executablePath, req, modelPath)
	if err != nil {
		e.log.Warn("inferencia local fallo", zap.String("model", modelName), zap.Error(err))
		return domain.InferenceResponse{}, err
	}

	return res, nil
}

func (e *Engine) Shutdown() error {
	return e.executor.Close()
}

func (e *Engine) resolveExecutablePath() (string, error) {
	asset, ok := SelectCurrentAsset()
	if !ok {
		return "", domain.ErrBinaryNotReady
	}

	path := e.binaryStorage.ExecutablePath(asset)
	if strings.TrimSpace(path) == "" {
		return "", domain.ErrBinaryNotReady
	}

	if complete, err := e.binaryStorage.IsComplete(asset); err != nil || !complete {
		if err != nil {
			e.log.Warn("no se pudo validar binario", zap.Error(err))
		}
		return "", domain.ErrBinaryNotReady
	}

	return path, nil
}

func ParseDurationSeconds(raw string, fallback time.Duration) time.Duration {
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || v <= 0 {
		return fallback
	}
	return time.Duration(v) * time.Second
}

func ParseInt(raw string, fallback int) int {
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

func ParseFloat(raw string, fallback float64) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}
