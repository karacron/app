package infrastructure

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/user/application"
	"github.com/authuser-dev/karacron/apps/core/src/util/log"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

func UserGetController(app *fiber.App, db *sqlx.DB) {
	logger := log.Log("user").Named("get")
	useCase := application.NewUserGetUseCase(db, logger)
	app.Get("/user", useCase.GetOverview)
	app.Get("/user/status", useCase.GetStatus)
	app.Get("/user/me", useCase.GetMe)
	app.Get("/user/validate-email", useCase.ValidateEmail)
}

