package models

import (
	"fmt"
	"path/filepath"
)

type Article struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	UploadDate  string `json:"uploadDate"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ImageText   string `json:"imageText"`
}
type ArticleInterface interface {
	GetTitle() string
	GetDescription() string
	GetImage() string
	GetImageText() string
}

func Addarticle(art ArticleInterface) error {
	rows, err := DB().Query("SHOW TABLES LIKE 'articles'")
	if err != nil {
		return fmt.Errorf("查询失败:%s", err)
	}
	if !rows.Next() {
		_, err := DB().Exec(`CREATE TABLE IF NOT EXISTS articles(
    		id INT AUTO_INCREMENT PRIMARY KEY,
    		title VARCHAR(255) NOT NULL  ,
      		uploadDate DATETIME DEFAULT CURRENT_TIMESTAMP ,
      		description VARCHAR(255) NOT NULL ,
      		image  BLOB,
    		imageText VARCHAR(255) NOT NULL

		)`)
		if err != nil {
			return fmt.Errorf("createError:%s", err)
		}
	}

	//准备预处理语句
	stmt, err := DB().Prepare("INSERT INTO articles(title, description,image,imageText) VALUES (?,?,?,?)")
	if err != nil {
		fmt.Println(37)
		return fmt.Errorf("stmtError:%s", err)
	}
	defer stmt.Close()
	//执行插入
	_, err = stmt.Exec(art.GetTitle(), art.GetDescription(), art.GetImage(), art.GetImageText())
	if err != nil {
		fmt.Println(44)
		return fmt.Errorf("execError:%s", err)
	}
	return nil
}

func Getarticle() ([]Article, error) {
	rows, err := DB().Query("SELECT * FROM articles")
	if err != nil {
		fmt.Println(26)
		fmt.Printf("错误：%s", err)
		return nil, fmt.Errorf("获取文章失败:%s", err)
	}
	defer rows.Close()
	var articles []Article
	// Iter results
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.Id, &article.Title, &article.UploadDate, &article.Description, &article.Image, &article.ImageText)
		if err != nil {
			return nil, fmt.Errorf("query error:%w", err)
		}
		article.Image = fmt.Sprintf("/images/%s", filepath.Base(article.Image))
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("query error:%w", err)
	}
	return articles, nil
}
