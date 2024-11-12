package models

import "database/sql"

type Article struct {
	ID      int            `json:"id"`
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Brief   string         `json:"brief"`
	Tags    []int          `json:"tags"`
	Date    string         `json:"date"`
	Url     sql.NullString `json:"url"`
}

type Tag struct {
	TagID   int    `json:"tagId" db:"id"`
	TagName string `json:"tagName" db:"name"`
	Count   int    `json:"count"`
}
