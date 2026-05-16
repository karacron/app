package infrastructure

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/assistant/application"
	"github.com/authuser-dev/karacron/apps/core/src/util/log"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

func AssistantGetController(app *fiber.App, db *sqlx.DB) {
	logger := log.Log("assistant").Named("get")
	useCase := application.NewAssistantGetUseCase(db, logger)

	app.Get("/assistant/personalities", useCase.GetPersonalities)
	app.Get("/assistant/tones", useCase.GetTones)
}

