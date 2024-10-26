package router

import (
	"post/controllers"
	"post/middleware"

	"github.com/gofiber/fiber/v2"
)

func Comment(post *fiber.App) {
	comment := post.Group("/comment")

	comment.Post("/comment", middleware.TokenControl(), controllers.Comment)
	comment.Delete("/delete-comment", middleware.TokenControl(), controllers.DeleteComment)
	comment.Post("/like-comment", middleware.TokenControl(), controllers.LikeCommand)
	comment.Post("/get-back-likecomment", middleware.TokenControl(), controllers.GetBackLikeCommand)
	comment.Post("/answer-comment", middleware.TokenControl(), controllers.AnswerComment)
	comment.Delete("/delete-answer", middleware.TokenControl(), controllers.DeleteAnswer)
	comment.Post("/get-all-comment/:id", middleware.TokenControl(), controllers.GetAllCommentsByPostID)
}
