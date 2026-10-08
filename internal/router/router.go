package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/auth"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/config"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/handler"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/middleware"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
)

func New(db *pgxpool.Pool, cfg config.Config) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
		response.UseJSONFieldNames()
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Dependencies
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiresIn)
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, jwtManager)

	healthH := handler.NewHealthHandler(db)
	authH := handler.NewAuthHandler(authSvc)

	// Public
	r.GET("/health", healthH.Check)
	api := r.Group("/api/v1")
	api.POST("/auth/login", authH.Login)

	// Butuh login
	protected := api.Group("", middleware.AuthRequired(jwtManager))
	// Khusus owner (dipakai mulai CRUD master & users)
	ownerOnly := protected.Group("", middleware.RequireRoles(middleware.RoleOwner))
	_ = ownerOnly
	protected.GET("/auth/me", authH.Me)

	return r
}
