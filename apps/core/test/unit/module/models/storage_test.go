package models_test

import (
	"os"
	"testing"

	"github.com/authuser-dev/karacron/apps/core/src/module/models/domain"
	"github.com/authuser-dev/karacron/apps/core/src/module/models/infrastructure"
)

func TestStorageCompleteAndIncompleteMarkers(t *testing.T) {
	tmp := t.TempDir()
	storage := infrastructure.NewLocalStorage(tmp)

	if err := storage.CreateDirectoryIfNotExists(); err != nil {
		t.Fatalf("CreateDirectoryIfNotExists fallo: %v", err)
	}

	model := domain.Model{Name: "tiny", URL: "https://example.test/model.bin"}
	if err := storage.EnsureModelDirectory(model); err != nil {
		t.Fatalf("EnsureModelDirectory fallo: %v", err)
	}

	if err := os.WriteFile(storage.ModelFilePath(model), []byte("ok"), 0o644); err != nil {
		t.Fatalf("escritura de archivo modelo fallo: %v", err)
	}

	if err := storage.MarkIncomplete(model); err != nil {
		t.Fatalf("MarkIncomplete fallo: %v", err)
	}

	hasIncomplete, err := storage.HasIncomplete(model)
	if err != nil {
		t.Fatalf("HasIncomplete fallo: %v", err)
	}
	if !hasIncomplete {
		t.Fatal("se esperaba marcador incompleto")
	}

	if err := storage.MarkComplete(model); err != nil {
		t.Fatalf("MarkComplete fallo: %v", err)
	}

	complete, err := storage.CheckIfComplete(model)
	if err != nil {
		t.Fatalf("CheckIfComplete fallo: %v", err)
	}
	if !complete {
		t.Fatal("se esperaba modelo completo")
	}
}
