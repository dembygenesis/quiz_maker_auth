package app

import (
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/config"
	"github.com/dembygenesis/quiz_maker_auth/src/v1/api/services/email_notification"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	app := fiber.New(fiber.Config{
		BodyLimit: 20971520,
	})

	app.Use(recover.New())

	mapUrlsV3(app)

	app.Static("/", "./public")
	app.Static("/prefix", "./public")
	app.Static("*", "./public/index.html")

	if config.IsDev == 0 {
		emailNotification := email_notification.NewEmailNotificationService()
		emailNotification.ExecuteCron()
	} else {
		fmt.Println("No one")
	}

	// ======================================================
	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		fmt.Println("Gracefully shutting down")
		app.Shutdown()
	}()

	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		fmt.Println("gg has error", err)
	}

	// ======================================================
	// Clean up tasks

	// Maybe a logger of some trigger? Or slack notif? I'm not sure... Inputs?
}