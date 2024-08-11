package router

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"goBlog/controller"
	_ "goBlog/docs"
	"goBlog/handler"
	"goBlog/views/tmpl"
	"net/http"
	"time"
)

var e *echo.Echo

type jwtCustomClaims struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

func NewServer() {

	e = echo.New()
	e.Renderer = tmpl.Temp()
	e.Debug = true
	e.Use(middleware.Recover()) // 主要用于拦截panic错误并且在控制台打印错误日志，避免echo程序直接崩溃
	e.Use(middleware.Logger())  // Logger中间件主要用于打印http请求日志
	e.Use(handler.CountMiddleware)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			echo.GET, echo.POST, echo.PUT, echo.DELETE,
		},
	}))
	e.Static("/images", "uploads")
	e.Static("/static", "static")

	e.POST("/login", controller.Login)
	e.POST("/register", controller.Register)
	e.GET("/docs/*", echoSwagger.WrapHandler)
	m := e.Group("")
	ApiMovie(m)
	r := e.Group("usr")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &jwtCustomClaims{}
		},
		SigningKey: []byte("secret"),
	}

	r.Use(echojwt.WithConfig(config))
	ApiCreate(r)
	//r.GET("/profile", restricted)

	s := &http.Server{
		Addr:         ":1323",
		ReadTimeout:  50 * time.Minute,
		WriteTimeout: 50 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}
func Restricted(c echo.Context) error {
	// User
	type User struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*jwtCustomClaims)
	user := &User{
		Account:  claims.Account,
		Password: claims.Password,
	}
	return c.JSON(http.StatusOK, user)
}
