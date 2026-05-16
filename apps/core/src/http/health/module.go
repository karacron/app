package health

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/health/infrastructure"

	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App) {
	infrastructure.HealthGetController(app)
}

