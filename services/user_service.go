package services

import (
	"blog/models"
	"blog/repositories"
)

type UserService interface {
	FindUser(username string, password string) (*models.User, error)
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	CreateVisitor(visitor *models.Visitor) (*models.Visitor, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) UserService {
	return &userService{userRepo: *userRepo}
}

func (us *userService) FindUser(username string, password string) (*models.User, error) {
	return us.userRepo.FindUser(username, password)
}

func (us *userService) CreateUser(user *models.User) error {
	return us.userRepo.CreateUser(user)
}

func (us *userService) GetUserByID(id int) (*models.User, error) {
	return us.userRepo.GetUserById(id)
}

func (us *userService) CreateVisitor(visitor *models.Visitor) (*models.Visitor, error) {
	return us.userRepo.CreateVisitor(visitor)
}

//func (us *userService) RefreshToken(c echo.Context) (token string, tokenRefresh string, err error) {
//	tokenData := c.Request().Header.Get("Authorization")
//	if tokenData == "" {
//		return "", "", echo.NewHTTPError(401, "Unauthorized")
//	}
//	// TODO: refresh token logic
//	tokenStr := strings.Split(tokenData, " ")[1]
//	_, claims, err := middlewares.ParseToken(tokenStr)
//	expireTime := jwt.NewNumericDate(time.Now().Add(10 * time.Second)) //token过期时间10秒，主要是测试方便
//	refreshTime := jwt.NewNumericDate(time.Now().Add(1 * time.Minute)) //刷新的时间限制，超过20秒重新登录
//	if err != nil {
//		return "", "", echo.NewHTTPError(500, err)
//	} else {
//		myClaims := models.CustomClaims{
//			claims.Username,
//			jwt.RegisteredClaims{
//				ExpiresAt: expireTime,
//			},
//		}
//		myClaimsRefrrsh := models.CustomClaims{
//			claims.Username,
//			jwt.RegisteredClaims{
//				ExpiresAt: refreshTime,
//			},
//		}
//		jwtKey := []byte("lyf123456")
//		tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
//		tokenStr, err := tokenObj.SignedString(jwtKey)
//		tokenFresh := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaimsRefrrsh)
//		tokenStrRefresh, err2 := tokenFresh.SignedString(jwtKey)
//		if err != nil && err2 != nil {
//			return "", "", echo.NewHTTPError(500, err)
//		} else {
//			return tokenStr, tokenStrRefresh, nil
//		}
//	}
//}
