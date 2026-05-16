package data

import "github.com/authuser-dev/karacron/apps/core/src/module/inference/domain"

// LocalBinaryMetadata define binarios de llama.cpp por plataforma para descarga automatica.
var LocalBinaryMetadata = domain.BinaryMetadata{
	Assets: []domain.BinaryAsset{
		{
			Name:       "llama.cpp",
			Version:    "b9174",
			OS:         "windows",
			Arch:       "amd64",
			URL:        "https://github.com/ggml-org/llama.cpp/releases/download/b9174/llama-b9174-bin-win-cpu-x64.zip",
			Archive:    "zip",
			Executable: "llama-server.exe",
		},
		{
			Name:       "llama.cpp",
			Version:    "b9174",
			OS:         "windows",
			Arch:       "arm64",
			URL:        "https://github.com/ggml-org/llama.cpp/releases/download/b9174/llama-b9174-bin-win-cpu-arm64.zip",
			Archive:    "zip",
			Executable: "llama-server.exe",
		},
		{
			Name:       "llama.cpp",
			Version:    "b9174",
			OS:         "linux",
			Arch:       "amd64",
			URL:        "https://github.com/ggml-org/llama.cpp/releases/download/b9174/llama-b9174-bin-ubuntu-x64.tar.gz",
			Archive:    "tar.gz",
			Executable: "llama-server",
		},
		{
			Name:       "llama.cpp",
			Version:    "b9174",
			OS:         "linux",
			Arch:       "arm64",
			URL:        "https://github.com/ggml-org/llama.cpp/releases/download/b9174/llama-b9174-bin-ubuntu-arm64.tar.gz",
			Archive:    "tar.gz",
			Executable: "llama-server",
		},
		{
			Name:       "llama.cpp",
			Version:    "b9174",
			OS:         "darwin",
			Arch:       "amd64",
			URL:        "https://github.com/ggml-org/llama.cpp/releases/download/b9174/llama-b9174-bin-macos-x64.tar.gz",
			Archive:    "tar.gz",
			Executable: "llama-server",
		},
		{
			Name:       "llama.cpp",
			Version:    "b9174",
			OS:         "darwin",
			Arch:       "arm64",
			URL:        "https://github.com/ggml-org/llama.cpp/releases/download/b9174/llama-b9174-bin-macos-arm64.tar.gz",
			Archive:    "tar.gz",
			Executable: "llama-server",
		},
	},
}
