package services

import (
	"blog/models"
	"blog/repositories"
)

type CommentService interface {
	CreateComment(comment *models.CommentChild) error
	GetCommentByArticleID(articleID int64, visitorID int64) ([]*models.Comment, error)
	ToggleCommentLike(like *models.Like) error
}

type commentService struct {
	commentRepository *repositories.CommentRepository
}

func NewCommentService(commentRepository *repositories.CommentRepository) CommentService {
	return &commentService{commentRepository}
}

func (s *commentService) CreateComment(comment *models.CommentChild) error {
	return s.commentRepository.CreateComment(comment)
}

func (s *commentService) GetCommentByArticleID(articleID int64, visitorID int64) ([]*models.Comment, error) {
	return s.commentRepository.GetCommentByArticleID(articleID, visitorID)
}

func (s *commentService) ToggleCommentLike(like *models.Like) error {
	return s.commentRepository.ToggleCommentLike(like)
}
