package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const ansiReset = "\033[0m"

// statusColor devuelve el cÃ³digo de status HTTP con color ANSI segÃºn su rango:
// 2xx=verde, 3xx=cyan, 4xx=amarillo, 5xx=rojo
func statusColor(code int) string {
	s := fmt.Sprintf("%d", code)
	switch {
	case code >= 200 && code < 300:
		return "\033[32m" + s + ansiReset // verde
	case code >= 300 && code < 400:
		return "\033[36m" + s + ansiReset // cyan
	case code >= 400 && code < 500:
		return "\033[33m" + s + ansiReset // amarillo
	default:
		return "\033[31m" + s + ansiReset // rojo
	}
}

// methodColor devuelve el mÃ©todo HTTP con color ANSI segÃºn su semÃ¡ntica:
// GET=verde, POST=azul, PUT=amarillo, PATCH=magenta, DELETE=rojo, HEAD/OPTIONS=gris
func methodColor(method string) string {
	switch method {
	case "GET":
		return "\033[32m[" + method + "]" + ansiReset // verde
	case "POST":
		return "\033[34m[" + method + "]" + ansiReset // azul
	case "PUT":
		return "\033[33m[" + method + "]" + ansiReset // amarillo
	case "PATCH":
		return "\033[35m[" + method + "]" + ansiReset // magenta
	case "DELETE":
		return "\033[31m[" + method + "]" + ansiReset // rojo
	case "HEAD", "OPTIONS":
		return "\033[90m[" + method + "]" + ansiReset // gris
	default:
		return "[" + method + "]"
	}
}

// HTTPLogger loguea cada peticiÃ³n HTTP con nuestro formato unificado.
func HTTPLogger(log *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		// Cuando un handler retorna un error, el ErrorHandler aÃºn no ha ejecutado,
		// por lo que c.Response().StatusCode() todavÃ­a es 200 (default).
		// El status real viene del error en sÃ­.
		status := c.Response().StatusCode()
		if err != nil {
			if e, ok := err.(*fiber.Error); ok {
				status = e.Code
			}
		}

		log.Info(fmt.Sprintf("%s %s \"%s\" (%s)",
			methodColor(c.Method()),
			statusColor(status),
			c.Path(),
			formatLatency(duration),
		))

		return err
	}
}

// formatLatency formatea una duraciÃ³n de forma legible.
// < 1ms â†’ Âµs  |  < 1s â†’ ms  |  >= 1s â†’ s
func formatLatency(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dÂµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%.1fms", float64(d.Microseconds())/1000)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

