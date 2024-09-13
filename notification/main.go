package main

import (
	"log"
	"notification/database"
	"notification/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	
	// RabbitMQ'ya mesaj gönderen endpoint
	router.Notification(app)

	log.Fatal(app.Listen(":9090"))
}
