package database

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func RunSeed(db *sqlx.DB, console *zap.Logger) error {

	console.Info("Iniciando seed")

	if err := execSQLFile(db, "03_seed.sql"); err != nil {
		return err
	}

	console.Info("seed aplicado", zap.String("file", "03_seed.sql"))
	return nil
}

