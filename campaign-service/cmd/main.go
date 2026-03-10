package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/likhon22/ad-system/campaign-service/config"
	httpserver "github.com/likhon22/ad-system/campaign-service/internal/adapter/inbound/http"
)

func main() {

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	app := httpserver.NewApp(cfg)
	go func() {
		if err := app.Run(); err != nil {
			log.Fatal("err", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		log.Fatal("forced shutdown", "err", err)

	}
	log.Println("server exited cleanly")
}
