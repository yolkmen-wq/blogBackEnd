package controller

import (
	"github.com/labstack/echo/v4"
	"goBlog/utils/spider"
)

func GoSpider(c echo.Context) error {
	spider.Create().Start()
	return nil

}
