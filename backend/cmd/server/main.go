package main

import (
	"log"
	"os"

	"github.com/tigercasino/backend/internal/config"
	"github.com/tigercasino/backend/internal/handlers"
	"github.com/tigercasino/backend/internal/models"
	"github.com/tigercasino/backend/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	models.SeedGames(db)
	models.SeedAdmin(db, "admin@tigercasino.com", "admin123")

	router := gin.Default()

	h := handlers.NewHandler(db, cfg)

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/logout", h.Logout)
			auth.POST("/refresh", h.RefreshToken)
			auth.POST("/otp/send", h.SendOTP)
			auth.POST("/otp/verify", h.VerifyOTP)
		}

		users := api.Group("/users")
		{
			users.GET("/me", h.GetCurrentUser)
			users.PUT("/me", h.UpdateProfile)
			users.POST("/2fa/setup", h.Setup2FA)
			users.POST("/2fa/verify", h.Verify2FA)
			users.POST("/2fa/disable", h.Disable2FA)
		}

		wallet := api.Group("/wallet")
		{
			wallet.GET("/balance", h.GetBalance)
			wallet.GET("/deposit/address", h.GetDepositAddress)
			wallet.POST("/withdraw", h.Withdraw)
			wallet.GET("/transactions", h.GetTransactions)
		}

		games := api.Group("/games")
		{
			games.GET("", h.GetGames)
			games.GET("/:id", h.GetGame)
			games.POST("/:id/bet", h.PlaceBet)
			games.GET("/history", h.GetBetHistory)
		}

		admin := api.Group("/admin")
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

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
