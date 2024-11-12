package repositories

import (
	"blog/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ArticleRepository struct {
	db *sqlx.DB
}

func NewArticleRepository(db *sqlx.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (ar *ArticleRepository) GetAllArticles(pageSize int, pageNumber int) ([]models.Article, error) {
	var articles []models.Article
	var query string

	// 获取总记录数
	var totalCount int
	err := ar.db.Get(&totalCount, "SELECT COUNT(*) FROM articles")
	if err != nil {
		return nil, err
	}
	// 计算偏移量
	offset := pageNumber * pageSize

	// 计算分页
	if (pageSize == 0 && pageNumber == 0) || offset >= totalCount {
		query = "SELECT * FROM articles"
		err := ar.db.Select(&articles, query)
		if err != nil {
			return nil, err
		}
		return articles, nil
	} else {
		// 查询指定数量的记录
		query = "SELECT * FROM articles LIMIT ?, ?"
		// 执行查询
		err := ar.db.Select(&articles, query, offset, pageSize)
		if err != nil {
			return nil, err
		}

		return articles, nil
	}

}

func (ar *ArticleRepository) GetArticleById(id int) (models.Article, error) {
	var article models.Article
	err := ar.db.Get(&article, "SELECT * FROM articles WHERE id = ?", id)
	if err != nil {
		return models.Article{}, err
	}
	return article, nil
}

func (ar *ArticleRepository) CreateArticle(article models.Article) error {
	result, err := ar.db.NamedExec("INSERT INTO articles (title, brief, content) VALUES (:title, :brief,  :content)", article)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	for _, tag := range article.Tags {
		_, err := ar.db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", id, tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ar *ArticleRepository) GetTags() ([]*models.Tag, error) {
	var tags []*models.Tag
	err := ar.db.Select(&tags, "SELECT  * FROM tags")
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		var count int
		err = ar.db.Get(&count, "SELECT COUNT(*) FROM article_tags where tag_id = ?", tag.TagID)
		tag.Count = count
		if err != nil {
			return nil, err
		}
	}
	return tags, nil
}
