package middleware

import (
	"auth/database"
	"auth/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TokenControl(c *fiber.Ctx) (models.User, error) {
	db := database.DB.Db
	authorizationHeader := c.Get("Authorization")

	if authorizationHeader == "" || len(authorizationHeader) < 7 || authorizationHeader[:7] != "Bearer " {
		return models.User{}, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No token provided"})
	}
	token := authorizationHeader[7:]

	var session models.Session
	if err := db.Where("token = ? AND is_active = ?", token, true).First(&session).Error; err != nil {
		return models.User{}, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token"})
	}

	var user models.User
	if err := db.Where("id = ?", session.UserID).First(&user); err != nil {
		return models.User{}, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Token süresini kontrol et
	if time.Now().Unix() > session.Expiry {
		return models.User{}, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token expired"})
	}

	// IP adresini kontrol et (opsiyonel, IP'yi kaydettiyseniz)
	if session.IP != c.IP() {
		return models.User{}, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "IP address mismatch"})
	}

	return user, nil

}
