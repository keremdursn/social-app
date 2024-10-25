package controllers

import (
	"post/database"
	"post/middleware"
	"post/models"

	"github.com/gofiber/fiber/v2"
)

func Comment(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	db := database.DB.Db
	//Yorum verilerini al
	var comment models.Comment
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid request"})
	}
	//Postu bul getir
	var post models.Post
	err = db.Where("id = ?", comment.PostID).First(&post).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Post not found"})
	}
	//Postun yorum sayısını güncelle
	post.CommentCount++
	err = db.Save(&post).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update post comment count"})
	}
	//Commenti kaydet
	comment.UserID = user.ID
	err = db.Create(&comment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create comment"})
	}
	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Comment created successfully", "data": comment})
}

func DeleteComment(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}
	db := database.DB.Db

	commentID := c.Params("comment_id")

	var comment models.Comment

	err = database.DB.Db.Where("id = ? AND user_id = ?", commentID, user.ID).First(&comment).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Comment not found or not authorized to delete"})
	}

	err = db.Delete(&comment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete comment"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Comment archived successfully"})
}

func LikeCommand(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}
	db := database.DB.Db

	// Bodyden post verilerini al
	var likeComment models.LikeComment
	if err := c.BodyParser(&likeComment); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	commentID := likeComment.CommentID
	userID := user.ID

	// Commenti bul getir
	var comment models.Comment
	err = db.First(&comment, commentID).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "fail", "message": "Comment not found"})
	}

	// Like'ın daha önce yapılıp yapılmadığını kontrol et
	var commentControl models.LikeComment
	err = db.Where("comment_id = ? AND user_id = ?", commentID, userID).First(&commentControl).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "You have already liked this comment"})
	}

	// Like'ı veritabanına kaydet
	likeComment.UserID = user.ID
	likeComment.CommentID = comment.CommentID
	err = db.Create(&likeComment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to like comment"})
	}

	// Commentin Like sayısını güncelle
	comment.LikeCount++
	err = db.Save(&comment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update comment like count"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Comment liked successfully", "commentlike": likeComment})

}

func GetBackLikeCommand(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}
	db := database.DB.Db

	var likeComment models.LikeComment
	if err := c.BodyParser(&likeComment); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	commentID := likeComment.CommentID
	userID := user.ID

	//Commenti bul
	var comment models.Comment
	err = db.Where("comment_id = ? ", commentID).First(&comment).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "fail", "message": "Comment not found"})
	}

	//Like controlü yap
	var controlLike models.LikeComment
	err = db.Where("comment_id = ? AND user_id = ?", commentID, userID).First(&controlLike).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "fail", "message": "Comment like not found"})
	}

	//Like ı sil
	err = db.Delete(&controlLike).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to remove comment like"})
	}

	//Like sayısını güncelle
	comment.LikeCount--
	err = db.Save(&comment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update comment like count"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Like removed successfully"})
}

func AnswerComment(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}
	db := database.DB.Db

	//body den answer isteğini al
	var answerComment models.AnswerComment
	if err := c.BodyParser(&answerComment); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	//yorumu bul
	var comment models.Comment
	err = db.Where("comment_id = ? ", answerComment.CommentID).First(&comment).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Comment not found"})
	}

	//answer structını doldur
	answerComment.UserID = user.ID
	answerComment.CommentID = comment.CommentID
	//answerı kaydet
	err = db.Create(&answerComment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not save the answer"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Answer added successfully"})
}

func GetAllCommentsByPostID(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}
	db := database.DB.Db

	// Post ID'yi route parametresinden alıyoruz
	postID := c.Params("post_id")

	// Belirtilen post ID'ye göre tüm yorumları veritabanından alıyoruz
	var comments []models.Comment
	err = db.Where("post_id = ?", postID).Preload("User").Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to fetch comments"})
	}

	// Eğer hiç yorum yoksa
	if len(comments) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No comments found for this post"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Comments fetched successfully", "data": comments, "user": user})
}

// delete answer eklenecek
func DeleteAnswer(c *fiber.Ctx) error {
	user, err := middleware.TokenControl(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}
	db := database.DB.Db

	answerID := c.Params("answer_id")

	var answer models.AnswerComment

	err = db.Where("answer_id = ? AND user_id = ?", answerID, user.ID).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Answer not found or not authorized to delete"})
	}

	err = db.Delete(&answer).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to delete answer"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Answer archived successfully"})
}
