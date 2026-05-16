// Paquete applog: logger unificado con formato coloreado para toda la app.
//
// Formato de salida:
//
//	[HH:MM:SS] [LEVEL] [scope] mensaje  {campos}
//
// Uso:
//
//	log := applog.New("health")
//	log.Info("peticiÃ³n recibida", zap.String("ip", "127.0.0.1"))
//
// Esto es un thin wrapper sobre go.uber.org/zap que configura el encoder
// con colores ANSI y el formato deseado. El tipo devuelto es *zap.Logger,
// por lo que es completamente compatible con el resto del proyecto.
package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CÃ³digos de color ANSI.
const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiGray   = "\033[90m"
	ansiGreen  = "\033[32m"
	ansiCyan   = "\033[36m"
	ansiYellow = "\033[33m"
	ansiRed    = "\033[31m"
)

func timeEnc(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(ansiGray + "[" + t.Format("15:04:05") + "]" + ansiReset)
}

func nameEnc(name string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(ansiBold + "[" + name + "]" + ansiReset)
}

// levelEnc: DEBUGâ†’cyan  INFOâ†’verde  WARNâ†’amarillo  ERRORâ†’rojo
func levelEnc(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var color, text string
	switch l {
	case zapcore.DebugLevel:
		color, text = ansiCyan, "DEBUG"
	case zapcore.InfoLevel:
		color, text = ansiGreen, "INFO"
	case zapcore.WarnLevel:
		color, text = ansiYellow, "WARN"
	case zapcore.ErrorLevel:
		color, text = ansiRed, "ERROR"
	default:
		// DPanic, Panic, Fatal
		color, text = ansiBold+ansiRed, l.CapitalString()
	}
	enc.AppendString(fmt.Sprintf("%s[%s]%s", color, text, ansiReset))
}

// New crea un *zap.Logger coloreado.
// Si no se pasa level, usa LOG_LEVEL del entorno (fallback: info).
// level soportado: "debug", "info", "warn", "error".
func Log(scope string, level ...string) *zap.Logger {
	encCfg := zapcore.EncoderConfig{
		TimeKey:          "T",
		LevelKey:         "L",
		NameKey:          "N",
		MessageKey:       "M",
		StacktraceKey:    "S",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeTime:       timeEnc,
		EncodeLevel:      levelEnc,
		EncodeName:       nameEnc,
		ConsoleSeparator: " ",
	}

	resolvedLevel := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	if resolvedLevel == "" {
		resolvedLevel = "info"
	}
	if len(level) > 0 && strings.TrimSpace(level[0]) != "" {
		resolvedLevel = strings.ToLower(strings.TrimSpace(level[0]))
	}

	// Convertir string a zapcore.Level
	var minLevel zapcore.Level
	switch resolvedLevel {
	case "debug":
		minLevel = zapcore.DebugLevel
	case "warn":
		minLevel = zapcore.WarnLevel
	case "error":
		minLevel = zapcore.ErrorLevel
	default: // "info" o cualquier otro
		minLevel = zapcore.InfoLevel
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encCfg),
		zapcore.AddSync(os.Stdout),
		minLevel,
	)

	// Stacktrace: WARN+ en debug, ERROR+ en el resto
	stacktraceLevel := zapcore.ErrorLevel
	if minLevel == zapcore.DebugLevel {
		stacktraceLevel = zapcore.WarnLevel
	}

	opts := []zap.Option{
		zap.AddStacktrace(stacktraceLevel),
	}
	if minLevel == zapcore.DebugLevel {
		opts = append(opts, zap.Development())
	}

	return zap.New(core, opts...).Named(scope)
}

