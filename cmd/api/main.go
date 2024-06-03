package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"vulh/soundcommunity/internal/models"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type application struct {
	logger *slog.Logger
	models *models.Models
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	loadEnvFile()

	db, err := connectDB()
	if err != nil {
		logger.Error("Cannot connect to database. Error: " + err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger: logger,
		models: models.NewModels(db),
	}

	server := &http.Server{
		Addr:     fmt.Sprintf(":%v", viper.Get("PORT")),
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	logger.Info(fmt.Sprintf("Server starting on port: %v", viper.Get("PORT")), "[METHOD]", "STARTING_SERVER")
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("Cannot start server. Error: " + err.Error())
		os.Exit(1)
	}
}
