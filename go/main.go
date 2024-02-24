package main

import (
	"log/slog"
	"net/http"
	"os"
	"powerit/db"
	"powerit/users"
	"powerit/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lmittmann/tint"
)

func main() {
	// Set up the logger
	w := os.Stderr
	var log slog.Level
	if utils.LOG_LEVEL == "info" {
		log = slog.LevelInfo
	} else {
		log = slog.LevelDebug
	}
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			AddSource:  true,
			Level:      log,
			TimeFormat: time.Kitchen,
		}),
	))

	// Connect to the database
	err := db.Connect()
	if err != nil {
		slog.Error("Error opening database", "db.Connect", err)
		panic(err)
	}
	slog.Info("Database connected")

	// Run migrations
	err = db.Migrations()
	if err != nil {
		slog.Error("Error running migrations", "db.Migrations", err)
		panic(err)
	}
	slog.Info("Migrations completed")

	// Run the HTTP server
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", utils.CLIENT_URL)
		id := 0
		err := db.Db.QueryRow("SELECT 1").Scan(&id)
		if err != nil {
			slog.Error("Error pinging database", "Db.QueryRow", err)
			return c.String(http.StatusInternalServerError, "Error pinging database")
		}
		return c.String(http.StatusOK, "Hello, World!")
	})
    e.GET("/auth", func(c echo.Context) error {
        return users.Auth(c)
    })
	e.GET("/oauth-login/:provider", func(c echo.Context) error {
		return users.OauthLogin(c)
	})
	e.GET("/oauth-callback/:provider", func(c echo.Context) error {
		return users.OauthCallback(c)
	})
	slog.Info("HTTP server listening on", "port", utils.HTTP_PORT)
	if utils.TLS == "true" {
		err = e.StartTLS(":"+utils.HTTP_PORT, utils.CERT_PATH, utils.KEY_PATH)
	} else {
		err = e.Start(":" + utils.HTTP_PORT)
	}
	if err != nil {
		slog.Error("Error serving HTTP", "e.Start", err)
		panic(err)
	}
}
