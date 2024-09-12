package router

import (
	"auth/controllers"

	"github.com/gofiber/fiber/v2"
)

func Post(app *fiber.App) {
	post := app.Group("/post")

	post.Post("/create-post", controllers.CreatePost)
	post.Put("/update-post/:id", controllers.UpdatePost)
	post.Delete("/delete/post/:id", controllers.DeletePost)
	post.Get("/getallpost", controllers.GetAllPost)
	post.Get("getpostbyid", controllers.GetPostByUserID)
	post.Post("/like-post", controllers.LikePost)
	post.Post("/get-back-like", controllers.GetBackLike)
}
