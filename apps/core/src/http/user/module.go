package user

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/user/infrastructure"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

func Register(app *fiber.App, db *sqlx.DB) {
	infrastructure.UserGetController(app, db)
	infrastructure.UserPostController(app, db)
	infrastructure.UserPatchController(app, db)
	infrastructure.UserDeleteController(app, db)
}

