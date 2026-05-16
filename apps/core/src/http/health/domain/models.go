package domain

// HealthCheckRes es la respuesta del health check
type HealthCheckRes struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

