package controllers

import (
	"auth/database"
	"auth/helpers"
	"auth/models"

	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {

	db := database.DB.Db
	var user models.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid request"})
	}

	//username kontrol
	err = helpers.UsernameControl(user.Username)
	if err == nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "username already taken", "data": err})
	}

	//mail kontrol
	err = helpers.MailControl(user.Mail)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "mail used by another one", "data": err})
	}

	//pasHass
	user.Password = helpers.HashPass(user.Password)

	err = db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "user couldn't created"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "user created", "data": user})
}
