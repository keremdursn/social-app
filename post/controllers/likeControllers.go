package controllers

import (
	"post/database"
	"post/models"

	"github.com/gofiber/fiber/v2"
)

func LikePost(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	db := database.DB.Db

	var likePost models.LikePost
	var post models.Post

	// Postu bul getir
	id := c.Params("id")
	err := db.Where("id = ?",id).First(&post).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "fail", "message": "Post not found"})
	}

	postID := id
	userID := user.ID

	// Like'ın daha önce yapılıp yapılmadığını kontrol et
	var controlLike models.LikePost
	err = db.Where("post_id = ? AND user_id = ?", postID, userID).First(&controlLike).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "You have already liked this post"})
	}

	// Like'ı veritabanına kaydet
	likePost.UserID = user.ID
	likePost.PostID = post.PostID
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
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	db := database.DB.Db

	var post models.Post

	// Postu bul getir
	id := c.Params("id")
	err := db.Where("id = ?",id).First(&post).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "fail", "message": "Post not found"})
	}

	postID := id
	userID := user.ID

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
