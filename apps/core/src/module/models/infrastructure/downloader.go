package infrastructure

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/authuser-dev/karacron/apps/core/src/module/models/domain"

	"go.uber.org/zap"
)

type HuggingFaceDownloader struct {
	client  *http.Client
	storage *LocalStorage
	log     *zap.Logger
}

func NewHuggingFaceDownloader(storage *LocalStorage, logger *zap.Logger) *HuggingFaceDownloader {
	return &HuggingFaceDownloader{
		client: &http.Client{Timeout: 0},
		storage: storage,
		log:     logger,
	}
}

func (d *HuggingFaceDownloader) Download(ctx context.Context, m domain.Model) (string, error) {
	if err := d.storage.EnsureModelDirectory(m); err != nil {
		return "", err
	}

	if err := d.storage.MarkIncomplete(m); err != nil {
		return "", err
	}

	targetPath := d.storage.ModelFilePath(m)
	tmpPath := targetPath + ".part"

	if err := os.Remove(tmpPath); err != nil && !os.IsNotExist(err) {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.URL, nil)
	if err != nil {
		return "", err
	}

	started := time.Now()
	resp, err := d.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("descarga fallo con status %d para %s", resp.StatusCode, m.URL)
	}

	file, err := os.Create(tmpPath)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	written, copyErr := io.Copy(io.MultiWriter(file, hasher), resp.Body)
	closeErr := file.Close()
	if copyErr != nil {
		return "", copyErr
	}
	if closeErr != nil {
		return "", closeErr
	}

	if m.SizeBytes > 0 && written != m.SizeBytes {
		return "", fmt.Errorf("tamano inesperado para %s: esperado=%d obtenido=%d", m.Name, m.SizeBytes, written)
	}

	if checksum := strings.TrimSpace(m.Checksum); checksum != "" {
		got := hex.EncodeToString(hasher.Sum(nil))
		if !strings.EqualFold(got, checksum) {
			return "", fmt.Errorf("checksum invalido para %s", m.Name)
		}
	}

	if err := os.Rename(tmpPath, targetPath); err != nil {
		return "", err
	}

	if err := d.storage.MarkComplete(m); err != nil {
		return "", err
	}

	d.log.Info("modelo descargado",
		zap.String("name", m.Name),
		zap.String("path", targetPath),
		zap.Int64("bytes", written),
		zap.Duration("duration", time.Since(started)),
	)

	return targetPath, nil
}
