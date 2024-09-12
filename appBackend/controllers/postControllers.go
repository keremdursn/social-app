package controllers

import (
	"auth/database"
	"auth/middleware"
	"auth/models"
	"io"

	"github.com/gofiber/fiber/v2"
)

//TODO: , getallpost, getuserspost

func CreatePost(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	db := database.DB.Db
	var Post models.Post

	// Çoklu resimleri al
	files, err := c.MultipartForm()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to get form data"})
	}

	fileHeaders := files.File["post_images"]
	if len(fileHeaders) == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "No files uploaded"})
	}

	var imageIDs []string
	var imageURLs []string

	// Her bir resmi yükleyin
	for _, file := range fileHeaders {
		openedFile, err := file.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to open file"})
		}
		defer openedFile.Close()

		imageBytes, err := io.ReadAll(openedFile)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to read file"})
		}

		imageID, imageURL, err := database.ConnectToCloudinary(imageBytes)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to upload to Cloudinary"})
		}

		imageIDs = append(imageIDs, imageID)
		imageURLs = append(imageURLs, imageURL)
	}

	// Post nesnesine resim URL'lerini ekle
	Post.ImageIDs = imageIDs
	Post.ImageURLs = imageURLs

	description := c.FormValue("description")
	Post.PostDesc = description
	Post.UserID = user.ID

	err = db.Create(&Post).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Post not created", "data": err})
	}
	return c.Status(201).JSON(fiber.Map{"Status": "Success", "Message": "Post created successfully", "data": Post})
}

func UpdatePost(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}
	db := database.DB.Db

	postID := c.Params("id")
	var post models.Post

	// Mevcut post'u veritabanından al
	if err := database.DB.Db.First(&post, postID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"Status": "Error", "Message": "Post not found"})
	}

	if post.UserID != user.ID {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "you are not allowed to update this post."})
	}

	oldImageIDs := post.ImageIDs

	// Eski resimleri sil
	if len(oldImageIDs) > 0 {
		err = database.DeleteFromCloudinary(oldImageIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete old images from Cloudinary"})
		}
	}

	// Resimleri güncelle
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to parse form"})
	}

	files := form.File["post_images"]
	if len(files) > 0 {
		var imageIDs []string
		var imageURLs []string

		for _, file := range files {
			openedFile, err := file.Open()
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to open file"})
			}
			defer openedFile.Close()

			imageBytes, err := io.ReadAll(openedFile)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to read file"})
			}

			imageID, imageURL, err := database.ConnectToCloudinary(imageBytes)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to upload to Cloudinary"})
			}

			imageIDs = append(imageIDs, imageID)
			imageURLs = append(imageURLs, imageURL)
		}

		post.ImageIDs = imageIDs
		post.ImageURLs = imageURLs
	}

	description := c.FormValue("description")
	post.PostDesc = description

	err = db.Save(&post).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update post", "data": err})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Post updated successfully", "data": post})
}

func DeletePost(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}
	db := database.DB.Db

	id := c.Params("id")
	var post models.Post

	err = db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"Status": "Error", "Message": "Post not found"})
	}

	//Kullanıcı Postun Sahibi mi kontrolü
	if post.UserID != user.ID {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not allowed to delete this post."})
	}

	post.IsActive = false

	err = db.Save(&post).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update post"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Post deleted successfully", "data": post})
}

func GetAllPost(c *fiber.Ctx) error {
	db := database.DB.Db
	var posts []models.Post
	err := db.Preload("User").Where("is_active = ?", true).Order("id DESC").Find(&posts).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to get posts"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Posts retrieved successfully", "data": posts})
}

func GetPostByUserID(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	db := database.DB.Db
	var posts []models.Post

	err = db.Preload("User").Where("user_id = ? AND is_active = ?", user.ID, true).Order("id DESC").Find(&posts).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to get posts"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Posts retrieved successfully", "data": posts})
}
