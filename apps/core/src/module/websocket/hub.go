// Paquete websocket: mÃ³dulo que gestiona las conexiones WebSocket activas.
// El Hub es el "registro central" de clientes conectados â€” equivalente al
// Gateway de NestJS con @WebSocketGateway pero sin decoradores, todo explÃ­cito.
package websocket

import (
	"sync" // sync.RWMutex: permite mÃºltiples lecturas simultÃ¡neas, escritura exclusiva

	// fasthttp/websocket es la librerÃ­a WebSocket base compatible con Fiber v3.
	// gofiber/contrib/websocket solo funciona con Fiber v2.
	"github.com/fasthttp/websocket"
	"go.uber.org/zap"
)

// Client representa una conexiÃ³n WebSocket activa.
// Cada cliente que se conecta a /ws crea una instancia de esta struct.
type Client struct {
	ID   string          // UUID Ãºnico asignado al conectarse
	Conn *websocket.Conn // La conexiÃ³n WebSocket real (leer/escribir mensajes)
}

// Hub mantiene el mapa de todos los clientes conectados.
// Es un singleton: se crea una vez en main.go y se pasa a las rutas.
//
// Â¿Por quÃ© RWMutex en vez de Mutex?
//   - RLock/RUnlock: mÃºltiples goroutines pueden leer el mapa al mismo tiempo.
//   - Lock/Unlock:   solo una goroutine puede escribir (conectar/desconectar).
// Esto evita bloqueos innecesarios cuando solo se estÃ¡ leyendo (ej: broadcast).
type Hub struct {
	mu      sync.RWMutex       // protege el mapa de clientes contra race conditions
	clients map[string]*Client // clientID â†’ *Client
	log     *zap.Logger
}

// NewHub crea e inicializa el Hub. Se llama desde main.go.
// El map vacÃ­o se crea aquÃ­ para que nunca sea nil.
func NewHub(log *zap.Logger) *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		log:     log,
	}
}

// Register aÃ±ade un cliente al mapa cuando se conecta.
// Usa Lock() (escritura exclusiva) porque modifica el mapa.
func (h *Hub) Register(client *Client) {
	h.mu.Lock()         // bloquea escritura
	defer h.mu.Unlock() // libera al salir de la funciÃ³n (defer = finally)

	h.clients[client.ID] = client

	h.log.Info("cliente WebSocket conectado",
		zap.String("clientID", client.ID),
		zap.Int("total_clientes", len(h.clients)),
	)
}

// Unregister elimina un cliente del mapa cuando se desconecta.
// TambiÃ©n usa Lock() porque modifica el mapa.
func (h *Hub) Unregister(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, clientID) // delete(map, clave) â€” equivalente a Map.delete() en JS

	h.log.Info("cliente WebSocket desconectado",
		zap.String("clientID", clientID),
		zap.Int("total_clientes", len(h.clients)),
	)
}

// Broadcast envÃ­a un mensaje a todos los clientes conectados EXCEPTO al emisor.
// Usa RLock() (lectura compartida) porque solo recorre el mapa sin modificarlo.
func (h *Hub) Broadcast(senderID string, message []byte) {
	h.mu.RLock()         // bloquea solo lectura (otros tambiÃ©n pueden leer al mismo tiempo)
	defer h.mu.RUnlock()

	for id, client := range h.clients {
		// Saltar al emisor: no hace falta que reciba su propio mensaje
		if id == senderID {
			continue
		}

		// WriteMessage envÃ­a el mensaje por la conexiÃ³n WebSocket.
		// fiberws.TextMessage = tipo 1 (texto UTF-8). TambiÃ©n existe BinaryMessage = 2.
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			h.log.Warn("error al enviar mensaje a cliente",
				zap.String("clientID", id),
				zap.Error(err),
			)
		}
	}

	h.log.Debug("broadcast enviado",
		zap.String("senderID", senderID),
		zap.Int("receptores", len(h.clients)-1),
	)
}

// Count devuelve el nÃºmero de clientes conectados actualmente.
// Ãštil para logs y para el endpoint /health si queremos incluirlo.
func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

