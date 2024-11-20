package repositories

import (
	"blog/models"
	"blog/utils"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// getNickname 获取昵称
func getNickname(db *sqlx.DB, userID int64) (string, error) {
	var nickname string
	err := db.QueryRow(`SELECT nickname FROM visitors WHERE id = ?`, userID).Scan(&nickname)
	if err != nil {
		return "", err
	}
	return nickname, nil
}

// getLikeCount 获取点赞数
func getLikeCount(db *sqlx.DB, commentID int64) (int64, error) {
	var count int64
	err := db.Get(&count, `SELECT COUNT(*) FROM visitor_like WHERE comment_id = ?`, commentID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateComment 创建评论
func (r *CommentRepository) CreateComment(comment *models.CommentChild) error {
	err := r.db.QueryRow(`SELECT id FROM post_comment_parent WHERE article_id =? && user_id =?`, comment.ArticleID, comment.UserID).Scan(&comment.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// 没有找到记录，处理此情况
			// 插入父评论
			_, err := r.db.NamedExec(`INSERT INTO post_comment_parent (article_id, user_id, content,create_time ) VALUES (:article_id,:user_id, :content, NOW())`, comment)
			if err != nil {
				fmt.Println(31, err)
				return err
			}
		}
	}

	// 插入子评论
	var query string
	if comment.ReplyID != nil {
		query = `INSERT INTO post_comment_child (article_id,parent_id, user_id,reply_id, content,create_time ) VALUES (:article_id, :parent_id, :user_id,:reply_id,:content, NOW())`

	} else {
		query = `INSERT INTO post_comment_child (article_id,parent_id, user_id, content,create_time ) VALUES (:article_id, :parent_id, :user_id,:content, NOW())`

	}
	_, err = r.db.NamedExec(query, comment)
	if err != nil {
		fmt.Println(45, err)
		return err
	}

	return err
}

// GetCommentByArticleID 根据文章ID获取评论
func (r *CommentRepository) GetCommentByArticleID(articleID int64, visitorID int64) ([]*models.Comment, error) {
	// 获取父评论
	comments := make([]*models.Comment, 0)
	rows, err := r.db.Query(`SELECT id, article_id, user_id, content, create_time FROM post_comment_parent WHERE article_id =? ORDER BY create_time DESC`, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		comment := new(models.Comment)
		var createTimeBytes []byte
		if err := rows.Scan(&comment.ID, &comment.ArticleID, &comment.UserID, &comment.Content, &createTimeBytes); err != nil {
			return nil, err
		}

		// 手动解析时间
		comment.CreateTime, err = utils.ParseTime(createTimeBytes)
		if err != nil {
			return nil, err
		}

		// 获取昵称
		comment.Nickname, err = getNickname(r.db, comment.UserID)
		if err != nil {
			return nil, err
		}

		// 获取点赞数
		comment.LikeCount, err = getLikeCount(r.db, comment.ID)
		if err != nil {
			return nil, err
		}

		// 判断该用户是否已经点赞过
		if visitorID == 0 {
			var isLiked bool
			fmt.Println(comment.ID, visitorID)
			err := r.db.Get(&isLiked, `SELECT COUNT(*) FROM visitor_like WHERE comment_id = ? AND visitor_id = ?`, comment.ID, visitorID)
			if err != nil {
				return nil, err
			}
			if isLiked {
				comment.IsLiked = 1
			} else {
				comment.IsLiked = 0
			}
		}

		comments = append(comments, comment)
	}

	// 获取子评论
	childComments := make([]*models.CommentChild, 0)
	rows, err = r.db.Query(`SELECT id, article_id, parent_id, user_id, reply_id, content, create_time FROM post_comment_child WHERE article_id =? ORDER BY create_time DESC`, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		childComment := new(models.CommentChild)
		var createTimeBytes []byte
		if err := rows.Scan(&childComment.ID, &childComment.ArticleID, &childComment.ParentID, &childComment.UserID, &childComment.ReplyID, &childComment.Content, &createTimeBytes); err != nil {
			return nil, err
		}

		// 手动解析时间
		childComment.CreateTime, err = utils.ParseTime(createTimeBytes)
		if err != nil {
			return nil, err
		}
		// 获取子评论的昵称
		childComment.Nickname, err = getNickname(r.db, childComment.UserID)
		if err != nil {
			return nil, err
		}
		// 获取子评论的点赞数
		childComment.LikeCount, err = getLikeCount(r.db, childComment.ID)
		if err != nil {
			return nil, err
		}

		childComments = append(childComments, childComment)
	}

	// 给父评论添加子评论
	for _, childComment := range childComments {
		// 给每个子评论添加昵称
		query := `SELECT nickname FROM visitors WHERE id = ?`
		row := r.db.QueryRow(query, childComment.UserID)
		err := row.Scan(&childComment.Nickname) // 直接扫描到 childComment.Nickname
		if err != nil {
			return nil, err
		}
		for _, parentComment := range comments {
			if parentComment.ID == *childComment.ParentID {
				parentComment.Children = append(parentComment.Children, *childComment)
				break
			}
		}
	}

	return comments, nil
}

// LikeComment 点赞评论
func (r *CommentRepository) ToggleCommentLike(like *models.Like) error {
	// 先查询是否已经点赞过
	var commentID int64
	query := `SELECT comment_id FROM visitor_like WHERE comment_id = ?`
	err := r.db.QueryRow(query, like.CommentID).Scan(&commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			// 点赞
			_, err := r.db.Exec(`INSERT INTO visitor_like (comment_id, visitor_id) VALUES (?, ?)`, like.CommentID, like.VisitorID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// 取消点赞
		_, err := r.db.Exec(`DELETE FROM visitor_like WHERE comment_id = ? AND visitor_id = ?`, like.CommentID, like.VisitorID)
		if err != nil {
			return err
		}
	}
	return nil
}
