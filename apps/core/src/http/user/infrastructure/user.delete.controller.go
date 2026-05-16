package infrastructure

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/user/application"
	"github.com/authuser-dev/karacron/apps/core/src/util/log"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

func UserDeleteController(app *fiber.App, db *sqlx.DB) {
	logger := log.Log("user").Named("delete")
	useCase := application.NewUserDeleteUseCase(db, logger)
	app.Delete("/user/me", useCase.UserDelete)
}

