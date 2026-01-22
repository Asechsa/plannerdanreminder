package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/routes"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed connect DB:", err)
	}
	defer conn.Close()

	// Email Service (SMTP)
	emailService := services.NewEmailService(
		cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom,
	)

	// Reminder scheduler (poll 60 sec)
	reminderService := services.NewReminderService(conn, emailService, 60)
	go reminderService.Start()

	r := gin.Default()
	routes.RegisterRoutes(r, conn, cfg)

	log.Println("Server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}
