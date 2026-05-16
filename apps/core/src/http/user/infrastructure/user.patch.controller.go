package infrastructure

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/user/application"
	"github.com/authuser-dev/karacron/apps/core/src/util/log"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

func UserPatchController(app *fiber.App, db *sqlx.DB) {
	logger := log.Log("user").Named("patch")
	useCase := application.NewUserPatchUseCase(db, logger)
	app.Patch("/user/me", useCase.UserPatch)
}

