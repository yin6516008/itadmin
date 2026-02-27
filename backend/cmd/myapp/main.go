package main

import (
	"flag"
	"log/slog"
	"os"

	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/infra"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("c", "config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}
	infra.InitLogger(cfg.App.Env)

	db, err := infra.NewDatabase(cfg.Database.DSN)
	if err != nil {
		slog.Error("failed to connect database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		slog.Error("failed to migrate database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler()

	router := gin.Default()
	router.Use(handler.CORS())
	router.Use(handler.RequestLogger())

	v1 := router.Group("/api/v1")
	authHandler.RegisterRoutes(v1)
	userHandler.RegisterRoutes(v1)

	slog.Info("server starting", slog.String("addr", cfg.Server.Addr))
	if err := router.Run(cfg.Server.Addr); err != nil {
		slog.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
