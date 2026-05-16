package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func ValidateDatabase(db *sqlx.DB, log *zap.Logger) error {
	if err := validateIntegrity(db); err != nil {
		log.Error("integridad sqlite invalida", zap.Error(err))
		return err
	}

	if err := validateSchema(db); err != nil {
		log.Error("schema sqlite invalido", zap.Error(err))
		return err
	}

	log.Info("validacion de base de datos OK")
	return nil
}

func validateIntegrity(db *sqlx.DB) error {
	var result string
	if err := db.Get(&result, "PRAGMA integrity_check;"); err != nil {
		return fmt.Errorf("no se pudo ejecutar integrity_check: %w", err)
	}
	if result != "ok" {
		return fmt.Errorf("integrity_check devolvio: %s", result)
	}
	return nil
}

func validateSchema(db *sqlx.DB) error {
	required := map[string][]string{
	/* 	"installations":    {"id", "singleton_key", "name", "slug", "timezone", "locale", "created_at", "updated_at"},
		"user_settings":    {"id", "installation_id", "language", "timezone", "onboarding_complete", "created_at", "updated_at"},
		"assistant_configs": {"id", "installation_id", "assistant_name", "memory_enabled", "autonomy_enabled", "created_at", "updated_at"},
		"assistant_sessions": {"id", "installation_id", "session_type", "status", "started_at", "last_activity_at"},
		"session_messages":  {"id", "session_id", "role", "content", "content_type", "status", "created_at"}, */
	}

	for table, cols := range required {
		exists, err := tableExists(db, table)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("tabla requerida no existe: %s", table)
		}

		actualCols, err := tableColumns(db, table)
		if err != nil {
			return err
		}

		for _, requiredCol := range cols {
			if _, ok := actualCols[requiredCol]; !ok {
				return fmt.Errorf("tabla %s sin columna requerida: %s", table, requiredCol)
			}
		}
	}

	return nil
}

func tableExists(db *sqlx.DB, table string) (bool, error) {
	var count int
	err := db.Get(&count, `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?;`, table)
	if err != nil {
		return false, fmt.Errorf("error validando existencia de tabla %s: %w", table, err)
	}
	return count > 0, nil
}

func tableColumns(db *sqlx.DB, table string) (map[string]struct{}, error) {
	type pragmaRow struct {
		CID       int    `db:"cid"`
		Name      string `db:"name"`
		Type      string `db:"type"`
		NotNull   int    `db:"notnull"`
		DfltValue any    `db:"dflt_value"`
		PK        int    `db:"pk"`
	}

	rows := []pragmaRow{}
	if err := db.Select(&rows, fmt.Sprintf("PRAGMA table_info(%s);", table)); err != nil {
		return nil, fmt.Errorf("error leyendo columnas de %s: %w", table, err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("tabla %s no tiene columnas", table)
	}

	result := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		result[row.Name] = struct{}{}
	}

	return result, nil
}

