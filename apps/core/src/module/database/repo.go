package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const defaultInstallationReference = "default"

// GenID genera un ID Ãºnico basado en UUID v4.
func GenID() string {
	return uuid.New().String()
}

// NowISO devuelve la fecha/hora actual en formato ISO 8601 UTC.
func NowISO() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

// EnsureInstallation garantiza que exista la installation singleton 'default'.
// Devuelve el id numÃ©rico de la installation.
func EnsureInstallation(db *sqlx.DB) (int64, error) {
	var id int64
	_, err := db.Exec(
		`INSERT OR IGNORE INTO installation (reference, version, platform, is_onboarded) VALUES (?, '0.1.0', 'windows', 0)`,
		defaultInstallationReference,
	)
	if err != nil {
		return 0, fmt.Errorf("EnsureInstallation: %w", err)
	}

	// SQLite puede conservar filas antiguas con id NULL cuando el esquema usaba SERIAL.
	if _, err := db.Exec(`UPDATE installation SET id = rowid WHERE reference = ? AND id IS NULL`, defaultInstallationReference); err != nil {
		return 0, fmt.Errorf("EnsureInstallation normalize: %w", err)
	}

	err = db.Get(&id, `SELECT COALESCE(CAST(id AS INTEGER), rowid) FROM installation WHERE reference = ? LIMIT 1`, defaultInstallationReference)
	if err != nil {
		return 0, fmt.Errorf("EnsureInstallation fetch: %w", err)
	}

	return id, nil
}

