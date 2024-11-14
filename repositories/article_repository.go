package repositories

import (
	"blog/models"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ArticleRepository struct {
	db *sqlx.DB
}

func NewArticleRepository(db *sqlx.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// GetArticleById 根据 ID 获取文章
func (ar *ArticleRepository) GetArticleById(id int) (models.Article, error) {
	var article models.Article
	err := ar.db.Get(&article, "SELECT * FROM articles WHERE id = ?", id)
	if err != nil {
		return models.Article{}, err
	}
	return article, nil
}

// GetAllArticles 获取文章列表
func (ar *ArticleRepository) GetAllArticles(pageSize int, pageNumber int, keyword string) ([]models.Article, error) {
	var articles []models.Article
	var query string

	// 获取总记录数
	var totalCount int
	if err := ar.db.Get(&totalCount, "SELECT COUNT(*) FROM articles"); err != nil {
		return nil, err
	}

	// 计算偏移量
	offset := pageNumber * pageSize

	// 如果 pageSize 和 pageNumber 都为 0，则返回所有文章
	if pageSize <= 0 || pageNumber < 0 || offset >= totalCount {
		query = "SELECT * FROM articles WHERE title LIKE ?"
	} else {
		query = "SELECT * FROM articles LIMIT ? OFFSET ? WHERE title LIKE ?"
	}

	// 执行查询
	var err error
	if pageSize <= 0 || pageNumber < 0 || offset >= totalCount {
		// 如果需要返回所有文章，执行不带参数的查询
		err = ar.db.Select(&articles, query, "%"+keyword+"%")
	} else {
		// 当限制数量和偏移量有效时，执行带参数的查询
		err = ar.db.Select(&articles, query, pageSize, offset, "%"+keyword+"%")
	}
	if err != nil {
		return nil, err
	}

	// 提取标签
	err = ar.loadArticleTags(&articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// loadArticleTags 用于批量加载文章标签
func (ar *ArticleRepository) loadArticleTags(articles *[]models.Article) error {

	// 创建占位符的字符串
	placeholders := make([]string, len(*articles))
	args := make([]interface{}, len(*articles))

	for i, article := range *articles {
		placeholders[i] = "?" // 或为适用的占位符格式
		args[i] = article.ID  // 将 ID 添加到 args
	}
	// 使用 SQL IN 语法进行一次查询
	query := fmt.Sprintf("SELECT article_id, tag_id FROM article_tags WHERE article_id IN (%s)", strings.Join(placeholders, ", "))

	rows, err := ar.db.Queryx(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	tagsMap := make(map[int][]int)
	for rows.Next() {
		var articleID int
		var tagID int
		if err := rows.Scan(&articleID, &tagID); err != nil {
			fmt.Println(97, err)
			return err
		}
		tagsMap[articleID] = append(tagsMap[articleID], tagID)
	}
	// 将标签添加到对应的文章中
	for i := range *articles {
		if tags, found := tagsMap[(*articles)[i].ID]; found {
			(*articles)[i].Tags = tags
		}
	}
	return nil
}

// CreateArticle 创建文章
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

// GetTags 获取所有标签
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

// CreateLink 创建外部链接
func (ar *ArticleRepository) CreateLink(link models.Link) error {
	_, err := ar.db.NamedExec("INSERT INTO links (title, link,link_type) VALUES (:title, :link, :link_type)", link)
	if err != nil {
		return err
	}
	return nil
}

// GetLinks 获取所有外部链接
func (ar *ArticleRepository) GetLinks() ([]models.Link, error) {
	var links []models.Link
	err := ar.db.Select(&links, "SELECT  * FROM links")
	if err != nil {
		return nil, err
	}
	return links, nil
}
