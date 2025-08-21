package repositories

import (
	"blog/models"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db     *sqlx.DB
	client *redis.Client
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// 登录
func (ur *UserRepository) FindUser(username string, password string) (*models.User, error) {
	var user models.User
	err := ur.db.Get(&user, "SELECT * FROM users WHERE username=? AND password=?", username, password)
	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

// 创建用户
func (ur *UserRepository) CreateUser(user *models.User) error {
	fmt.Println(32, user)
	// 验证用户名是否存在
	var count int
	err := ur.db.Get(&count, "SELECT COUNT(*) FROM users WHERE username=?", user.Username)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Username already exists")
	}
	_, err = ur.db.NamedExec("INSERT INTO users (username,password,email) VALUES (:username,:password,:email)", user)
	if err != nil {
		return err
	}
	return nil
}

// 根据ID获取用户
func (ur *UserRepository) GetUserById(id int) (*models.User, error) {
	stmt, err := ur.db.Prepare("SELECT * FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(id).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) CreateVisitor(visitor *models.Visitor) (*models.Visitor, error) {
	var exists bool
	query := `SELECT COUNT(*) > 0 FROM visitors WHERE ip = ? OR nickname = ?`
	err := ur.db.QueryRow(query, visitor.IP, visitor.Nickname).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("Visitor already exists")
	}
	result, err := ur.db.NamedExec("INSERT INTO visitors (nickname,email,ip ) VALUES (:nickname,:email,:ip )", visitor)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	visitor.ID = int(id)
	return visitor, nil
}
