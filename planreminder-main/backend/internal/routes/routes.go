package routes

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool, cfg config.Config) {
	authHandler := handlers.NewAuthHandler(db, cfg)
	cardHandler := handlers.NewCardHandler(db)
	subHandler := handlers.NewSubCardHandler(db)
	taskHandler := handlers.NewTaskHandler(db)
	pushHandler := handlers.NewPushHandler(db)

	api := r.Group("/api")
	{
		// public route
		api.POST("/auth/google", authHandler.GoogleLogin)
	}

	// protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Cards
		protected.GET("/cards", cardHandler.GetCards)
		protected.POST("/cards", cardHandler.CreateCard)
		protected.DELETE("/cards/:id", cardHandler.DeleteCard)

		// SubCards
		protected.GET("/cards/:card_id/subcards", subHandler.GetSubCards)
		protected.POST("/cards/:card_id/subcards", subHandler.CreateSubCard)
		protected.DELETE("/subcards/:id", subHandler.DeleteSubCard)

		// Tasks
		protected.GET("/subcards/:subcard_id/tasks", taskHandler.GetTasks)
		protected.POST("/subcards/:subcard_id/tasks", taskHandler.CreateTask)
		protected.PUT("/tasks/:id/status", taskHandler.UpdateTaskStatus)
		protected.DELETE("/tasks/:id", taskHandler.DeleteTask)

		// Push
		protected.POST("/push/subscribe", pushHandler.Subscribe)
		protected.POST("/push/unsubscribe", pushHandler.Unsubscribe)
	}
}
