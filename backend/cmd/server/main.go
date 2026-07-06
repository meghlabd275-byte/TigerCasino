package main

import (
	"log"
	"os"

	"github.com/tigercasino/backend/internal/config"
	"github.com/tigercasino/backend/internal/handlers"
	"github.com/tigercasino/backend/internal/middleware"
	"github.com/tigercasino/backend/internal/models"
	"github.com/tigercasino/backend/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Transaction{},
		&models.Game{},
		&models.Bet{},
		&models.Session{},
		&models.AuditLog{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed default games
	models.SeedGames(db)

	// Seed admin user
	models.SeedAdmin(db, cfg.AdminEmail, cfg.AdminPassword)

	// Initialize Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Initialize handlers
	h := handlers.NewHandler(db, cfg)

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/logout", h.Logout)
			auth.POST("/refresh", h.RefreshToken)
			auth.POST("/otp/send", h.SendOTP)
			auth.POST("/otp/verify", h.VerifyOTP)
		}

		// User routes (protected)
		users := api.Group("/users")
		users.Use(middleware.Auth(cfg.JWTSecret))
		{
			users.GET("/me", h.GetCurrentUser)
			users.PUT("/me", h.UpdateProfile)
			users.POST("/2fa/setup", h.Setup2FA)
			users.POST("/2fa/verify", h.Verify2FA)
			users.POST("/2fa/disable", h.Disable2FA)
		}

		// Wallet routes (protected)
		wallet := api.Group("/wallet")
		wallet.Use(middleware.Auth(cfg.JWTSecret))
		{
			wallet.GET("/balance", h.GetBalance)
			wallet.GET("/deposit/address", h.GetDepositAddress)
			wallet.POST("/withdraw", h.Withdraw)
			wallet.GET("/transactions", h.GetTransactions)
		}

		// Game routes (protected)
		games := api.Group("/games")
		{
			games.GET("", h.GetGames)
			games.GET("/:id", h.GetGame)
			games.POST("/:id/bet", middleware.Auth(cfg.JWTSecret), h.PlaceBet)
			games.GET("/history", middleware.Auth(cfg.JWTSecret), h.GetBetHistory)
		}

		// Admin routes (admin only)
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(cfg.JWTSecret))
		admin.Use(middleware.AdminOnly())
		{
			admin.GET("/dashboard", h.GetAdminStats)
			admin.GET("/users", h.GetUsers)
			admin.PUT("/users/:id", h.UpdateUser)
			admin.GET("/transactions", h.GetAllTransactions)
			admin.POST("/transactions/:id/approve", h.ApproveTransaction)
			admin.POST("/transactions/:id/reject", h.RejectTransaction)
			admin.GET("/games", h.GetAdminGames)
			admin.PUT("/games/:id", h.UpdateGame)
			admin.GET("/audit-logs", h.GetAuditLogs)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
