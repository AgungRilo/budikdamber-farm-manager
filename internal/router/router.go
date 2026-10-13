package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/AgungRilo/budikdamber-farm-manager/docs/swagger"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/auth"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/config"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/handler"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/middleware"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	pondRepo := repository.NewPondRepository(db)
	pondSvc := service.NewPondService(pondRepo)
	pondH := handler.NewPondHandler(pondSvc)
	userH := handler.NewUserHandler(service.NewUserService(userRepo))
	cycleH := handler.NewCycleHandler(service.NewCycleService(repository.NewCycleRepository(db)))
	speciesH := handler.NewFishSpeciesHandler(service.NewFishSpeciesService(repository.NewFishSpeciesRepository(db)))
	feedH := handler.NewFeedHandler(service.NewFeedService(repository.NewFeedRepository(db)))
	healthH := handler.NewHealthHandler(db)
	authH := handler.NewAuthHandler(authSvc)

	// Public
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", healthH.Check)
	api := r.Group("/api/v1")
	api.POST("/auth/login", authH.Login)

	// Butuh login
	protected := api.Group("", middleware.AuthRequired(jwtManager))
	// Khusus owner (dipakai mulai CRUD master & users)
	ownerOnly := protected.Group("", middleware.RequireRoles(middleware.RoleOwner))
	protected.GET("/ponds", pondH.List)
	protected.GET("/ponds/:id", pondH.Get)
	ownerOnly.POST("/ponds", pondH.Create)
	ownerOnly.PUT("/ponds/:id", pondH.Update)
	ownerOnly.DELETE("/ponds/:id", pondH.Delete)
	protected.GET("/auth/me", authH.Me)

	// Users (khusus owner)
	ownerOnly.GET("/users", userH.List)
	ownerOnly.GET("/users/:id", userH.Get)
	ownerOnly.POST("/users", userH.Create)
	ownerOnly.PUT("/users/:id", userH.Update)
	ownerOnly.PATCH("/users/:id/status", userH.SetStatus)

	// Siklus: semua role bisa melihat, owner yang membuat
	protected.GET("/cycles", cycleH.List)
	protected.GET("/cycles/:id", cycleH.Get)
	ownerOnly.POST("/cycles", cycleH.Create)

	// fish species
	protected.GET("/fish-species", speciesH.List)
	protected.GET("/fish-species/:id", speciesH.Get)
	ownerOnly.POST("/fish-species", speciesH.Create)
	ownerOnly.PUT("/fish-species/:id", speciesH.Update)
	ownerOnly.DELETE("/fish-species/:id", speciesH.Delete)

	// feeds
	protected.GET("/feeds", feedH.List)
	protected.GET("/feeds/:id", feedH.Get)
	ownerOnly.POST("/feeds", feedH.Create)
	ownerOnly.PUT("/feeds/:id", feedH.Update)
	ownerOnly.DELETE("/feeds/:id", feedH.Delete)

	return r
}
