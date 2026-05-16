package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/authuser-dev/karacron/apps/core/data"
	"github.com/authuser-dev/karacron/apps/core/src/http/assistant"
	"github.com/authuser-dev/karacron/apps/core/src/http/health"
	"github.com/authuser-dev/karacron/apps/core/src/http/system"
	"github.com/authuser-dev/karacron/apps/core/src/http/user"
	"github.com/authuser-dev/karacron/apps/core/src/module/database"
	"github.com/authuser-dev/karacron/apps/core/src/module/inference"
	"github.com/authuser-dev/karacron/apps/core/src/module/models"
	"github.com/authuser-dev/karacron/apps/core/src/module/websocket"
	"github.com/authuser-dev/karacron/apps/core/src/server/config"
	"github.com/authuser-dev/karacron/apps/core/src/server/middleware"
	logging "github.com/authuser-dev/karacron/apps/core/src/util/log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// 1. CARGAR CONFIGURACIÓN DE ENTORNO
	hasDotEnv := true
	if err := godotenv.Load(); err != nil {
		hasDotEnv = false
	}

	// 2. CARGAR CONFIGURACIÓN TIPADA
	cfg := config.Load()

	// 3. INICIALIZAR LOGGER GLOBAL
	console := logging.Log("server")
	defer console.Sync()

	if !hasDotEnv {
		console.Warn("sin .env, usando variables del entorno")
	}

	console.Info("Kara Services arrancando",
		zap.String("env", cfg.Env),
		zap.String("port", cfg.Port),
		zap.String("databasePath", cfg.DatabasePath),
		zap.String("defaultProvider", cfg.DefaultProvider),
	)

	// 4. CONECTAR Y PREPARAR BASE DE DATOS
	db, err := database.Connect(cfg.DatabasePath, logging.Log("db"))
	if err != nil {
		console.Fatal("no se pudo conectar a sqlite", zap.Error(err))
	}
	defer database.Close(logging.Log("db"))

	// 4.1 CARGAR MODELOS DE IA LOCALES
	modelModule := models.New(filepath.Join("data", "models"), data.LocalModelMetadata, logging.Log("models"))
	go modelModule.Loader.LoadAllModels(context.Background())

	// 4.2 PREPARAR BINARIO DE INFERENCIA LOCAL (llama.cpp)
	inferenceModule := inference.New(filepath.Join("data", "binaries"), data.LocalBinaryMetadata, cfg, logging.Log("inference"))
	if !strings.EqualFold(cfg.InferenceEnabled, "false") {
		go func() {
			if _, err := inferenceModule.Downloader.EnsureBinary(context.Background()); err != nil {
				inferenceModule.Service.MarkReady(false)
				console.Warn("no se pudo preparar llama.cpp", zap.Error(err))
				return
			}
			inferenceModule.Service.MarkReady(true)
			console.Info("llama.cpp listo para inferencia local")
		}()
	}
 
	// 5. EJECUTAR MIGRACIONES
	if err := database.RunMigrations(db, logging.Log("db")); err != nil {
		console.Fatal("fallaron migraciones sqlite", zap.Error(err))
	}

	// 6. VALIDAR INTEGRIDAD Y SCHEMA DE LA BASE DE DATOS
	if err := database.ValidateDatabase(db, logging.Log("db")); err != nil {
		console.Fatal("la base de datos no cuadra con el schema esperado", zap.Error(err))
	}

	// 7. CARGAR DATOS DATOS NECESARIOS INICIALES (SEED)
	if err := database.RunSeed(db, logging.Log("db")); err != nil {
		console.Warn("seed sqlite fallido", zap.Error(err))
	}

	// 8. CREAR HUB DE WEBSOCKET
	hub := websocket.NewHub(logging.Log("ws"))

	// 9. CONFIGURAR SERVIDOR FIBER
	// Crea la instancia principal del router HTTP con:
	// - Nombre de la aplicación
	// - Manejador centralizado de errores (convierte panics y errores en responses HTTP)
	app := fiber.New(fiber.Config{
		AppName: "kara-services",
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Errores 502 (Bad Gateway) suelen ser problemas de red/transporte, no de lógica
			if code == fiber.StatusBadGateway {
				console.Debug("error de red/transporte",
					zap.Int("status", code),
					zap.String("path", c.Path()),
					zap.Error(err),
				)
				return c.SendStatus(code)
			}

			// Log como error si es 5xx, como debug si es 4xx
			if code >= 500 {
				console.Error("error en handler",
					zap.Int("status", code),
					zap.String("path", c.Path()),
					zap.Error(err),
				)
			} else {
				console.Debug("respuesta no exitosa",
					zap.Int("status", code),
					zap.String("path", c.Path()),
				)
			}
			return c.Status(code).JSON(fiber.Map{
				"error":  err.Error(),
				"status": code,
			})
		},
	})

	// 10. REGISTRAR MIDDLEWARE DE RECUPERACIÓN DE PANICS
	app.Use(recover.New())

	// 11. REGISTRAR MIDDLEWARE DE SEGURIDAD
	middleware.Register(app, cfg, logging.Log("security"))

	// 12. REGISTRAR MIDDLEWARE DE LOGGING HTTP
	app.Use(middleware.HTTPLogger(logging.Log("http")))

	// 13. REGISTRAR RUTAS DE DOMINIOS DE NEGOCIO
	health.Register(app)
	assistant.Register(app, db)
	user.Register(app, db)
	system.Register(app)

	// 14. REGISTRAR RUTAS WEBSOCKET
	websocket.Register(app, hub, db)

	// 16. INICIAR SERVIDOR HTTP + MANEJO DE APAGADO LIMPIO
	addr := ":" + cfg.Port
	console.Info("Servidor HTTP listo", zap.String("addr", addr))

	listenErr := make(chan error, 1)
	go func() {
		listenErr <- app.Listen(addr)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case sig := <-sigCh:
		console.Info("senal de cierre recibida", zap.String("signal", sig.String()))
	case err := <-listenErr:
		if err != nil {
			console.Fatal("Error arrancando el servidor", zap.Error(err))
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
		console.Warn("error en shutdown de fiber", zap.Error(err))
	}

	if err := inferenceModule.Service.Shutdown(); err != nil {
		console.Warn("error cerrando inference service", zap.Error(err))
	}

	console.Info("apagado completado")
}
