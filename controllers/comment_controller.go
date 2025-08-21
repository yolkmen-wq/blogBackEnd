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

type CommentController struct {
	commentService services.CommentService
}

func NewCommentController(commentService services.CommentService) *CommentController {
	return &CommentController{commentService}
}

func (c *CommentController) CreateComment(ctx echo.Context) error {
	comment := new(models.CommentChild)

	if err := ctx.Bind(comment); err != nil {
		fmt.Println(25, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid request payload")
	}
	if err := c.commentService.CreateComment(comment); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create comment")
	}
	response := config.Response{
		Success: true,
		Message: "评论成功",
		Data:    "",
	}
	return ctx.JSON(http.StatusOK, response)
}

func (c *CommentController) GetCommentByArticleID(ctx echo.Context) error {
	articleId, err := strconv.ParseInt(ctx.QueryParam("articleId"), 10, 64)

	var visitorID int64
	if ctx.QueryParam("visitorID") != "" {
		visitorID, err = strconv.ParseInt(ctx.QueryParam("visitorID"), 10, 64)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get comments")
	}
	comments, err := c.commentService.GetCommentByArticleID(articleId, visitorID)
	if err != nil {
		fmt.Println(51, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get comments")
	}
	response := config.Response{
		Success: true,
		Message: "获取评论成功",
		Data:    comments,
	}
	return ctx.JSON(http.StatusOK, response)
}

func (c *CommentController) ToggleCommentLike(ctx echo.Context) error {
	//commentId, err := strconv.ParseInt(ctx.QueryParam("commentId"), 10, 64)
	like := new(models.Like)

	if err := ctx.Bind(like); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid request payload")
	}

	if err := c.commentService.ToggleCommentLike(like); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to toggle like")
	}
	response := config.Response{
		Success: true,
		Message: "点赞成功",
		Data:    "",
	}
	return ctx.JSON(http.StatusOK, response)
}
