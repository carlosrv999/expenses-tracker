package config

import "os"

type Config struct {
	DBURL          string
	ServerAddress  string
	MigrationsPath string
}

func Load() *Config {
	return &Config{
		DBURL:          getenv("DATABASE_URL", "postgres://expense_app:my_secure_password@localhost:5432/expense_tracker?sslmode=disable"),
		ServerAddress:  getenv("SERVER_ADDRESS", ":8080"),
		MigrationsPath: getenv("MIGRATIONS_PATH", "file://migrations"),
	}
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
