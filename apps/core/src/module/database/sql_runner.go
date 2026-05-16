package database

import (
	"embed"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

func execSQLFile(db *sqlx.DB, filename string) error {
	raw, err := sqlFiles.ReadFile("sql/" + filename)
	if err != nil {
		return fmt.Errorf("no se pudo leer %s: %w", filename, err)
	}

	script := strings.TrimSpace(string(raw))
	if script == "" {
		return nil
	}

	// SQLite no acepta DEFAULT NOW(); usamos el equivalente portable.
	script = strings.ReplaceAll(script, "DEFAULT NOW()", "DEFAULT CURRENT_TIMESTAMP")

	if _, err := db.Exec(script); err != nil {
		return fmt.Errorf("error ejecutando %s: %w", filename, err)
	}

	return nil
}

func execSQLFiles(db *sqlx.DB, filenames []string) error {
	for _, filename := range filenames {
		if err := execSQLFile(db, filename); err != nil {
			return err
		}
	}
	return nil
}

