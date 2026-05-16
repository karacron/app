package database

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func RunMigrations(db *sqlx.DB, log *zap.Logger) error {
	files := []string{"01_schema.sql", "02_functions.sql", "03_seed.sql"}
	if err := execSQLFiles(db, files); err != nil {
		log.Error("migration fallida", zap.String("files", "01_schema.sql,02_functions.sql,03_seed.sql"), zap.Error(err))
		return err
	}

	log.Info("migraciones aplicadas", zap.String("files", "01_schema.sql,02_functions.sql,03_seed.sql"))

	return nil
}

