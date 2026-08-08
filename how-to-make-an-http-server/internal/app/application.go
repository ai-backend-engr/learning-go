package app

import (
	"log/slog"

	"example.com/htmahs/internal/config"
	"example.com/htmahs/internal/database"
	"example.com/htmahs/internal/server"
)

type Application struct {
	Config   config.Config
	Server   *server.Server
	Logger   *slog.Logger
	Database *database.Postgres
}

func (a *Application) Run() {}
