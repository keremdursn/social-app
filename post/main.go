package main

import (
	"post/database"
	"post/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDb()

	app := fiber.New()
	router.Post(app)
	router.Comment(app)
	app.Listen(":8081")
}
