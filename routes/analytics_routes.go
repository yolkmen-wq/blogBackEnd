package routes

import (
	"blog/config"
	"blog/controllers"
	"blog/repositories"
	"blog/services"
	"github.com/labstack/echo/v4"
)

func SetupAnalyticsRoutes(e *echo.Echo) {
	analyticController := controllers.NewAnalyticsController(services.NewAnalyticsService(repositories.NewAnalyticsRepository(config.RedisClient())))

	e.POST("/countVisits", analyticController.CountViews)
}
