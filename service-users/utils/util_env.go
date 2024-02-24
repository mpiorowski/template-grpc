package utils

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
	HTTP_PORT            = mustHaveEnv("HTTP_PORT")
	GRPC_PORT            = mustHaveEnv("GRPC_PORT")
	COOKIE_DOMAIN        = mustHaveEnv("COOKIE_DOMAIN")
	CLIENT_URL           = mustHaveEnv("CLIENT_URL")
	SERVER_HTTP          = mustHaveEnv("SERVER_HTTP")
	REDIS_URL            = mustHaveEnv("REDIS_URL")
	REDIS_PASSWORD       = mustHaveEnv("REDIS_PASSWORD")
    TLS                  = mustHaveEnv("TLS")
	CERT_PATH            = mustHaveEnv("CERT_PATH")
	KEY_PATH             = mustHaveEnv("KEY_PATH")
	JWT_SECRET           = mustHaveEnv("JWT_SECRET")
	STRIPE_API_KEY       = mustHaveEnv("STRIPE_API_KEY")
	STRIPE_PRICE_ID      = mustHaveEnv("STRIPE_PRICE_ID")
	GOOGLE_CLIENT_ID     = mustHaveEnv("GOOGLE_CLIENT_ID")
	GOOGLE_CLIENT_SECRET = mustHaveEnv("GOOGLE_CLIENT_SECRET")
	GITHUB_CLIENT_ID     = mustHaveEnv("GITHUB_CLIENT_ID")
	GITHUB_CLIENT_SECRET = mustHaveEnv("GITHUB_CLIENT_SECRET")
)
