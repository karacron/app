package domain

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// UserUseCase contiene las dependencias comunes inyectadas para todos los casos de uso.
// Es el equivalente a un base service en NestJS que otros services heredan.
type UserUseCase struct {
	DB  *sqlx.DB
	Log *zap.Logger
}

// NewUserUseCase construye una instancia base del caso de uso.
func NewUserUseCase(db *sqlx.DB, logger *zap.Logger) *UserUseCase {
	return &UserUseCase{DB: db, Log: logger}
}

