package infrastructure

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/system/application"
	"github.com/authuser-dev/karacron/apps/core/src/util/log"

	"github.com/gofiber/fiber/v3"
)


func ControllerSystemGet(app *fiber.App) {
	logger := log.Log("system").Named("get")
	useCase := application.NewSystemGetUseCase(logger)

	app.Get("/system", useCase.SystemGetOverview)
	app.Get("/system/os", useCase.SystemGetOs)
	app.Get("/system/cpu", useCase.SystemGetCpu)
	app.Get("/system/memory", useCase.SystemGetMemory)
	app.Get("/system/network", useCase.SystemGetNetwork)
	app.Get("/system/storage", useCase.SystemGetStorage)
	app.Get("/system/uptime", useCase.SystemGetUptime)
}

