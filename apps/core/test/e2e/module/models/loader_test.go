package models_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/authuser-dev/karacron/apps/core/src/module/models/application"
	"github.com/authuser-dev/karacron/apps/core/src/module/models/domain"
	"github.com/authuser-dev/karacron/apps/core/src/module/models/infrastructure"

	"go.uber.org/zap"
)

func TestLoaderDownloadsConfiguredModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/a.bin":
			_, _ = w.Write([]byte("aaaa"))
		case "/b.bin":
			_, _ = w.Write([]byte("bbbb"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	storage := infrastructure.NewLocalStorage(t.TempDir())
	downloader := infrastructure.NewHuggingFaceDownloader(storage, zap.NewNop())
	service := application.NewModelService()

	meta := domain.ModelMetadata{Models: []domain.Model{
		{Name: "model-a", URL: server.URL + "/a.bin", SizeBytes: 4},
		{Name: "model-b", URL: server.URL + "/b.bin", SizeBytes: 4},
	}}
	loader := application.NewLoader(meta, downloader, storage, service, zap.NewNop())

	loader.LoadAllModels(context.Background())

	for _, name := range []string{"model-a", "model-b"} {
		if !loader.IsModelReady(name) {
			t.Fatalf("modelo no listo: %s", name)
		}
		if path, ok := loader.GetModelPath(name); !ok || path == "" {
			t.Fatalf("ruta de modelo no disponible: %s", name)
		}
	}
}
