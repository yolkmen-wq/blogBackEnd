package routes

import (
	"blog/config"
	"blog/controllers"
	"blog/repositories"
	"blog/services"
	"github.com/labstack/echo/v4"
)

func setupCommentRoutes(e *echo.Echo) {
	controller := controllers.NewCommentController(services.NewCommentService(repositories.NewCommentRepository(config.DB())))
	e.POST("/createComment", controller.CreateComment)
	e.GET("/getComments", controller.GetCommentByArticleID)
	e.POST("/toggleCommentLike", controller.ToggleCommentLike)
}
