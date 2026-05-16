package application

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/health/domain"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type HealthGetUseCase struct {
	*domain.HealthUseCase
}

func NewHealthGetUseCase(logger *zap.Logger) *HealthGetUseCase {
	return &HealthGetUseCase{
		HealthUseCase: domain.NewHealthUseCase(logger),
	}
}

func (u *HealthGetUseCase) HealthGet(c fiber.Ctx) error {
	u.Log.Info("health check recibido",
		zap.String("ip", c.IP()),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
	)

	response := domain.HealthCheckRes{
		Status:  "ok",
		Service: "kara-services",
	}

	return c.JSON(response)
}

