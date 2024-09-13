package controllers

import (
	"notification/config"
	"notification/database"
	"notification/models"

	"github.com/gofiber/fiber/v2"
)

func SendMessage(c *fiber.Ctx) error {
	db := database.DB.Db

	var req models.MessageLog
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	// RabbitMQ bağlantısını başlat
	conn := config.ConnectRabbitmq()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open a channel"})
	}
	defer ch.Close()

	// RabbitMQ'dan mesajları tüketen işlem
	go config.ConsumeMessages(ch, "myqueue")

	// Gelen JSON mesajı kuyruğa gönder
	if err := config.PublishMessage(ch, "myqueue", req.Message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to publish message"})
	}

	logEntry := models.MessageLog{
		QueueName: "myqueue",
		Message:   req.Message,
	}
	if err := db.Create(&logEntry).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save message to database"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Message sent and saved", "data": req.Message})
}
