package database

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	_ "modernc.org/sqlite"
)

var DB *sqlx.DB

// Connect abre la base SQLite y configura pragmas recomendados.
func Connect(path string, log *zap.Logger) (*sqlx.DB, error) {
	// SQLite crea el archivo .db, pero no crea directorios padres.
	if err := ensureDBDir(path); err != nil {
		return nil, err
	}

	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// SQLite trabaja mejor con un writer a la vez para evitar locks.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	DB = db
	log.Info("sqlite conectado", zap.String("path", path))
	return db, nil
}

func Close(log *zap.Logger) {
	if DB == nil {
		return
	}
	if err := DB.Close(); err != nil {
		log.Error("error cerrando sqlite", zap.Error(err))
	}
}

func ensureDBDir(path string) error {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" || trimmed == ":memory:" {
		return nil
	}

	dir := filepath.Dir(trimmed)
	if dir == "." || dir == "" {
		return nil
	}

	return os.MkdirAll(dir, 0o755)
}

