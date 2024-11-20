package controllers

import (
	"blog/config"
	"blog/models"
	"blog/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid request data"})
	}

	err := uc.userService.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (uc *UserController) CreateVisitor(c echo.Context) error {
	visitor := new(models.Visitor)
	if err := c.Bind(visitor); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid request data"})
	}

	result, err := uc.userService.CreateVisitor(visitor)
	if err != nil {
		response := config.Response{
			Success: false,
			Message: err.Error(),
			Data:    "",
		}
		return c.JSON(http.StatusOK, response)
	}
	response := config.Response{
		Success: true,
		Message: "Visitor created successfully",
		Data:    result,
	}

	return c.JSON(http.StatusOK, response)
}
