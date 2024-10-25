package router

import (
	"auth/controllers"
	"auth/middleware"

	"github.com/gofiber/fiber/v2"
)

func User(app *fiber.App) {
	user := app.Group("/user")

	user.Post("/signup", controllers.SignUp)
	user.Post("/login", controllers.LogIn)
	user.Get("/logout", middleware.TokenControl(), controllers.LogOut)
	user.Put("/update-user", middleware.TokenControl(), controllers.UpdateAccount)
	user.Put("/update-photo", middleware.TokenControl(), controllers.UpdatePhoto)
	user.Put("/change-password", middleware.TokenControl(), controllers.ChangePassword)
	user.Delete("/delete-account/", middleware.TokenControl(), controllers.DeleteAccount)
	user.Get("/:id", middleware.TokenControl(), controllers.GetUserByID)
}
