package router

import (
	"github.com/labstack/echo/v4"
	"goBlog/controller"
)

func ApiCreate(api *echo.Group) {
	api.GET("/getArticle", controller.GetArticle)
	api.POST("/publish", controller.Publish)
	api.GET("/profile", Restricted)
}

func ApiMovie(api *echo.Group) {
	api.GET("", controller.Index)
	api.GET("/display", controller.Display)
	api.GET("/movie", controller.Movie)
	api.GET("/search", controller.Search)
	api.GET("/play", controller.Play)
	api.GET("/about", controller.About)
	api.GET("/spider", controller.GoSpider)
}
