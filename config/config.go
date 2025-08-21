package config

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

// 定义一个结构体用于表示要返回的数据
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type AppConfig struct {
	DatabaseUrl string
	Port        string
	Username    string
	Password    string
}

var Config AppConfig
var db *sqlx.DB
var client *redis.Client
var err error

// 定义一个密钥
var SecretKey = []byte("lyf123456")
var CaptchaSource = "0123456789abcdefghijklmnopqrstuvwxyz"

func InitConfig() {
	Config = AppConfig{
		DatabaseUrl: "47.121.201.137",
		Port:        "3306",
		Username:    "root",
		Password:    "root",
	}
	fmt.Println("Config initialized")
}

func DB() *sqlx.DB {
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
func newDB() (*sqlx.DB, error) {
	//配置数据库连接信息
	dsn := "root:wuqi9457@tcp(47.121.201.137:3306)/my_blog"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(100)                // 最大打开连接数
	db.SetMaxIdleConns(50)                 // 最大空闲连接数
	db.SetConnMaxLifetime(time.Minute * 5) // 连接的最大生命周期

	//检查连接是否成功

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connect to Mysql database")
	return db, nil
}

func RedisClient() *redis.Client {
	if client == nil {
		newClient := newRedis()
		client = newClient
	}
	return client
}

func newRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "47.121.201.137:6011",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// 测试连接
	_, err := client.Ping().Result()
	if err != nil {
		log.Println("Could not connect to Redis: %v", err)
	} else {
		log.Println("Successfully connected to Redis")
	}
	return client
}

// 过滤不需要验证的路由
var ExcludeRoutes = []string{"/login", "/register", "/getCaptcha", "/articles", "/getArticleById", "/createArticle", "/getTags", "/createLink", "/getLinks", "/countVisits"}
