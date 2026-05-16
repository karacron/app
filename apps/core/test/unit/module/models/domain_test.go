package models_test

import (
	"testing"

	"github.com/authuser-dev/karacron/apps/core/src/module/models/domain"
)

func TestDownloadStatusValues(t *testing.T) {
	if domain.DownloadPending == "" || domain.DownloadDownloading == "" || domain.DownloadComplete == "" || domain.DownloadFailed == "" {
		t.Fatal("los estados de descarga no deben estar vacios")
	}
}

func TestModelMetadataStoresModels(t *testing.T) {
	meta := domain.ModelMetadata{
		Models: []domain.Model{{
			Name:      "tiny-model",
			URL:       "https://example.test/model.bin",
			Version:   "v1",
			SizeBytes: 12,
			Checksum:  "abc",
		}},
	}

	if len(meta.Models) != 1 {
		t.Fatalf("se esperaba 1 modelo, se obtuvo %d", len(meta.Models))
	}
	if meta.Models[0].Name != "tiny-model" {
		t.Fatalf("nombre invalido: %s", meta.Models[0].Name)
	}
}
