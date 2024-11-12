package controllers

import (
	"blog/models"
	"blog/services"
	"fmt"
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
		fmt.Println(23, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid request data"})
	}

	err := uc.userService.CreateUser(user)
	if err != nil {
		fmt.Println(29, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, user)
}

func (uc *UserController) GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}
