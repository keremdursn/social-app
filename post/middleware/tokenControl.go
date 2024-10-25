package middleware

import (
	"post/database"
	"post/models"

	"github.com/gofiber/fiber/v2"
)

func TokenControl() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.DB.Db
		authorizationHeader := c.Get("Authorization")

		if authorizationHeader == "" || len(authorizationHeader) < 7 || authorizationHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No token provided"})
		}
		token := authorizationHeader[7:]

		var session models.Session
		if err := db.Where("token = ?", token).First(&session).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token"})
		}

		var user models.User
		if err := db.Where("id = ?", session.UserID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		c.Locals("user", user)

		return c.Next()

	}

}