package server

import (
	"catify/internal/config"
	"database/sql"
	"net/http"
	"os"

	v1 "catify/internal/delivery/http/v1"
	postgres "catify/internal/repository/postgres"

	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	cfg    *config.Config
	router *http.ServeMux
	logger zerolog.Logger
	db     *sql.DB
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Initialize() error {
	a.logger = zerolog.New(os.Stdout)
	a.router = http.NewServeMux()

	if err := a.setHandler(); err != nil {
		a.logger.Error().Err(err).Msg("Failed to set handlers")
		return err
	}
	return nil
}

func (a *App) setHandler() error {

	dbConn, err := postgres.NewDbConnInstance(&a.cfg.Repository)
	if err != nil {
		a.logger.Error().Err(err).Msg("Connection to database failed")
		return err
	}
	a.db = dbConn
	a.logger.Info().Msg("Database connected successfully")

	a.router.Handle("GET /swagger/", httpSwagger.WrapHandler)

	v1.SetHandler(a.router, a.logger)

	return nil
}
