package controllers

import (
	"auth/database"
	"auth/helpers"
	"auth/middleware"
	"auth/models"
	"io"
	"log"

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
	err = db.Where("mail = ?", login.Mail).First(&user).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "failed", "message": "Mail or password is wrong!"})
	}

	if err := helpers.CheckPass(user.Password, login.Password); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "failed", "message": "Username or password is wrong!"})
	}

	token, err := middleware.CreateToken(user.Mail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Token oluşturulamadı")
	}

	var session models.Session
	session.UserID = user.ID
	session.Token = token


	// Oturumu veritabanına kaydet
	// session := models.Session{
	// 	UserID:   user.ID,
	// 	Token:    token,
	// 	Expiry:   expiry,
	// 	IP:       ip,
	// 	IsActive: true,
	// }
	err = db.Create(&session).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session oluşturulamadı")
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Login success", "data": user, "token": token})
}

func LogOut(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}
	db := database.DB.Db
	var session models.Session

	// err = db.Where("user_id = ? ", user.ID).First(&session).Error
	// if err != nil {
	// 	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "bulamadi"})
	// }
	// session.IsActive = false
	log.Println("---------------buraya bak-------------------", session)
	// Session'ı veritabanından sil
	err := db.Raw("DELETE FROM sessions WHERE user_id= ?", user.ID).Scan(&session).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "silemedi"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Logout successful"})
}

func ChangePassword(c *fiber.Ctx) error {
	// Kullanıcıyı token ile doğrula
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	// Şifre değişiklik isteğini parse et
	changePassword := new(models.ChangePassword)
	if err := c.BodyParser(&changePassword); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid request"})
	}

	// Eski şifreyi doğrula
	err := helpers.CheckPass(user.Password, changePassword.OldPassword)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "faild", "message": "Old password is incorrect"})
	}

	// Yeni şifrelerin uyumunu kontrol et
	if changePassword.NewPassword1 != changePassword.NewPassword2 {
		return c.Status(401).JSON(fiber.Map{"status": "faild", "message": "New passwords do not match"})
	}

	// Şifrenin karmaşıklığını kontrol et (opsiyonel)
	if len(changePassword.NewPassword1) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Password must be at least 8 characters long"})
	}

	// Yeni şifreyi hashle
	hashedNewPassword := helpers.HashPass(changePassword.NewPassword1)

	// Şifreyi güncelle
	user.Password = hashedNewPassword
	db := database.DB.Db
	if err := db.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update password"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Password changed successfully"})
}

func DeleteAccount(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	db := database.DB.Db

	err := db.Delete(&user).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "failed to delete user", "data:": nil})
	}

	var session models.Session
	err = db.Where("user_id = ?", user.ID).Delete(&session).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "failed to delete session", "data:": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "account and session has been successfully deleted. "})
}

func GetUserByID(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}
	db := database.DB.Db

	id := c.Params("id")

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "user not found"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User Found", "data": user})
}

func UpdateAccount(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}
	log.Println("User ID:", user.ID)

	db := database.DB.Db

	//İstekten kullanıcı ID sini al
	// id := c.Params("id")

	// err := db.Where("id = ?", id).First(&user).Error
	// if err != nil {
	// 	return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "user not found"})
	// }

	// updateUserData := new(model.UpdateUser)
	// c.BodyParser(&updateUserData)
	// if err != nil {
	// 	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "json bodyparse edilemedi.", "data": err})
	// }

	// user.Username = updateUserData.Username
	// user.Name = updateUserData.Name
	// user.Surname = updateUserData.Surname

	name := c.FormValue("name")
	surname := c.FormValue("surname")
	username := c.FormValue("username")
	mail := c.FormValue("mail")

	if len(name) != 0 {
		user.Name = name
	}
	if len(surname) != 0 {
		user.Surname = surname
	}
	if len(username) != 0 {
		user.Username = username
	}
	if len(mail) != 0 {
		user.Mail = mail
	}

	err := db.Save(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "server error!"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "user data has been successfully updated "})
}

func UpdatePhoto(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": ""})
	}

	db := database.DB.Db

	//Kullanıcıdan profil fotoğrafı al
	file, err := c.FormFile("profile_image")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to upload file"})
	}

	// Dosyayı açma ve okuma
	openedFile, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to open file"})
	}
	defer openedFile.Close()

	imageBytes, err := io.ReadAll(openedFile)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to read file"})
	}

	// Cloudinary'e yükleme
	imageID, imageURL, err := database.ConnectToCloudinary(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to upload to Cloudinary"})
	}

	user.PPId = imageID
	user.PPUrl = imageURL

	err = db.Save(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update user profile photo"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "error", "message": "Photo uploaded successfully", "ppurl": imageURL})
}
