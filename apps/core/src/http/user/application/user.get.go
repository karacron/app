package application

import (
	"context"
	"github.com/authuser-dev/karacron/apps/core/src/http/user/domain"
	dbmod "github.com/authuser-dev/karacron/apps/core/src/module/database"
	"github.com/authuser-dev/karacron/apps/core/src/server/middleware"
	"net"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type UserGetUseCase struct {
	*domain.UserUseCase
}

func NewUserGetUseCase(db *sqlx.DB, logger *zap.Logger) *UserGetUseCase {
	return &UserGetUseCase{
		UserUseCase: domain.NewUserUseCase(db, logger),
	}
}

func (u *UserGetUseCase) GetOverview(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"timestamp": dbmod.NowISO(),
		"links": fiber.Map{
			"overview": "/user",
			"status":   "/user/status",
			"me":       "/user/me",
			"create":   "/user",
			"update":   "/user/me",
			"remove":   "/user/me",
		},
	})
}

func (u *UserGetUseCase) GetStatus(c fiber.Ctx) error {
	installID, err := dbmod.EnsureInstallation(u.DB)
	if err != nil {
		u.Log.Error("EnsureInstallation", zap.Error(err))
		return c.JSON(fiber.Map{
			"userExists":          false,
			"onboardingCompleted": false,
			"userId":              nil,
			"requiresOtp":         false,
		})
	}

	var row domain.UserStatusRow
	err = u.DB.Get(&row, `SELECT id, onboarding_complete FROM user_settings WHERE installation_id = ? LIMIT 1`, installID)
	if err != nil {
		return c.JSON(fiber.Map{
			"userExists":          false,
			"onboardingCompleted": false,
			"userId":              nil,
			"requiresOtp":         false,
		})
	}

	return c.JSON(fiber.Map{
		"userExists":          true,
		"onboardingCompleted": row.OnboardingComplete == 1,
		"userId":              row.ID,
		"requiresOtp":         false,
	})
}

func (u *UserGetUseCase) GetMe(c fiber.Ctx) error {
	installID, err := dbmod.EnsureInstallation(u.DB)
	if err != nil {
		u.Log.Error("EnsureInstallation", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	var row domain.UserSettingsRow
	if err := u.DB.Get(&row, `SELECT * FROM user_settings WHERE installation_id = ? LIMIT 1`, installID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Usuario no encontrado", "status": 404})
	}

	return c.JSON(row)
}

func (u *UserGetUseCase) ValidateEmail(c fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ParÃ¡metro email requerido", "status": 400})
	}

	errs := middleware.ValidateRequest(&domain.EmailReq{Email: email})
	valid := len(errs) == 0
	if !valid {
		return c.JSON(fiber.Map{"email": email, "valid": false, "reason": "invalid_format"})
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return c.JSON(fiber.Map{"email": email, "valid": false, "reason": "invalid_format"})
	}

	domain := strings.TrimSpace(strings.ToLower(parts[1]))
	domain = strings.TrimSuffix(domain, ".")
	if domain == "" || strings.HasPrefix(domain, ".") || !strings.Contains(domain, ".") {
		return c.JSON(fiber.Map{"email": email, "valid": false, "reason": "invalid_format"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mxRecords, err := net.DefaultResolver.LookupMX(ctx, domain)
	if err != nil {
		if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.IsNotFound {
			return c.JSON(fiber.Map{"email": email, "valid": false, "reason": "no_mx_record"})
		}
		return c.JSON(fiber.Map{"email": email, "valid": false, "reason": "dns_error"})
	}

	if len(mxRecords) == 0 {
		return c.JSON(fiber.Map{"email": email, "valid": false, "reason": "no_mx_record"})
	}

	return c.JSON(fiber.Map{"email": email, "valid": true})
}

