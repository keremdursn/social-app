package controllers

import (
	"auth/database"
	"auth/helpers"
	"auth/middleware"
	"auth/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {

	db := database.DB.Db
	var user models.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid request"})
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

func LogIn(c *fiber.Ctx) error {
	db := database.DB.Db
	var login models.Login
	var user models.User

	//login verilerini aldık
	err := c.BodyParser(&login)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid request"})
	}

	//kullaniciyi db den bul
	err = db.Where("username = ?", login.Username).First(&user).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "failed", "message": "Username or password is wrong!"})
	}

	if err := helpers.CheckPass(user.Password, login.Password); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "failed", "message": "Username or password is wrong!"})
	}

	token, err := middleware.CreateToken(login.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Token oluşturulamadı")
	}

	// Token'in son kullanma tarihini belirle (örneğin, 72 saat geçerlilik süresi)
	expiry := time.Now().Add(72 * time.Hour).Unix()

	// IP adresini al
	ip := c.IP()

	// Oturumu veritabanına kaydet
	session := models.Session{
		UserID:   int(user.ID),
		Token:    token,
		Expiry:   expiry,
		IP:       ip,
		IsActive: true,
	}
	err = db.Create(&session).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session oluşturulamadı")
	} 
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Login success", "data": user, "token": token})
}
