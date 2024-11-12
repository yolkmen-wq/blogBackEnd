package main

import (
	"blog/config"
	"blog/routes"
	_ "blog/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

func main() {
	config.InitConfig()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			echo.GET, echo.POST, echo.PUT, echo.DELETE,
		},
	}))
	routes.InitRoutes(e)

	s := &http.Server{
		Addr:         ":1323",
		ReadTimeout:  50 * time.Minute,
		WriteTimeout: 50 * time.Minute,
	}

	e.StartServer(s)

	e.Logger.Fatal(e.StartServer(s))
}
