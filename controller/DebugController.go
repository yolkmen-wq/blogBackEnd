package controller

import (
	"github.com/labstack/echo/v4"
	"goBlog/utils/spider/tian_kong"
)

func Debug(c echo.Context) {

	go tian_kong.DelAllListCacheKey()

}
