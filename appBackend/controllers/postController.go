package controllers

import (
	"auth/database"
	"auth/middleware"
	"auth/models"
	"io"

	"github.com/gofiber/fiber/v2"
)

//TODO: , update post, deletepost, getallpost, getuserspost

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
