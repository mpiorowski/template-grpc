package system

import (
	"os"
)

func mustHaveEnv(key string) string {
	env := os.Getenv("ENV")
	value := os.Getenv(key)
	if value == "" && env != "test" {
		panic("Missing environment variable: " + key)
	}
	return value
}

var (
	LOG_LEVEL            = mustHaveEnv("LOG_LEVEL")
	GRPC_PORT            = mustHaveEnv("GRPC_PORT")
	JWT_SECRET           = mustHaveEnv("JWT_SECRET")
	TURSO_URL            = mustHaveEnv("TURSO_URL")
	TLS                  = mustHaveEnv("TLS")
	CERT_PATH            = mustHaveEnv("CERT_PATH")
	KEY_PATH             = mustHaveEnv("KEY_PATH")
)
