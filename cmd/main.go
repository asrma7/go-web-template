package main

import (
	"github.com/asrma7/go-web-template/internal/handlers"
	"github.com/asrma7/go-web-template/internal/repositories"
	"github.com/asrma7/go-web-template/internal/routes"
	"github.com/asrma7/go-web-template/internal/services"
	"github.com/asrma7/go-web-template/pkg/config"
	"github.com/asrma7/go-web-template/pkg/database"
	"github.com/asrma7/go-web-template/pkg/logs"
	"github.com/asrma7/go-web-template/pkg/redis"
	"github.com/asrma7/go-web-template/pkg/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	logs.InitLogger()

	cfg := config.LoadConfig()

	db, err := database.ConnectDB(cfg)
	if err != nil {
		logs.Error("Failed to connect to database", map[string]any{"error": err})
		return
	}

	redisClient := redis.InitRedisClient(cfg)
	if redisClient == nil {
		logs.Error("Failed to connect to Redis", nil)
		return
	}

	if cfg.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()
	r.Use(utils.NewCors())

	userRepo := repositories.NewUserRepository(db)

	authHandler := handlers.NewAuthHandler(services.NewAuthService(cfg, &userRepo, redisClient))

	routes.RegisterRoutes(r, authHandler)

	logs.Info("Starting server", map[string]any{
		"port": cfg.Port,
		"env":  cfg.Environment,
	})

	if err := r.Run(":" + cfg.Port); err != nil {
		logs.Logger.Fatal("Failed to start server", map[string]any{
			"error": err,
		})
	}
}
