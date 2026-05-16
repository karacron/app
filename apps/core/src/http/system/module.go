package system

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/system/infrastructure"

	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App) {
	infrastructure.ControllerSystemGet(app)
}

