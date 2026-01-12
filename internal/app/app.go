package app

import (
	"catify/internal/config"
	"catify/internal/server"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const cfgPath = "./config/config.json"

func Start() {
	cfg, err := config.GetConfig(cfgPath)
	if err != nil {
		fmt.Println("config error:", err)
		return
	}
	app := server.NewApp(cfg)
	if err := app.Initialize(); err != nil {
		fmt.Println("app initialization error:", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-interrupt
		cancel()
	}()
	app.Run(ctx)
}
