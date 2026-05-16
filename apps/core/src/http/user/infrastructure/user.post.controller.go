package infrastructure

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/user/application"
	"github.com/authuser-dev/karacron/apps/core/src/util/log"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

func UserPostController(app *fiber.App, db *sqlx.DB) {
	logger := log.Log("user").Named("post")
	useCase := application.NewUserPostUseCase(db, logger)
	app.Post("/user", useCase.UserPost)
}

