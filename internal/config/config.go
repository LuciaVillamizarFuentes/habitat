// Package config carga la configuración de la aplicación desde variables de entorno.
package config

import (
	"os"
	"time"
)

type Config struct {
	Env             string
	HTTPAddr        string
	DatabaseURL     string
	ShutdownTimeout time.Duration
}

// Load lee la configuración del entorno con valores por defecto sensatos para desarrollo.
// TODO(T-07): validar la configuración y fallar rápido si falta algo obligatorio.
func Load() Config {
	return Config{
		Env:             getEnv("APP_ENV", "development"),
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://habitat:habitat@localhost:5432/habitat?sslmode=disable"),
		ShutdownTimeout: 10 * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
