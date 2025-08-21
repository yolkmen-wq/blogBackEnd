package routes

import (
	"blog/config"
	"blog/controllers"
	"blog/repositories"
	"blog/services"
	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(e *echo.Echo) {
	userController := controllers.NewUserController(services.NewUserService(repositories.NewUserRepository(config.DB())))

	e.POST("/users", userController.CreateUser)
	e.GET("/users/:id", userController.GetUserByID)
	e.POST("/registerVisitor", userController.CreateVisitor)
	e.POST("/login", userController.Login)
	e.POST("/register", userController.CreateUser)
	e.POST("/getCaptcha", userController.CreateCaptcha)
}
