package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrBinaryNotReady = errors.New("llama.cpp binary no disponible")
	ErrModelNotReady  = errors.New("modelo local no disponible")
	ErrEmptyPrompt    = errors.New("prompt vacio")
)

type InferenceRequest struct {
	ModelName   string
	Prompt      string
	MaxTokens   int
	Temperature float64
}

type InferenceResponse struct {
	ModelName string
	Content   string
	Latency   time.Duration
}

type BinaryAsset struct {
	Name       string
	Version    string
	OS         string
	Arch       string
	URL        string
	Archive    string
	Executable string
	Checksum   string
}

type BinaryMetadata struct {
	Assets []BinaryAsset
}

type Executor interface {
	Run(ctx context.Context, executablePath string, req InferenceRequest, modelPath string) (InferenceResponse, error)
	Close() error
}
