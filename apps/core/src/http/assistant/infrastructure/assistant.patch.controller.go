package infrastructure

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/assistant/application"
	"github.com/authuser-dev/karacron/apps/core/src/util/log"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

func AssistantPatchController(app *fiber.App, db *sqlx.DB) {
	logger := log.Log("assistant").Named("patch")
	useCase := application.NewAssistantPatchUseCase(db, logger)

	app.Patch("/assistant/config", useCase.PatchConfig)
}

