package assistant

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/assistant/infrastructure"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

func Register(app *fiber.App, db *sqlx.DB) {
	infrastructure.AssistantGetController(app, db)
	infrastructure.AssistantPatchController(app, db)
}

