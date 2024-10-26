package router

import (
	"post/controllers"
	"post/middleware"

	"github.com/gofiber/fiber/v2"
)

func Post(app *fiber.App) {
	post := app.Group("/post")

	post.Post("/create-post", middleware.TokenControl(), controllers.CreatePost)
	post.Put("/update-post/:id", middleware.TokenControl(), controllers.UpdatePost)
	post.Delete("/delete/post/:id", middleware.TokenControl(), controllers.DeletePost)
	post.Get("/getallpost", controllers.GetAllPost)
	post.Get("getpostbyid", middleware.TokenControl(), controllers.GetPostByUserID)
	post.Post("/like-post", middleware.TokenControl(), controllers.LikePost)
	post.Post("/get-back-like", middleware.TokenControl(), controllers.GetBackLike)
}
