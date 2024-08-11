package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var db *sql.DB
var err error

func DB() *sql.DB {
	if db == nil {
		newDb, err := newDB()
		if err != nil {
			log.Fatal(err)
		}
		//设置数据库连接池的最大打开连接数和最大空闲连接数
		newDb.SetMaxOpenConns(10)
		newDb.SetMaxIdleConns(5)
		newDb.SetConnMaxLifetime(time.Hour)
		db = newDb
	}
	return db
}
func newDB() (*sql.DB, error) {
	//配置数据库连接信息
	dsn := "root:wuqi9457@tcp(47.121.201.137:3306)/my_blog"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	//检查连接是否成功

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connect to Mysql database")
	return db, nil
}
