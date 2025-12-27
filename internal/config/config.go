package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            string
	Env             string
	MongoDBURI      string
	MongoDBDatabase string
	SessionSecret   string
	SessionMaxAge   int
	LogLevel        string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		Env:             getEnv("ENV", "development"),
		MongoDBURI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBDatabase: getEnv("MONGODB_DATABASE", "chessdrill"),
		SessionSecret:   getEnv("SESSION_SECRET", "change-this-to-a-secure-random-string"),
		SessionMaxAge:   getEnvInt("SESSION_MAX_AGE", 604800),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
