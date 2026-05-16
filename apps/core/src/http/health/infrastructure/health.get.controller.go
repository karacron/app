package infrastructure

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/health/application"
	"github.com/authuser-dev/karacron/apps/core/src/util/log"

	"github.com/gofiber/fiber/v3"
)


func HealthGetController(app *fiber.App) {
	logger := log.Log("health").Named("get")
	useCase := application.NewHealthGetUseCase(logger)
	app.Get("/", useCase.HealthGet)
	app.Get("/health", useCase.HealthGet)
}

