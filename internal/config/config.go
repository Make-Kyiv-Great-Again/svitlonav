package config

import "os"

type Config struct {
	DatabaseURL string
	ValhallaURL string
	HTTPPort    string
}

func Load() Config {
	return Config{
		DatabaseURL: mustEnv("DATABASE_URL"),
		ValhallaURL: mustEnv("VALHALLA_URL"),
		HTTPPort:    envOr("HTTP_PORT", "8080"),
	}
}

func mustEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		panic("missing required env var: " + key)
	}
	return v
}

func envOr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
