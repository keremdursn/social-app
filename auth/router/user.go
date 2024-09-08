package router

import (
	"auth/controllers"

	"github.com/gofiber/fiber/v2"
)

func User(app *fiber.App) {
	user := app.Group("/user")

	user.Post("/signup", controllers.SignUp)
	user.Post("/login", controllers.LogIn)
	user.Get("/logout", controllers.LogOut)
}
