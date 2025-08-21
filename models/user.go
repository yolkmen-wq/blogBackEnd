package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mojocn/base64Captcha"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Email    string `json:"email"  form:"email"`
}

type Visitor struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	IP       string `json:"ip"`
}

type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// configJsonBody json request body.
type ConfigJsonBody struct {
	Id            string
	CaptchaType   string `json:"captchaType"`
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}
