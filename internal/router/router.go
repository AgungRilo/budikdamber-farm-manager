package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/AgungRilo/budikdamber-farm-manager/docs/swagger"
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
	}
	response.UseJSONFieldNames() // berlaku di semua environment

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// ===== Dependencies =====
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiresIn)
	userRepo := repository.NewUserRepository(db)

	authH := handler.NewAuthHandler(service.NewAuthService(userRepo, jwtManager))
	userH := handler.NewUserHandler(service.NewUserService(userRepo))
	pondH := handler.NewPondHandler(service.NewPondService(repository.NewPondRepository(db)))
	speciesH := handler.NewFishSpeciesHandler(service.NewFishSpeciesService(repository.NewFishSpeciesRepository(db)))
	feedH := handler.NewFeedHandler(service.NewFeedService(repository.NewFeedRepository(db)))
	cycleH := handler.NewCycleHandler(service.NewCycleService(repository.NewCycleRepository(db)))
	healthH := handler.NewHealthHandler(db)
	feedLogH := handler.NewFeedLogHandler(service.NewFeedLogService(repository.NewFeedLogRepository(db)))
	mortalityH := handler.NewMortalityLogHandler(service.NewMortalityLogService(repository.NewMortalityLogRepository(db)))
	saleH := handler.NewSaleHandler(service.NewSaleService(repository.NewSaleRepository(db)))

	// ===== Public =====
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", healthH.Check)

	// path persis sesuai form final project (tanpa v1)
	r.POST("/api/users/login", authH.Login)
	r.POST("/api/users/register", userH.Register)

	api := r.Group("/api/v1")
	api.POST("/users/login", authH.Login)
	api.POST("/users/register", userH.Register)

	// ===== Butuh login (status & role dicek ke DB tiap request) =====
	protected := api.Group("", middleware.AuthRequired(jwtManager, userRepo))
	ownerOnly := protected.Group("", middleware.RequireRoles(middleware.RoleOwner))

	protected.GET("/auth/me", authH.Me)

	// Users (khusus owner)
	ownerOnly.GET("/users", userH.List)
	ownerOnly.GET("/users/:id", userH.Get)
	ownerOnly.POST("/users", userH.Create)
	ownerOnly.PUT("/users/:id", userH.Update)
	ownerOnly.PATCH("/users/:id/status", userH.SetStatus)

	// Ponds
	protected.GET("/ponds", pondH.List)
	protected.GET("/ponds/:id", pondH.Get)
	ownerOnly.POST("/ponds", pondH.Create)
	ownerOnly.PUT("/ponds/:id", pondH.Update)
	ownerOnly.DELETE("/ponds/:id", pondH.Delete)

	// Fish species
	protected.GET("/fish-species", speciesH.List)
	protected.GET("/fish-species/:id", speciesH.Get)
	ownerOnly.POST("/fish-species", speciesH.Create)
	ownerOnly.PUT("/fish-species/:id", speciesH.Update)
	ownerOnly.DELETE("/fish-species/:id", speciesH.Delete)

	// Feeds
	protected.GET("/feeds", feedH.List)
	protected.GET("/feeds/:id", feedH.Get)
	ownerOnly.POST("/feeds", feedH.Create)
	ownerOnly.PUT("/feeds/:id", feedH.Update)
	ownerOnly.DELETE("/feeds/:id", feedH.Delete)

	// Cycles: semua role bisa melihat, owner yang membuat
	protected.GET("/cycles", cycleH.List)
	protected.GET("/cycles/:id", cycleH.Get)
	ownerOnly.POST("/cycles", cycleH.Create)

	// Log pakan: owner & worker mencatat, hanya owner yang menghapus
	protected.GET("/cycles/:id/feed-logs", feedLogH.List)
	protected.POST("/cycles/:id/feed-logs", feedLogH.Create)
	ownerOnly.DELETE("/feed-logs/:id", feedLogH.Delete)

	// mortality log
	protected.GET("/cycles/:id/mortality-logs", mortalityH.List)
	protected.POST("/cycles/:id/mortality-logs", mortalityH.Create)
	ownerOnly.DELETE("/mortality-logs/:id", mortalityH.Delete)

	//sales
	protected.GET("/cycles/:id/sales", saleH.List)
	protected.POST("/cycles/:id/sales", saleH.Create)
	ownerOnly.PATCH("/sales/:id/paid", saleH.SetPaid)
	ownerOnly.DELETE("/sales/:id", saleH.Delete)

	ownerOnly.POST("/cycles/:id/close", cycleH.Close)
	ownerOnly.POST("/cycles/:id/reopen", cycleH.Reopen)

	return r
}
