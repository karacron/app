package models_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/authuser-dev/karacron/apps/core/src/module/models/domain"
	"github.com/authuser-dev/karacron/apps/core/src/module/models/infrastructure"

	"go.uber.org/zap"
)

func TestDownloaderCompleteDownload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("model-bytes"))
	}))
	defer server.Close()

	storage := infrastructure.NewLocalStorage(t.TempDir())
	if err := storage.CreateDirectoryIfNotExists(); err != nil {
		t.Fatalf("CreateDirectoryIfNotExists fallo: %v", err)
	}

	downloader := infrastructure.NewHuggingFaceDownloader(storage, zap.NewNop())
	model := domain.Model{Name: "tiny", URL: server.URL + "/model.bin"}

	path, err := downloader.Download(context.Background(), model)
	if err != nil {
		t.Fatalf("Download fallo: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("archivo descargado no encontrado: %v", err)
	}

	complete, err := storage.CheckIfComplete(model)
	if err != nil {
		t.Fatalf("CheckIfComplete fallo: %v", err)
	}
	if !complete {
		t.Fatal("se esperaba marcador de descarga completa")
	}
}

func TestDownloaderKeepsIncompleteMarkerOnInterruptedDownload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 32; i++ {
			_, _ = w.Write([]byte("chunk"))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(15 * time.Millisecond)
		}
	}))
	defer server.Close()

	storage := infrastructure.NewLocalStorage(t.TempDir())
	if err := storage.CreateDirectoryIfNotExists(); err != nil {
		t.Fatalf("CreateDirectoryIfNotExists fallo: %v", err)
	}

	downloader := infrastructure.NewHuggingFaceDownloader(storage, zap.NewNop())
	model := domain.Model{Name: "slow", URL: server.URL + "/slow.bin"}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()

	if _, err := downloader.Download(ctx, model); err == nil {
		t.Fatal("se esperaba error por cancelacion de contexto")
	}

	hasIncomplete, err := storage.HasIncomplete(model)
	if err != nil {
		t.Fatalf("HasIncomplete fallo: %v", err)
	}
	if !hasIncomplete {
		t.Fatal("se esperaba marcador incompleto tras falla de descarga")
	}
}
