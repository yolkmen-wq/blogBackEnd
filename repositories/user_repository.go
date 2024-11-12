package repositories

import (
	"blog/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

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
