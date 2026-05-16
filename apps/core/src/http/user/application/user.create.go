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

// UserPostUseCase es el equivalente cercano a un service/use-case en NestJS.
// Recibe sus dependencias por composiciÃ³n y concentra la lÃ³gica del caso de uso.
type UserPostUseCase struct {
	*domain.UserUseCase
}

// NewUserPostUseCase funciona como un constructor manual.
// En Go no hay inyecciÃ³n de dependencias automÃ¡tica como en NestJS; se arma explÃ­citamente.
func NewUserPostUseCase(db *sqlx.DB, logger *zap.Logger) *UserPostUseCase {
	return &UserPostUseCase{
		UserUseCase: domain.NewUserUseCase(db, logger),
	}
}

// UserPost ejecuta el flujo de creaciÃ³n del usuario.
// Piensa en este mÃ©todo como el cuerpo de un service method llamado desde un controller.
func (u *UserPostUseCase) UserPost(c fiber.Ctx) error {
	var req domain.CreateUserReq
	// 1. Parsear el body HTTP al DTO del dominio.
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON invÃ¡lido", "status": 400})
	}
	// 2. Validar el DTO. Esto cumple un rol parecido a class-validator en NestJS.
	if errs := middleware.ValidateRequest(&req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ValidaciÃ³n fallida", "validationErrors": errs, "status": 400})
	}

	// 3. Asegurar que exista la instalaciÃ³n base sobre la que se crearÃ¡ el usuario.
	installID, err := dbmod.EnsureInstallation(u.DB)
	if err != nil {
		u.Log.Error("EnsureInstallation", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	// 4. En esta versiÃ³n solo se permite un usuario por instalaciÃ³n.
	var existing string
	if u.DB.Get(&existing, `SELECT id FROM user_settings WHERE installation_id = ? LIMIT 1`, installID) == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Solo se permite un usuario activo en esta versiÃ³n", "status": 409})
	}

	// 5. Normalizar valores opcionales antes de persistir.
	// AquÃ­ se aplican defaults y aliases de campos para no delegar esa lÃ³gica al controller.
	firstName := req.FirstName
	if firstName == nil {
		firstName = req.Name
	}
	bio := req.Bio
	if bio == nil {
		bio = req.Description
	}
	lang := "en"
	if req.Language != nil {
		lang = *req.Language
	} else if req.Locale != nil {
		lang = *req.Locale
	}
	tz := "UTC"
	if req.Timezone != nil {
		tz = *req.Timezone
	}
	onboarding := 0
	if req.OnboardingComplete != nil && *req.OnboardingComplete {
		onboarding = 1
	} else if req.OnboardingCompleted != nil && *req.OnboardingCompleted {
		onboarding = 1
	}

	var preferencesJSON *string
	if req.Preferences != nil {
		payload, err := json.Marshal(req.Preferences)
		if err != nil {
			u.Log.Error("marshal preferences", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Preferences inv\u00e1lidas", "status": 400})
		}
		value := string(payload)
		preferencesJSON = &value
	}

	// 6. Persistir el nuevo registro en la tabla user_settings.
	id := dbmod.GenID()
	_, err = u.DB.Exec(
		`INSERT INTO user_settings
			(id, installation_id, first_name, last_name, display_name, avatar_url, email, birth_date, bio, language, timezone, date_format, time_format, onboarding_complete, preferences)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, installID, firstName, req.LastName, req.DisplayName, req.AvatarURL,
		req.Email, req.BirthDate, bio, lang, tz,
		req.DateFormat, req.TimeFormat, onboarding, preferencesJSON,
	)
	if err != nil {
		u.Log.Error("insert user_settings", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "status": 400})
	}

	// 7. Reconsultar el registro creado y devolver la representaciÃ³n final al cliente.
	var row domain.UserSettingsRow
	_ = u.DB.Get(&row, `SELECT * FROM user_settings WHERE id = ?`, id)
	return c.Status(fiber.StatusCreated).JSON(row)
}

