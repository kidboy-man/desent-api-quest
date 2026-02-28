package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

type Config struct {
	Port      string
	JWTSecret string
	GinMode   string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", generateDefaultSecret()),
		GinMode:   getEnv("GIN_MODE", "release"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func generateDefaultSecret() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "DUMMY-JWT-SECRET"
	}
	return hex.EncodeToString(b)
}
