package websocket

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/authuser-dev/karacron/apps/core/src/util/log"

	"github.com/authuser-dev/karacron/apps/core/src/module/inference"
	inferencedomain "github.com/authuser-dev/karacron/apps/core/src/module/inference/domain"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const chatChunkSize = 240

var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

// wsIncoming es el sobre de todos los mensajes entrantes: {"event":"...", "data": <cualquier JSON>}
type wsIncoming struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// wsChatRequest corresponde al evento "chat": {"message":"...", "context":{...}}
type wsChatRequest struct {
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context,omitempty"`
}

type wsMessageRequest struct {
	Message string `json:"message"`
}

// send escribe un mensaje JSON en la conexiÃ³n de forma segura.
func send(conn *websocket.Conn, log *zap.Logger, clientID string, payload any) {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Error("ws: error serializando mensaje", zap.String("clientID", clientID), zap.Error(err))
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
		log.Warn("ws: error escribiendo mensaje", zap.String("clientID", clientID), zap.Error(err))
	}
}

// handleMessage procesa el evento "message": eco al emisor + broadcast al resto.
func handleMessage(conn *websocket.Conn, hub *Hub, db *sqlx.DB, log *zap.Logger, clientID string, raw json.RawMessage) {
	send(conn, log, clientID, map[string]interface{}{
		"event": "message",
		"data":  raw,
	})

	broadcast, _ := json.Marshal(map[string]interface{}{
		"event": "message",
		"from":  clientID,
		"data":  raw,
	})
	hub.Broadcast(clientID, broadcast)

	prompt := extractPrompt(raw)
	if strings.TrimSpace(prompt) == "" {
		return
	}

	response, err := inferPrompt(context.Background(), prompt)
	if err != nil {
		send(conn, log, clientID, map[string]interface{}{
			"event": "message:reply",
			"data": map[string]interface{}{
				"ok":    false,
				"error": err.Error(),
			},
		})
		return
	}

	send(conn, log, clientID, map[string]interface{}{
		"event": "message:reply",
		"data": map[string]interface{}{
			"ok":      true,
			"content": response.Content,
			"model":   response.ModelName,
		},
	})

	persistExchange(db, log, prompt, response)
}

// chunkString divide una cadena en trozos de tamaÃ±o fijo (igual que NestJS WebsocketService).
func chunkString(s string, size int) []string {
	if s == "" {
		return []string{""}
	}
	var chunks []string
	runes := []rune(s)
	for start := 0; start < len(runes); start += size {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
	}
	return chunks
}

// buildFallbackMarkdownResponse mantiene el formato markdown cuando la inferencia no esta disponible.
func buildFallbackMarkdownResponse(req wsChatRequest, inferErr error) string {
	msg := strings.TrimSpace(req.Message)
	if msg == "" {
		msg = "Sin contenido"
	}
	contextSection := ""
	if req.Context != nil {
		b, _ := json.MarshalIndent(req.Context, "", "  ")
		contextSection = fmt.Sprintf("\n\n## Contexto\n\n```json\n%s\n```", string(b))
	}
	errorSection := ""
	if inferErr != nil {
		errorSection = fmt.Sprintf("\n\n## Estado de inferencia\n\n%s", inferErr.Error())
	}
	return fmt.Sprintf("# Respuesta de chat\n\n## Mensaje recibido\n\n%s%s%s", msg, contextSection, errorSection)
}

// handleChat procesa el evento "chat": emite chunks de markdown + chat:done, devuelve chat:accepted.
// Replica exactamente la lÃ³gica de WebsocketGateway.handleChat + WebsocketService.handleChat.
func handleChat(conn *websocket.Conn, db *sqlx.DB, log *zap.Logger, clientID string, raw json.RawMessage) {
	var req wsChatRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		send(conn, log, clientID, map[string]interface{}{
			"event":   "error",
			"message": "payload invÃ¡lido para chat",
		})
		return
	}

	log.Debug("ws: chat request", zap.String("clientID", clientID), zap.String("message", req.Message))

	response, err := inferPrompt(context.Background(), req.Message)
	markdown := ""
	if err != nil {
		markdown = buildFallbackMarkdownResponse(req, err)
	} else {
		markdown = fmt.Sprintf("# Respuesta del modelo\n\n%s", strings.TrimSpace(response.Content))
		persistExchange(db, log, req.Message, response)
	}
	chunks := chunkString(markdown, chatChunkSize)

	for i, chunk := range chunks {
		send(conn, log, clientID, map[string]interface{}{
			"event": "chat:chunk",
			"data": map[string]interface{}{
				"content": chunk,
				"index":   i,
				"total":   len(chunks),
			},
		})
	}

	send(conn, log, clientID, map[string]interface{}{
		"event": "chat:done",
		"data": map[string]interface{}{
			"total":  len(chunks),
			"format": "markdown",
		},
	})

	// chat:accepted â€” respuesta final equivalente al return del gateway NestJS
	send(conn, log, clientID, map[string]interface{}{
		"event": "chat:accepted",
		"data": map[string]interface{}{
			"total":  len(chunks),
			"format": "markdown",
		},
	})
}

func inferPrompt(ctx context.Context, prompt string) (inferencedomain.InferenceResponse, error) {
	service := inference.DefaultService()
	if service == nil || !service.IsReady() {
		return inferencedomain.InferenceResponse{}, inferencedomain.ErrBinaryNotReady
	}

	return service.Generate(ctx, inferencedomain.InferenceRequest{
		ModelName: "qwen3-1.7b-light",
		Prompt:    strings.TrimSpace(prompt),
	})
}

