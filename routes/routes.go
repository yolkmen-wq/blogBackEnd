package routes

import "github.com/labstack/echo/v4"

func InitRoutes(e *echo.Echo) {
	SetupUserRoutes(e)
	SetupArticleRoutes(e)
	SetupAnalyticsRoutes(e)
}
