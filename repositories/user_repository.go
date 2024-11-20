package repositories

import (
	"blog/models"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// 创建用户
func (ur *UserRepository) CreateUser(user *models.User) error {
	stmt, err := ur.db.Prepare("INSERT INTO users (username, password) VALUES (?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Password)
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
		fmt.Println(54, err)
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
