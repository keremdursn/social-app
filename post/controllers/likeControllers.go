package controllers

import (
	"post/database"
	"post/middleware"
	"post/models"

	"github.com/gofiber/fiber/v2"
)

func LikePost(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	db := database.DB.Db

	// Bodyden post verilerini al
	var likePost models.LikePost
	if err := c.BodyParser(&likePost); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	postID := likePost.PostID
	userID := user.ID

	// Postu bul getir
	var post models.Post
	err = db.First(&post, postID).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "fail", "message": "Post not found"})
	}

	// Like'ın daha önce yapılıp yapılmadığını kontrol et
	var controlLike models.LikePost
	err = db.Where("post_id = ? AND user_id = ?", postID, userID).First(&controlLike).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "You have already liked this post"})
	}

	// Like'ı veritabanına kaydet
	likePost.UserID = user.ID
	err = db.Create(&likePost).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to like post"})
	}

	// Postun Like sayısını güncelle
	post.LikeCount++
	err = db.Save(&post).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update post like count"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Post liked successfully", "like": likePost})
}

func GetBackLike(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	db := database.DB.Db

	var like models.LikePost
	if err := c.BodyParser(&like); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	postID := like.PostID
	userID := user.ID

	//Postu bul
	var post models.Post
	err = db.Where("id = ?", postID).First(&post).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "fail", "message": "Post not found"})
	}

	//Like kontrolü yap
	var controlLike models.LikePost
	err = db.Where("post_id = ? AND user_id = ?", postID, userID).First(&controlLike).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "fail", "message": "Post like not found"})
	}

	//Beğeniyi sil
	err = db.Delete(&controlLike).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to remove post like"})
	}

	//Postun beğeni sayısını güncelle
	post.LikeCount--
	err = db.Save(&post).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update post like count"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Like removed successfully"})
}
