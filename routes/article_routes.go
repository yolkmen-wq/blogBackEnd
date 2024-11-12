package routes

import (
	"blog/config"
	"blog/controllers"
	"blog/repositories"
	"blog/services"
	"github.com/labstack/echo/v4"
)

func SetupArticleRoutes(e *echo.Echo) {
	articleController := controllers.NewAricleController(services.NewArticleService(repositories.NewArticleRepository(config.DB())))

	e.GET("/articles", articleController.GetAllArticles)
	e.GET("/getArticleById", articleController.GetArticleById)
	e.POST("/createArticle", articleController.CreateArticle)
	e.POST("/getTags", articleController.GetTags)
}
