package infrastructure

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/authuser-dev/karacron/apps/core/src/module/inference/domain"

	"go.uber.org/zap"
)

type BinaryDownloader struct {
	metadata domain.BinaryMetadata
	storage  *BinaryStorage
	client   *http.Client
	log      *zap.Logger
}

func NewBinaryDownloader(metadata domain.BinaryMetadata, storage *BinaryStorage, logger *zap.Logger) *BinaryDownloader {
	return &BinaryDownloader{
		metadata: metadata,
		storage:  storage,
		client:   &http.Client{},
		log:      logger,
	}
}

func (d *BinaryDownloader) EnsureBinary(ctx context.Context) (string, error) {
	asset, ok := d.selectAsset()
	if !ok {
		return "", fmt.Errorf("no hay asset de llama.cpp para %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	if err := d.storage.CreateBaseDirectoryIfNotExists(); err != nil {
		return "", err
	}
	if err := d.storage.EnsurePlatformDirectory(asset); err != nil {
		return "", err
	}

	if complete, err := d.storage.IsComplete(asset); err == nil && complete {
		if err := validateExecutableRuntime(ctx, d.storage.ExecutablePath(asset)); err == nil {
		return d.storage.ExecutablePath(asset), nil
		}
		d.log.Warn("binario marcado completo pero no ejecutable, reinstalando", zap.String("path", d.storage.ExecutablePath(asset)))
	}

	if err := d.storage.MarkIncomplete(asset); err != nil {
		return "", err
	}

	archivePath := d.storage.ArchivePath(asset)
	if err := d.downloadArchive(ctx, asset, archivePath); err != nil {
		return "", err
	}

	extractDir, err := d.storage.PrepareExtractionDirectory(asset)
	if err != nil {
		return "", err
	}

	if err := extractArchive(archivePath, extractDir, asset.Archive); err != nil {
		return "", err
	}

	foundPath, err := findFileByName(extractDir, asset.Executable)
	if err != nil {
		return "", err
	}

	runtimeDir := filepath.Dir(foundPath)
	platformDir := d.storage.PlatformDirectory(asset)
	if err := copyDirectoryContents(runtimeDir, platformDir); err != nil {
		return "", err
	}
	targetPath := d.storage.ExecutablePath(asset)

	if runtime.GOOS != "windows" {
		if err := os.Chmod(targetPath, 0o755); err != nil {
			return "", err
		}
	}

	if err := validateExecutableRuntime(ctx, targetPath); err != nil {
		return "", err
	}

	if err := d.storage.MarkComplete(asset); err != nil {
		return "", err
	}

	d.log.Info("llama.cpp listo", zap.String("path", targetPath), zap.String("version", asset.Version))
	return targetPath, nil
}

func (d *BinaryDownloader) selectAsset() (domain.BinaryAsset, bool) {
	for _, asset := range d.metadata.Assets {
		if asset.OS == runtime.GOOS && asset.Arch == runtime.GOARCH {
			return asset, true
		}
	}
	return domain.BinaryAsset{}, false
}

func (d *BinaryDownloader) downloadArchive(ctx context.Context, asset domain.BinaryAsset, archivePath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.URL, nil)
	if err != nil {
		return err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("descarga de binario fallo con status %d", resp.StatusCode)
	}

	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(io.MultiWriter(file, hasher), resp.Body); err != nil {
		return err
	}

	if checksum := strings.TrimSpace(asset.Checksum); checksum != "" {
		got := hex.EncodeToString(hasher.Sum(nil))
		if !strings.EqualFold(got, checksum) {
			return fmt.Errorf("checksum invalido para %s", asset.Name)
		}
	}

	return nil
}

func extractArchive(archivePath, targetDir, archiveType string) error {
	switch archiveType {
	case "zip":
		return unzip(archivePath, targetDir)
	case "tar.gz":
		return untarGz(archivePath, targetDir)
	default:
		return fmt.Errorf("formato de archivo no soportado: %s", archiveType)
	}
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("entrada zip invalida: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}

		in, err := f.Open()
		if err != nil {
			return err
		}

		out, err := os.Create(path)
		if err != nil {
			in.Close()
			return err
		}

		if _, err := io.Copy(out, in); err != nil {
			out.Close()
			in.Close()
			return err
		}

		if err := out.Close(); err != nil {
			in.Close()
			return err
		}
		if err := in.Close(); err != nil {
			return err
		}
	}

	return nil
}

func untarGz(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, hdr.Name)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("entrada tar invalida: %s", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			out, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			if err := out.Close(); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}

			linkTarget := strings.TrimSpace(hdr.Linkname)
			if linkTarget == "" {
				return fmt.Errorf("enlace simbolico invalido: destino vacio para %s", hdr.Name)
			}
			if filepath.IsAbs(linkTarget) {
				return fmt.Errorf("enlace simbolico invalido: destino absoluto %s", hdr.Linkname)
			}

			if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
				return err
			}
			if err := os.Symlink(linkTarget, target); err != nil {
				return err
			}
		case tar.TypeLink:
			hardlinkTarget := filepath.Join(filepath.Dir(target), hdr.Linkname)
			hardlinkTarget = filepath.Clean(hardlinkTarget)
			if !strings.HasPrefix(hardlinkTarget, filepath.Clean(dest)+string(os.PathSeparator)) {
				return fmt.Errorf("hardlink invalido: %s -> %s", hdr.Name, hdr.Linkname)
			}
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
				return err
			}
			if err := os.Link(hardlinkTarget, target); err != nil {
				return err
			}
		}
	}

	return nil
}

func findFileByName(root, name string) (string, error) {
	var found string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.EqualFold(d.Name(), name) {
			found = path
			return io.EOF
		}
		return nil
	})
	if err != nil && err != io.EOF {
		return "", err
	}
	if found == "" {
		return "", fmt.Errorf("no se encontro ejecutable %s en el archivo", name)
	}
	return found, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func copyDirectoryContents(srcDir, dstDir string) error {
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return err
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			if err := copyDirectoryContents(srcPath, dstPath); err != nil {
				return err
			}
			continue
		}

		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}

func validateExecutableRuntime(ctx context.Context, executablePath string) error {
	absPath, err := filepath.Abs(executablePath)
	if err != nil {
		return err
	}

	var lastErr error
	for attempt := 1; attempt <= 2; attempt++ {
		runCtx, cancel := context.WithTimeout(ctx, 25*time.Second)

		cmd := exec.CommandContext(runCtx, absPath, "--version")
		cmd.Dir = filepath.Dir(absPath)
		out, err := cmd.CombinedOutput()
		cancel()

		if err == nil {
			return nil
		}

		trimmed := strings.TrimSpace(string(out))
		if trimmed == "" {
			trimmed = err.Error()
		}

		if runCtx.Err() == context.DeadlineExceeded {
			lastErr = fmt.Errorf("binario invalido: validacion excedio tiempo limite: %s", trimmed)
		} else {
			lastErr = fmt.Errorf("binario invalido: %s", trimmed)
		}

		if !strings.Contains(strings.ToLower(trimmed), "signal: killed") && runCtx.Err() != context.DeadlineExceeded {
			return lastErr
		}

		time.Sleep(250 * time.Millisecond)
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("binario invalido")
	}
	return lastErr
}
