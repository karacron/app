package infrastructure

import (
	"encoding/json"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		// Primera fase: permitir cualquier origen para simplificar pruebas locales.
		return true
	},
}

// WebSocketController registra la ruta GET /ws.
func WebSocketController(app *fiber.App, log *zap.Logger) {
	app.Get("/ws", func(c fiber.Ctx) error {
		return handleWebSocket(c, log)
	})
}

func handleWebSocket(c fiber.Ctx, log *zap.Logger) error {
	return upgrader.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {

		log.Info("ws: nueva conexion entrante", zap.String("remoteAddr", conn.RemoteAddr().String()))

		defer conn.Close()

		log.Info("ws: conexion aceptada", zap.String("remoteAddr", conn.RemoteAddr().String()))
		// Cada cliente recibe un identificador unico al conectar.
		clientID := uuid.New().String()

		log.Info("ws: conexion establecida",
			zap.String("clientID", clientID),
			zap.String("remoteAddr", conn.RemoteAddr().String()),
		)

		defer log.Info("ws: conexion cerrada",
			zap.String("clientID", clientID),
			zap.String("remoteAddr", conn.RemoteAddr().String()),
		)

		// Mensaje inicial para confirmar que el handshake fue correcto.
		if err := conn.WriteJSON(map[string]interface{}{
			"event": "connected",
			"data": map[string]string{
				"clientId": clientID,
			},
		}); err != nil {
			log.Error("ws: error al enviar mensaje inicial", zap.Error(err))
			return
		}

		// Bucle simple de lectura y eco.
		for {
			_, rawMessage, err := conn.ReadMessage()
			if err != nil {
				// Si falla la lectura, asumimos desconexion y cerramos.
				return
			}

			var payload map[string]interface{}
			if err := json.Unmarshal(rawMessage, &payload); err != nil {
				// Si el cliente manda algo que no es JSON, devolvemos error.
				_ = conn.WriteJSON(map[string]interface{}{
					"event": "error",
					"data": map[string]string{
						"message": "JSON invalido",
					},
				})
				continue
			}

			// Primera fase: devolvemos eco del mensaje al mismo cliente.
			if err := conn.WriteJSON(map[string]interface{}{
				"event": "echo",
				"data":  payload,
			}); err != nil {
				return
			}
		}
	})
}