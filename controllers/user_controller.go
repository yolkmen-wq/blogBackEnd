package controllers

import (
	"blog/config"
	"blog/models"
	"blog/services"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mojocn/base64Captcha"
	"net/http"
	"strconv"
	"time"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

var store = base64Captcha.DefaultMemStore
var client = config.RedisClient()

// 创建验证码
func CreateCaptcha() (string, string, error) {
	var param = models.ConfigJsonBody{
		CaptchaType: "string",
		VerifyValue: "",
		DriverAudio: &base64Captcha.DriverAudio{},
		DriverString: &base64Captcha.DriverString{
			Height:          60,
			Width:           240,
			Length:          5,
			NoiseCount:      20,
			ShowLineOptions: 3,
			Fonts:           []string{"wqy-microhei.ttc"},
			Source:          config.CaptchaSource,
		},
		DriverMath: &base64Captcha.DriverMath{},
		DriverChinese: &base64Captcha.DriverChinese{
			Height:          60,
			Width:           240,
			Length:          6,
			NoiseCount:      0,
			ShowLineOptions: 2,
			Fonts:           []string{"wqy-microhei.ttc"},
			Source:          "你好，世界！",
		},
		DriverDigit: &base64Captcha.DriverDigit{},
	}
	var driver base64Captcha.Driver

	//create base64 encoding captcha
	switch param.CaptchaType {
	case "audio":
		driver = param.DriverAudio
	case "string":
		driver = param.DriverString.ConvertFonts()
	case "math":
		driver = param.DriverMath.ConvertFonts()
	case "chinese":
		driver = param.DriverChinese.ConvertFonts()
	default:
		driver = param.DriverDigit
	}
	//store := client.Set("captcha_store", "123456", 10*60)
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := captcha.Generate()
	if err != nil {
		return "", "", err

	}
	// 获取验证码的答案
	verifyValue := store.Get(id, true)

	client.Set("captcha_store"+id, verifyValue, time.Minute*2)
	return id, b64s, err

}

// 校验验证码
func CaptchaVerifyHandle(code string, id string) error {
	verifyValue := client.Get("captcha_store" + id)
	if verifyValue.Val() == "" {
		return errors.New("验证码已过期")
	}
	if verifyValue.Val() != code {
		return errors.New("验证码错误")
	}
	return nil
}

func (uc *UserController) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	code := c.FormValue("code")
	captchaId := c.FormValue("captchaId")
	expireTime := jwt.NewNumericDate(time.Now().Add(2 * time.Hour)) //token过期时间10秒，主要是测试方便
	user, err := uc.userService.FindUser(username, password)
	if err != nil {
		response := config.Response{
			Success: false,
			Message: err.Error(),
			Data:    "",
		}
		return c.JSON(500, response)
	}
	if user == nil {
		response := config.Response{
			Success: false,
			Message: "用户名或密码错误",
			Data:    "",
		}
		return c.JSON(200, response)
	}

	err = CaptchaVerifyHandle(code, captchaId)
	if err != nil {
		response := config.Response{
			Success: false,
			Message: err.Error(),
			Data:    "",
		}
		return c.JSON(500, response)
	}
	claims := models.CustomClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: expireTime,                     // 24小时后过期
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 设置签发时间
			NotBefore: jwt.NewNumericDate(time.Now()), // 设置生效时间
			Issuer:    "test",                         // 设置签发人
			Subject:   "somebody",                     // 设置主题
			ID:        "1",                            // 设置ID
			Audience:  []string{"somebody_else"},      // 设置受众
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(config.SecretKey)
	if err != nil {
		response := config.Response{
			Success: false,
			Message: "登录失败",
			Data:    "",
		}
		return c.JSON(200, response)
	} else {
		response := config.Response{
			Success: true,
			Message: "登录成功",
			Data: echo.Map{
				"userInfo": user,
				"token":    tokenStr, //数据请求的token
			},
		}
		return c.JSON(200, response)
	}
}

func (uc *UserController) CreateUser(c echo.Context) error {
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid request data"})
	}
	err := uc.userService.CreateUser(user)

	if err != nil {
		Response := config.Response{
			Success: false,
			Message: err.Error(),
			Data:    "",
		}
		return c.JSON(http.StatusInternalServerError, Response)
	}
	response := config.Response{
		Success: true,
		Message: "User created successfully",
		Data:    "",
	}
	return c.JSON(http.StatusOK, response)
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

func (uc *UserController) CreateCaptcha(c echo.Context) error {
	id, b64s, err := CreateCaptcha()

	response := config.Response{
		Success: true,
		Message: "Captcha created successfully",
		Data: echo.Map{
			"id":   id,
			"b64s": b64s,
		},
	}
	if err != nil {
		response := config.Response{
			Success: true,
			Message: err.Error(),
			Data:    "",
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, response)
}
