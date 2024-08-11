package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	Account  string
	Password string
}

func QueryUser(acc string, pwd string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//执行查询
	fmt.Println(2222, acc, pwd)

	row := DB().QueryRowContext(ctx, "SELECT acc,pwd FROM users WHERE acc=? AND pwd=?", acc, pwd)
	if err != nil {
		return nil, fmt.Errorf("query error:%w", err)
	}
	var user User
	err = row.Scan(&user.Account, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		} else {
			return &user, nil
		}
	}
	fmt.Println(34, user)
	////处理结果集
	//var users []User
	//for rows.Next() {
	//	var user User
	//	err := rows.Scan(&user.Account, &user.Password)
	//	if err != nil {
	//		return nil, fmt.Errorf("query error:%w", err)
	//	}
	//	users = append(users, user)
	//}
	//if err = rows.Err(); err != nil {
	//	return nil, fmt.Errorf("query error:%w", err)
	//}
	return &user, nil
}

func InsertUser(acc string, pwd string) error {
	smt, err := DB().Prepare("INSERT INTO users(acc,pwd) VALUES(?,?)")
	if err != nil {
		return fmt.Errorf("db error:%w", err)
	}
	res, err := smt.Exec(acc, pwd)
	if err != nil {
		return fmt.Errorf("插入失败:%w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("插入失败:%w", err)
	}
	_, err = DB().Exec("INSERT INTO user_roles (user_id,role_id) VALUES (?,?)", id, 2)
	if err != nil {
		return fmt.Errorf("注册失败:%w", err)
	}
	return nil
}
