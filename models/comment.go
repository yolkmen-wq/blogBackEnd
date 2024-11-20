package models

import (
	"time"
)

type Comment struct {
	ID         int64          `json:"id" db:"id"`
	ArticleID  int64          `json:"articleId" db:"article_id"`
	UserID     int64          `json:"userId" db:"user_id"`
	Content    string         `json:"content" db:"content"`
	LikeCount  int64          `json:"likeCount" db:"like_count"`
	LikeList   []Like         `json:"likeList" `
	IsDelete   int8           `json:"isDelete" db:"is_delete"`
	CreateTime time.Time      `json:"createTime" db:"create_time"`
	Children   []CommentChild `json:"children" db:"children"`
	Nickname   string         `json:"nickname"`
	IsLiked    int8           `json:"isLiked" `
}

type CommentChild struct {
	ParentID *int64 `json:"parentId" db:"parent_id"`
	ReplyID  *int64 `json:"replyId" db:"reply_id"`
	Comment
}

type Like struct {
	ID        int64 `json:"id" db:"id"`
	CommentID int64 `json:"commentId" db:"comment_id"`
	VisitorID int64 `json:"userId" db:"visitor_id"`
}
