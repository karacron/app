package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// validatorInstance es el singleton del validador
var validatorInstance *validator.Validate

// init inicializa el validador una sola vez
func init() {
	validatorInstance = validator.New()
}

// GetValidator retorna la instancia global del validador
func GetValidator() *validator.Validate {
	return validatorInstance
}

// ValidateStruct valida un struct contra sus tags de validaciÃ³n
// Retorna un map de errores por campo si hay validaciÃ³n fallida
// Ejemplo: {"name": "Key: 'CreateUserReq.Name' Error:Field validation for 'Name' failed on the 'required' tag"}
func ValidateStruct(data interface{}) map[string]string {
	errors := make(map[string]string)

	if err := validatorInstance.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			tag := fieldError.Tag()
			errors[fieldName] = "ValidaciÃ³n fallida: " + tag
		}
	}

	return errors
}

// ValidatorMiddleware es un middleware factory que retorna un handler de Fiber
// para validar el body de una request
//
// Uso:
//   app.Post("/users", middleware.ValidatorMiddleware(CreateUserReq{}), handler)
func ValidatorMiddleware(dataType interface{}) fiber.Handler {
	return func(c fiber.Ctx) error {
		// El body ya serÃ¡ parseado por BodyParser en el handler anterior
		// Este middleware es principalmente para documentaciÃ³n/ejemplo
		return c.Next()
	}
}

// ValidateRequest es una funciÃ³n helper para validar request body dentro de un handler
//
// Uso en handler:
//   func CreateUser(c fiber.Ctx) error {
//       req := CreateUserReq{}
//       if err := c.BindJSON(&req); err != nil {
//           return c.Status(400).JSON(fiber.Map{"error": err.Error()})
//       }
//       if errors := middleware.ValidateRequest(&req); len(errors) > 0 {
//           return c.Status(400).JSON(fiber.Map{"validationErrors": errors})
//       }
//       // ... resto del handler
//   }
func ValidateRequest(req interface{}) map[string]string {
	return ValidateStruct(req)
}

// Tags de validaciÃ³n disponibles (ejemplos):
// required     - campo obligatorio
// max=N        - mÃ¡ximo N caracteres/elementos
// min=N        - mÃ­nimo N caracteres/elementos
// email        - validar formato email
// url          - validar formato URL
// numeric      - solo nÃºmeros
// alpha        - solo letras
// alphanum     - letras y nÃºmeros
// gte=N        - mayor o igual a N (para nÃºmeros)
// lte=N        - menor o igual a N (para nÃºmeros)
// gt=N         - mayor que N
// lt=N         - menor que N
// len=N        - exactamente N caracteres
// oneof=a b c  - solo valores a, b o c
//
// Ejemplo de struct:
//   type CreateUserReq struct {
//       Name     string `json:"name" validate:"required,max=100"`
//       Email    string `json:"email" validate:"required,email"`
//       Age      int    `json:"age" validate:"gte=0,lte=150"`
//       Role     string `json:"role" validate:"required,oneof=admin user guest"`
//   }

