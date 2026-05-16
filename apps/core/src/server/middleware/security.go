package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"go.uber.org/zap"

	"github.com/authuser-dev/karacron/apps/core/src/server/config"
)

// Register aplica los middlewares de seguridad en el orden correcto:
// helmet â†’ cors â†’ limiter
func Register(app *fiber.App, cfg *config.Config, log *zap.Logger) {

	// 1. Helmet: cabeceras HTTP de seguridad.
	// Activa por defecto: X-XSS-Protection, X-Content-Type-Options: nosniff,
	// X-Frame-Options: SAMEORIGIN, Referrer-Policy: no-referrer, HSTS.
	app.Use(helmet.New())

	// 2. CORS: controla quÃ© orÃ­genes pueden llamar a la API.
	// AllowCredentials permite cookies/tokens de sesiÃ³n desde el frontend.
	origins := strings.Split(cfg.CorsOrigins, ",")
	app.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           3600, // cachÃ© del preflight: 1 hora
	}))

	// 3. Rate Limiter: mÃ¡ximo N peticiones por ventana de tiempo por IP.
	// Cuando se supera el lÃ­mite responde 429 en vez de procesar la peticiÃ³n.
	max, _ := strconv.Atoi(cfg.RateLimitMax)
	windowSec, _ := strconv.Atoi(cfg.RateLimitWindow)
	app.Use(limiter.New(limiter.Config{
		Max:        max,
		Expiration: time.Duration(windowSec) * time.Second,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP() // lÃ­mite por IP
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":  "too many requests",
				"status": fiber.StatusTooManyRequests,
			})
		},
	}))

	log.Info("seguridad registrada",
		zap.String("corsOrigins", cfg.CorsOrigins),
		zap.String("rateLimitMax", cfg.RateLimitMax),
		zap.String("rateLimitWindow", cfg.RateLimitWindow+"s"),
	)
}

