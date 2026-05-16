package domain

// DownloadStatus representa el estado actual de descarga de un modelo.
type DownloadStatus string

const (
	DownloadPending     DownloadStatus = "pending"
	DownloadDownloading DownloadStatus = "downloading"
	DownloadComplete    DownloadStatus = "complete"
	DownloadFailed      DownloadStatus = "failed"
)

// Model define la metadata necesaria para descargar y resolver un modelo local.
type Model struct {
	Name      string
	URL       string
	Version   string
	SizeBytes int64
	Checksum  string
	LocalPath string
}

// ModelMetadata agrupa el listado de modelos configurados en local.
type ModelMetadata struct {
	Models []Model
}
