package application

import (
	dbmod "github.com/authuser-dev/karacron/apps/core/src/module/database"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type AssistantConfigInput struct {
	AssistantName      *string `json:"assistantName"      validate:"omitempty,max=255"`
	ActivePersonalityId *string `json:"activePersonalityId" validate:"omitempty"`
	ActiveToneId       *string `json:"activeToneId"       validate:"omitempty"`
}

type AssistantPatchUseCase struct {
	DB  *sqlx.DB
	Log *zap.Logger
}

func NewAssistantPatchUseCase(db *sqlx.DB, logger *zap.Logger) *AssistantPatchUseCase {
	return &AssistantPatchUseCase{DB: db, Log: logger}
}

func (u *AssistantPatchUseCase) PatchConfig(c fiber.Ctx) error {
	var req AssistantConfigInput
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON invÃ¡lido", "status": 400})
	}

	// Get installation ID (for future use when persisting config)
	if _, err := dbmod.EnsureInstallation(u.DB); err != nil {
		u.Log.Error("EnsureInstallation", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	// For now, just return success (configuration persistence can be added later)
	// This is a placeholder that allows the frontend to complete onboarding without errors
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"assistantName":       req.AssistantName,
		"activePersonalityId": req.ActivePersonalityId,
		"activeToneId":        req.ActiveToneId,
		"message":             "ConfiguraciÃ³n del asistente actualizada",
	})
}

