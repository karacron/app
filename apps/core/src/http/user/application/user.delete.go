package application

import (
	"github.com/authuser-dev/karacron/apps/core/src/http/user/domain"
	dbmod "github.com/authuser-dev/karacron/apps/core/src/module/database"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type UserDeleteUseCase struct {
	*domain.UserUseCase
}

func NewUserDeleteUseCase(db *sqlx.DB, logger *zap.Logger) *UserDeleteUseCase {
	return &UserDeleteUseCase{
		UserUseCase: domain.NewUserUseCase(db, logger),
	}
}

func (u *UserDeleteUseCase) UserDelete(c fiber.Ctx) error {
	installID, err := dbmod.EnsureInstallation(u.DB)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	res, err := u.DB.Exec(`DELETE FROM user_settings WHERE installation_id = ?`, installID)
	if err != nil {
		u.Log.Error("delete user_settings", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Usuario no encontrado", "status": 404})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

