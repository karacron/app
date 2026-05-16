package data

import "github.com/authuser-dev/karacron/apps/core/src/module/models/domain"

// LocalModelMetadata define los modelos locales que se descargaran en segundo plano.
// Las URLs deben apuntar al archivo real usando /resolve/{revision}/{archivo}.
var LocalModelMetadata = domain.ModelMetadata{
	Models: []domain.Model{
		{
			Name:    "qwen3-1.7b-light",
			URL:     "https://huggingface.co/Qwen/Qwen3-1.7B-GGUF/resolve/main/Qwen3-1.7B-Q8_0.gguf",
			Version: "main",
		},
		{
			Name:    "qwen3-4b-instruct",
			URL:     "https://huggingface.co/unsloth/Qwen3-4B-Instruct-2507-GGUF/resolve/main/Qwen3-4B-Instruct-2507-Q4_K_M.gguf",
			Version: "main",
		},
		{
			Name:    "qwen2.5-coder-7b-instruct",
			URL:     "https://huggingface.co/itlwas/Qwen2.5-Coder-7B-Instruct-Q4_K_M-GGUF/resolve/main/qwen2.5-coder-7b-instruct-q4_k_m.gguf",
			Version: "main",
		},
	},
}
