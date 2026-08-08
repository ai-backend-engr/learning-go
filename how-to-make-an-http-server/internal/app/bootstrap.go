package app

import (
	"example.com/htmahs/internal/config"
	"example.com/htmahs/internal/database"
	"example.com/htmahs/internal/domain/user"
	"example.com/htmahs/internal/logger"
	"example.com/htmahs/internal/repository/postgres"
)

func New() (*Application, error) {
	cfg := config.Load()

	log := logger.New()

	db, err := database.New(cfg.Database)
	if err != nil {
		return nil, err
	}

	userRepository := postgres.NewUserRepository(db.DB)

	userService := user.NewService(userRepository)

	userHandler := httpTransport.NewUserHandler(userService)

	app := Application{
		Config:   cfg,
		Logger:   log,
		Database: db,
		Server:   srv,
	}

	return &app, nil
}
