package controllers

import (
	"blog/config"
	"blog/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AnalyticsController struct {
	analyticsService services.AnalyticsService
}

func NewAnalyticsController(analyticsService services.AnalyticsService) *AnalyticsController {
	return &AnalyticsController{
		analyticsService: analyticsService,
	}
}

func (ac *AnalyticsController) CountViews(c echo.Context) error {
	count, err := ac.analyticsService.CountViews()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// 创建一个 Response 对象
	response := config.Response{
		Success: true,
		Message: "请求成功",
		Data:    map[string]int64{"count": count},
	}
	return c.JSON(http.StatusOK, response)
}
