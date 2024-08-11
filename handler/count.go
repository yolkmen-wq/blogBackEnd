package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

var totalVisits int64 = 0

func CountMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		totalVisits++
		c.Response().Header().Add("requests", fmt.Sprintf("%d", totalVisits))
		return next(c)
	}
}
