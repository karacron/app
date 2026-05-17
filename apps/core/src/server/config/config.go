// Paquete config: lee las variables de entorno y devuelve una struct tipada.
// Equivalente al AppConfigService de NestJS - un único lugar donde vive la config.
package config

import "os"

// Config contiene toda la configuración de la aplicación.
type Config struct {
	Port              string
	Env               string
	LogLevel          string // debug, info, warn, error
	DatabasePath      string
	WorkerConcurrency string
	DefaultProvider   string
	FallbackProvider  string
	CorsOrigins       string
	RateLimitMax      string
	RateLimitWindow   string
	InferenceEnabled  string
	InferenceModel    string
	InferenceTimeout  string
	InferenceMaxTokens string
	InferenceTemperature string
	InferenceServerHost string
	InferenceServerPort string
	InferenceServerPortMin string
	InferenceServerPortMax string
}

// Load lee las variables de entorno y devuelve la Config.
func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "3000"),
		Env:               getEnv("APP_ENV", "development"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		DatabasePath:      getEnv("DATABASE_PATH", "data/kara.db"),
		WorkerConcurrency: getEnv("WORKER_CONCURRENCY", "4"),
		DefaultProvider:   getEnv("DEFAULT_PROVIDER", "local"),
		FallbackProvider:  getEnv("FALLBACK_PROVIDER", "openrouter"),
		CorsOrigins:       getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:4000"),
		RateLimitMax:      getEnv("RATE_LIMIT_MAX", "100"),
		RateLimitWindow:   getEnv("RATE_LIMIT_WINDOW", "60"),
		InferenceEnabled:  getEnv("INFERENCE_ENABLED", "true"),
		InferenceModel:    getEnv("INFERENCE_MODEL", "qwen3-1.7b-light"),
		InferenceTimeout:  getEnv("INFERENCE_TIMEOUT_SECONDS", "120"),
		InferenceMaxTokens: getEnv("INFERENCE_MAX_TOKENS", "256"),
		InferenceTemperature: getEnv("INFERENCE_TEMPERATURE", "0.7"),
		InferenceServerHost: getEnv("INFERENCE_SERVER_HOST", "127.0.0.1"),
		InferenceServerPort: getEnv("INFERENCE_SERVER_PORT", "32111"),
		InferenceServerPortMin: getEnv("INFERENCE_SERVER_PORT_MIN", "32111"),
		InferenceServerPortMax: getEnv("INFERENCE_SERVER_PORT_MAX", "32130"),
	}
}

// getEnv es un helper privado (minúscula = no exportado fuera del paquete).
// Si la variable de entorno 'key' existe y no está vacía la devuelve,
// si no, devuelve el 'fallback'. Equivalente al patrón process.env.X ?? 'default'.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

