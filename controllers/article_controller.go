package controllers

import (
	"blog/config"
	"blog/models"
	"blog/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ArticleController struct {
	service services.ArticleService
}

func NewAricleController(service services.ArticleService) *ArticleController {
	return &ArticleController{service: service}
}

func (ac *ArticleController) GetAllArticles(c echo.Context) error {
	// 获取分页参数
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	pageNum, err := strconv.Atoi(c.QueryParam("pageNum"))
	keyword := c.QueryParam("keyword")
	if err != nil {
		articles, err := ac.service.GetAllArticles(0, 0, keyword)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error while fetching articles: %v", err))
		}
		return c.JSON(http.StatusOK, articles)
	}

	articles, err := ac.service.GetAllArticles(pageSize, pageNum, keyword)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error while fetching articles: %v", err))
	}
	response := config.Response{
		Success: true,
		Message: "请求成功",
		Data:    articles,
	}

	return c.JSON(http.StatusOK, response)
}

func (ac *ArticleController) GetArticleById(c echo.Context) error {
	id, _ := strconv.Atoi(c.QueryParam("id"))
	article, err := ac.service.GetArticleById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error while fetching article: %v", err))
	}
	return c.JSON(http.StatusOK, article)
}

func (ac *ArticleController) CreateArticle(c echo.Context) error {
	article := new(models.Article)
	if err := c.Bind(article); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Error while binding article: %v", err))
	}
	err := ac.service.CreateArticle(*article)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error while creating article: %v", err))
	}
	return c.JSON(http.StatusOK, article)
}

func (ac *ArticleController) GetTags(c echo.Context) error {
	tags, err := ac.service.GetTags()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error while fetching tags: %v", err))
	}
	// 创建一个 Response 对象
	response := config.Response{
		Success: true,
		Message: "请求成功",
		Data:    tags,
	}
	return c.JSON(http.StatusOK, response)
}

func (ac *ArticleController) CreateLink(c echo.Context) error {
	link := models.Link{}
	if err := c.Bind(&link); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error while binding link: %v", err))
	}
	err := ac.service.CreateLink(link)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error while creating link: %v", err))
	}
	return c.JSON(http.StatusOK, link)
}

func (ac *ArticleController) GetLinks(c echo.Context) error {
	links, err := ac.service.GetLinks()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error while fetching links: %v", err))
	}
	// 创建一个 Response 对象
	response := config.Response{
		Success: true,
		Message: "请求成功",
		Data:    links,
	}
	return c.JSON(http.StatusOK, response)
}
