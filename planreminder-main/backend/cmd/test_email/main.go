package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/services"
)

func main() {
	cfg := config.Load()

	emailSvc := services.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
		cfg.SMTPFrom,
	)

	// ganti ini ke email kamu sendiri buat test
	to := "yosshana34@gmail.com"

	err := emailSvc.Send(
		to,
		"Test SMTP PlanReminder ✅",
		"Halo! Kalau email ini masuk berarti SMTP kamu sukses 🚀",
	)

	if err != nil {
		log.Fatal("FAILED:", err)
	}

	log.Println("SUCCESS: Email sent!")
}

