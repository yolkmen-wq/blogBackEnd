package controller

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"goBlog/models"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// @Tags 用户
// @Description 登录
// @Accept  json
// @Param account formData string true "账号"
// @Param password formData string true "密码"
// @Success 200 "获取信息成功"
// @Failure 400	"获取信息失败"
// @Router /login [post]
type jwtCustomClaims struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}
type User struct {
	Account  string `form:"account"`
	Password string `form:"password"`
}

func Login(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(1111, user)
	data, err := models.QueryUser(user.Account, user.Password)
	fmt.Println(data)
	if err != nil {
		log.Println("QueryUser error:", err)
	}
	if data != nil {
		claims := &jwtCustomClaims{
			Account:  user.Account,
			Password: user.Password,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code": 500,
		"msg":  "该用户未注册！",
	})
}

// @Tags 用户
// @Description 注册
// @Accept json
// @Param account formData string true "账号"
// @Param password formData string true "密码"
// @Success 200 "获取信息成功"
// @Failure 400	"获取信息失败"
// @Router /register [post]
func Register(c echo.Context) error {
	acc := c.FormValue("account")
	pwd := c.FormValue("password")
	err := models.InsertUser(acc, pwd)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"msg": "注册失败！",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"msg": "注册成功！",
	})
}

// @Tags 用户
// @Title GetProfile
// @Description 获取用户信息
// @Success 200 "获取信息成功"
// @Failure 400	"获取信息失败"
// Router /profile [get]
func GetProfile(c echo.Context) error {
	fmt.Println(c.Get("user"))
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return errors.New("failed to cast claims as jwt.MapClaims")
	}
	fmt.Println(user.Claims)
	claims := user.Claims.(*jwtCustomClaims)

	return c.JSON(http.StatusOK, claims)
}

// @Tags 文章
// @Title Publish
// @Description 发布
// @Accept  json
// @Param artic formData string true "文章内容"
// @Success 200 "获取信息成功"
// @Failure 400	"获取信息失败"
// @Router /publish [post]
//func Publish(c echo.Context) error {
//	//// 从请求参数里获取 team 和 member 的值
//	//team := c.QueryParam("team")
//	//member := c.QueryParam("member")
//	article := c.FormValue("article")
//	file, err := os.Create("test.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()
//	for i := 0; i < 5; i++ {
//		file.WriteString("ab\n")
//		file.Write([]byte("cd\n"))
//	}
//	return c.String(http.StatusOK, article)
//}

// @Tags 文章
// @Title getArticle
// @Description 获取文章
// @Success 200 "获取信息成功"
// @Failure 400	"获取信息失败"
// @Router /getArticle [get]
func GetArticle(c echo.Context) error {
	result, err := models.Getarticle()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "获取文章失败！",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": result,
	})
}

// @Tags 文章
// @Title save
// @Description 发布文章
// @Success 200 "获取信息成功"
// @Failure 400	"获取信息失败"
// @Router /save [post]
type Article struct {
	title       string
	description string
	image       string
	imageText   string
}

func (a Article) GetTitle() string {
	return a.title
}
func (a Article) GetDescription() string {
	return a.description
}
func (a Article) GetImage() string {
	return a.image
}
func (a Article) GetImageText() string {
	return a.imageText
}
func Publish(c echo.Context) error {

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Failed to get image",
		})
	}
	//打开文件
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to open image",
		})
	}
	defer src.Close()

	//读取文件数据
	fileData, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to read image",
		})
	}
	//保存文件到文件系统
	filePath := filepath.Join("uploads", file.Filename)
	err = os.WriteFile(filePath, fileData, 0644)
	if err != nil {
		fmt.Printf("错误：%s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to save image",
		})
	}
	art := Article{
		title:       c.FormValue("title"),
		description: c.FormValue("description"),
		image:       filePath,
		imageText:   file.Filename,
	}

	err = models.Addarticle(art)
	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"msg": "新增成功！",
	})
}
