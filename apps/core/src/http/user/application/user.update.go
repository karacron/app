package application

import (
	"encoding/json"

	"github.com/authuser-dev/karacron/apps/core/src/http/user/domain"
	dbmod "github.com/authuser-dev/karacron/apps/core/src/module/database"
	"github.com/authuser-dev/karacron/apps/core/src/server/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type UserPatchUseCase struct {
	*domain.UserUseCase
}

func NewUserPatchUseCase(db *sqlx.DB, logger *zap.Logger) *UserPatchUseCase {
	return &UserPatchUseCase{
		UserUseCase: domain.NewUserUseCase(db, logger),
	}
}

func (u *UserPatchUseCase) UserPatch(c fiber.Ctx) error {
	var req domain.UpdateUserReq
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON invÃ¡lido", "status": 400})
	}
	if errs := middleware.ValidateRequest(&req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ValidaciÃ³n fallida", "validationErrors": errs, "status": 400})
	}

	installID, err := dbmod.EnsureInstallation(u.DB)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	var id string
	if err := u.DB.Get(&id, `SELECT id FROM user_settings WHERE installation_id = ? LIMIT 1`, installID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Usuario no encontrado", "status": 404})
	}

	fields := map[string]interface{}{}

	// Handle aliased fields (with coalesce fallback to primary fields)
	firstName := req.FirstName
	if req.Name != nil && firstName == nil {
		firstName = req.Name
	}
	if firstName != nil {
		fields["first_name"] = *firstName
	}

	bio := req.Bio
	if req.Description != nil && bio == nil {
		bio = req.Description
	}
	if bio != nil {
		fields["bio"] = *bio
	}

	onboardingComplete := req.OnboardingComplete
	if req.OnboardingCompleted != nil && onboardingComplete == nil {
		onboardingComplete = req.OnboardingCompleted
	}

	// Remaining fields
	if req.LastName != nil {
		fields["last_name"] = *req.LastName
	}
	if req.DisplayName != nil {
		fields["display_name"] = *req.DisplayName
	}
	if req.AvatarURL != nil {
		fields["avatar_url"] = *req.AvatarURL
	}
	if req.Email != nil {
		fields["email"] = *req.Email
	}
	if req.BirthDate != nil {
		fields["birth_date"] = *req.BirthDate
	}
	if req.Language != nil {
		fields["language"] = *req.Language
	}
	if req.Timezone != nil {
		fields["timezone"] = *req.Timezone
	}
	if req.DateFormat != nil {
		fields["date_format"] = *req.DateFormat
	}
	if req.TimeFormat != nil {
		fields["time_format"] = *req.TimeFormat
	}
	if onboardingComplete != nil {
		v := 0
		if *onboardingComplete {
			v = 1
		}
		fields["onboarding_complete"] = v
	}
	if req.Preferences != nil {
		payload, err := json.Marshal(req.Preferences)
		if err != nil {
			u.Log.Error("marshal preferences", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Preferences inv\u00e1lidas", "status": 400})
		}
		fields["preferences"] = string(payload)
	}

	for col, val := range fields {
		if _, err := u.DB.Exec("UPDATE user_settings SET "+col+" = ? WHERE id = ?", val, id); err != nil {
			u.Log.Error("update user_settings", zap.String("col", col), zap.Error(err))
			return fiber.ErrInternalServerError
		}
	}

	var row domain.UserSettingsRow
	_ = u.DB.Get(&row, `SELECT * FROM user_settings WHERE id = ?`, id)
	return c.JSON(row)
}

