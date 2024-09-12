package main

import (
	"auth/database"
	"auth/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDb()

	app := fiber.New()
	router.User(app)
	router.Post(app)
	router.Comment(app)
	app.Listen(":8080")
}
