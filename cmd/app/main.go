package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mibrgmv/document-service/docs"
	"github.com/mibrgmv/document-service/internal/app/server"
	"github.com/mibrgmv/document-service/internal/config"
)

// @title Document Server API
// @version 1.0
// @description API for document management service
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	s := server.New(cfg)
	go func() {
		if err := s.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	sig := <-sigCh
	log.Printf("shutting down. received signal: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}
