package controllers

import (
	"log"
	"notification/config"

	"github.com/gofiber/fiber/v2"
)

func SendMessage(c *fiber.Ctx) error {
	// RabbitMQ bağlantısını başlat
	conn := config.ConnectRabbitmq()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// RabbitMQ'dan mesajları tüketen işlem
	go config.ConsumeMessages(ch, "myqueue")

	message := c.FormValue("message")
	config.PublishMessage(ch, "myqueue", message)
	return c.SendString("Message sent: " + message)
}