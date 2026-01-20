package server

import (
	"catify/internal/config"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	v1 "catify/internal/delivery/http/v1"
	postgres "catify/internal/repository/postgres"
	service "catify/internal/service"

	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	cfg         *config.Config
	router      *http.ServeMux
	logger      zerolog.Logger
	userService v1.UserService
	goalService v1.GoalService
	db          *sql.DB
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

	dbConn, err := postgres.NewDbConnInstance(a.cfg.Repository)
	if err != nil {
		a.logger.Error().Err(err).Msg("Connection to database failed")
		return err
	}
	a.db = dbConn
	a.logger.Info().Msg("Database connected successfully")

	userRepo := postgres.NewUserRepository(a.db)
	goalRepo := postgres.NewGoalRepository(a.db)

	a.userService = service.NewUserService(userRepo, a.logger)
	a.goalService = service.NewGoalService(goalRepo, a.logger)

	a.router.Handle("GET /swagger/", httpSwagger.WrapHandler)

	v1.SetHandler(a.logger, a.userService, a.goalService, a.router)

	return nil
}

func (a *App) Run(ctx context.Context) {
	var err error
	defer a.closeConnections()

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", a.cfg.App.Port),
		Handler: a.router,
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error().Err(err).Msg("Server failed")
			cancel()
			return
		}
	}()

	a.logger.Info().Msgf("Server is running on %s", a.cfg.App.Port)
	<-ctx.Done()
	a.logger.Info().Msg("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		a.logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	a.logger.Info().Msg("Server exited properly")
}

func (a *App) closeConnections() {
	if err := a.db.Close(); err != nil {
		a.logger.Error().Err(err).Msg("failed to close db")
	}
}
