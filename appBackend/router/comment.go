package router

import (
	"auth/controllers"

	"github.com/gofiber/fiber/v2"
)

func Comment(post *fiber.App) {
	comment := post.Group("/comment")

	comment.Post("/comment", controllers.Comment)
	comment.Delete("/delete-comment", controllers.DeleteComment)
	comment.Post("/like-comment", controllers.LikeCommand)
	comment.Post("/get-back-likecomment", controllers.GetBackLikeCommand)
	comment.Post("/answer-comment", controllers.AnswerComment)
	comment.Delete("/delete-answer", controllers.DeleteAnswer)
	comment.Post("/get-all-comment/:id", controllers.GetAllCommentsByPostID)
}
