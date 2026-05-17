package websocket

import (
	"github.com/authuser-dev/karacron/apps/core/src/module/websocket/infrastructure"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// Connection expone el punto de entrada publico del modulo websocket.
func Connection(app *fiber.App,  log *zap.Logger) {
	infrastructure.WebSocketController(app, log)
}
