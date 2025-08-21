package middlewares

import (
	"blog/config"
	"blog/models"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"time"
)

// 解析和验证 Token
func validateToken(tokenString string) (*models.CustomClaims, error) {
	claims := &models.CustomClaims{}

	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.SecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("无效的 Token")
	}

	return claims, nil
}

// 延长 Token 有效期的函数
func extendToken(tokenString string, newExpiryDuration *jwt.NumericDate) (string, error) {
	// 验证当前 Token
	claims, err := validateToken(tokenString)
	if err != nil {
		return "", err // Token 无效
	}

	// 生成新的 Token
	return createToken(claims, newExpiryDuration)
}

// 生成新的 Token
func createToken(claims *models.CustomClaims, expiresIn *jwt.NumericDate) (string, error) {

	// 创建一个新的 Token
	newClaims := &models.CustomClaims{
		claims.Username,
		jwt.RegisteredClaims{
			ExpiresAt: expiresIn,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	// 对 Token 进行签名
	return token.SignedString(config.SecretKey)
}

func AuthMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		for _, route := range config.ExcludeRoutes {
			if c.Path() == route {
				return next(c)
			}
		}
		response := config.Response{
			Success: false,
			Message: "Unauthorized",
			Data:    "",
		}
		// 验证token
		if c.Request().Header.Get("Authorization") == "" {

			return c.JSON(200, response)
		}
		tokenData := c.Request().Header.Get("Authorization")
		claims, err := validateToken(tokenData)
		now := time.Now().Unix()
		expiry := claims.RegisteredClaims.ExpiresAt.Unix()
		if expiry-now < 1800 {
			// 延长token有效期
			newExpiryDuration := jwt.NewNumericDate(time.Now().Add(time.Hour * 2))
			newToken, err := extendToken(tokenData, newExpiryDuration)
			if err != nil {
				return c.JSON(401, response)
			}

			c.Response().Header().Set("Authorization", "Bearer "+newToken)
			claims, err = validateToken(newToken)
		}
		if err != nil {
			return c.JSON(401, response)
		}
		//// 验证权限
		//if claims.Role!= "admin" {
		//	return c.JSON(403, "Forbidden")
		//}
		return next(c)
	}

}
