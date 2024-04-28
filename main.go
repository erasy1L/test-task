package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/erazr/test-task/config"
	db "github.com/erazr/test-task/db"
	api "github.com/erazr/test-task/http"

	"github.com/erazr/test-task/pkg"
	"github.com/erazr/test-task/services"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Println("Press Ctrl+C to exit")

	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.NewMongoDB().ConnectDB(ctx, "backtask", cfg.MongoUrl)
	if err != nil {
		log.Fatal(err)
	}

	tokenManager := pkg.NewManager(cfg)

	userRepository, err := db.NewUserRepository(database.Database().Collection("users"))
	if err != nil {
		log.Fatal(err)
	}

	userService := services.NewUserService(userRepository, cfg.RefreshTokenTTL, tokenManager)

	handler := api.NewHandler(cfg, database, userService)
	handler.RunHttp(ctx, cfg.Port, cfg.SwaggerPath)

	<-stop
	log.Println("Shutting down...")
	ctx.Done()
}