func extractPrompt(raw json.RawMessage) string {
	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		return strings.TrimSpace(asString)
	}

	var req wsMessageRequest
	if err := json.Unmarshal(raw, &req); err == nil {
		return strings.TrimSpace(req.Message)
	}

	var generic map[string]interface{}
	if err := json.Unmarshal(raw, &generic); err == nil {
		if msg, ok := generic["message"].(string); ok {
			return strings.TrimSpace(msg)
		}
	}

	return ""
}

// RegisterRoutes aÃ±ade la ruta WebSocket /ws a la app de Fiber.
func Register(app *fiber.App, hub *Hub, db *sqlx.DB) {
	log := log.Log("ws")
	app.Get("/ws", func(c fiber.Ctx) error {
		return upgrader.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {
			clientID := uuid.New().String()

			client := &Client{ID: clientID, Conn: conn}
			hub.Register(client)
			defer hub.Unregister(clientID)

			log.Info("ws: conexiÃ³n establecida",
				zap.String("clientID", clientID),
				zap.String("remoteAddr", conn.RemoteAddr().String()),
			)

			// Bienvenida
			send(conn, log, clientID, map[string]interface{}{
				"event":    "connected",
				"clientId": clientID,
			})

			// Bucle de lectura de mensajes
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
						log.Info("ws: cliente cerrÃ³ la conexiÃ³n", zap.String("clientID", clientID))
					} else {
						log.Warn("ws: conexiÃ³n cerrada inesperadamente",
							zap.String("clientID", clientID), zap.Error(err))
					}
					return
				}

				log.Debug("ws: mensaje recibido",
					zap.String("clientID", clientID), zap.String("msg", string(msg)))

				var incoming wsIncoming
				if err := json.Unmarshal(msg, &incoming); err != nil {
					send(conn, log, clientID, map[string]interface{}{
						"event":   "error",
						"message": "JSON invÃ¡lido",
					})
					continue
				}

				switch incoming.Event {
				case "message":
					handleMessage(conn, hub, db, log, clientID, incoming.Data)
				case "chat":
					handleChat(conn, db, log, clientID, incoming.Data)
				default:
					send(conn, log, clientID, map[string]interface{}{
						"event":   "error",
						"message": fmt.Sprintf("evento desconocido: %s", incoming.Event),
					})
				}
			}
		})
	})

	log.Info("ws: ruta /ws registrada")
}

func persistExchange(db *sqlx.DB, log *zap.Logger, prompt string, response inferencedomain.InferenceResponse) {
	if db == nil {
		return
	}

	userID, assistantID, err := ensureDefaultActors(db)
	if err != nil {
		log.Warn("ws: no se pudo asegurar actores por defecto", zap.Error(err))
		return
	}

	sessionID, err := ensureActiveSession(db, userID, assistantID)
	if err != nil {
		log.Warn("ws: no se pudo asegurar sesion activa", zap.Error(err))
		return
	}

	modelID := sql.NullInt64{}
	if id, err := findModelID(db, response.ModelName); err == nil {
		modelID = sql.NullInt64{Int64: id, Valid: true}
	}

	if _, err := db.Exec(
		`INSERT INTO messages (session_id, role, content, content_type) VALUES (?, 'user', ?, 'text')`,
		sessionID,
		strings.TrimSpace(prompt),
	); err != nil {
		log.Warn("ws: no se pudo guardar mensaje user", zap.Error(err))
	}

	latencyMs := int(response.Latency.Milliseconds())
	if _, err := db.Exec(
		`INSERT INTO messages (session_id, model_id, role, content, content_type, latency_ms) VALUES (?, ?, 'assistant', ?, 'text', ?)`,
		sessionID,
		modelID,
		strings.TrimSpace(response.Content),
		latencyMs,
	); err != nil {
		log.Warn("ws: no se pudo guardar mensaje assistant", zap.Error(err))
	}
}

func ensureDefaultActors(db *sqlx.DB) (int64, int64, error) {
	var userID int64
	err := db.Get(&userID, `SELECT id FROM users ORDER BY id ASC LIMIT 1`)
	if err != nil {
		res, insErr := db.Exec(
			`INSERT INTO users (name, surname, email, birthdate) VALUES ('Kara', 'Local', 'kara.local@localhost', '2000-01-01')`,
		)
		if insErr != nil {
			return 0, 0, insErr
		}
		userID, _ = res.LastInsertId()
	}

	var assistantID int64
	err = db.Get(&assistantID, `SELECT id FROM assistants ORDER BY id ASC LIMIT 1`)
	if err != nil {
		res, insErr := db.Exec(
			`INSERT INTO assistants (reference, name, is_local, is_active, is_primary) VALUES ('local-qwen', 'Local Qwen', 1, 1, 1)`,
		)
		if insErr != nil {
			return 0, 0, insErr
		}
		assistantID, _ = res.LastInsertId()
	}

	return userID, assistantID, nil
}

func ensureActiveSession(db *sqlx.DB, userID, assistantID int64) (int64, error) {
	var sessionID int64
	err := db.Get(&sessionID, `SELECT id FROM sessions WHERE user_id = ? AND assistant_id = ? AND status = 'active' ORDER BY id DESC LIMIT 1`, userID, assistantID)
	if err == nil {
		return sessionID, nil
	}

	res, err := db.Exec(
		`INSERT INTO sessions (user_id, assistant_id, title, status) VALUES (?, ?, 'WebSocket Chat', 'active')`,
		userID,
		assistantID,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func findModelID(db *sqlx.DB, modelName string) (int64, error) {
	var modelID int64
	err := db.Get(&modelID, `SELECT id FROM models WHERE reference = ? OR name = ? ORDER BY id ASC LIMIT 1`, modelName, modelName)
	if err != nil {
		return 0, err
	}
	return modelID, nil
}



