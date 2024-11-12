package services

import (
	"blog/models"
	"blog/repositories"
)

type ArticleService interface {
	GetAllArticles(pageSize int, pageNumber int) ([]models.Article, error)
	GetArticleById(id int) (models.Article, error)
	CreateArticle(article models.Article) error
	GetTags() ([]*models.Tag, error)
}

type articleService struct {
	articleRepo repositories.ArticleRepository
}

func NewArticleService(articleRepo *repositories.ArticleRepository) ArticleService {
	return &articleService{
		articleRepo: *articleRepo,
	}
}

func (as *articleService) GetAllArticles(pageSize int, pageNumber int) ([]models.Article, error) {
	return as.articleRepo.GetAllArticles(pageSize, pageNumber)
}

func (as *articleService) GetArticleById(id int) (models.Article, error) {
	return as.articleRepo.GetArticleById(id)
}

func (as *articleService) CreateArticle(article models.Article) error {
	return as.articleRepo.CreateArticle(article)
}

func (as *articleService) GetTags() ([]*models.Tag, error) {
	return as.articleRepo.GetTags()
}
