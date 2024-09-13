package router

import (
	"notification/controllers"

	"github.com/gofiber/fiber/v2"
)

func Notification(app *fiber.App) {
	msg := app.Group("/notification")

	msg.Post("/send-message", controllers.SendMessage)
}
