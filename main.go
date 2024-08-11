package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"goBlog/config"
	"goBlog/router"
	"goBlog/utils"
	"goBlog/utils/spider"
	"os"
)

func init() {
	viper.SetConfigType("json")
	readConfig := errors.New("未定义配置文件")

	if _, err := os.Stat("./app.json"); os.IsNotExist(err) {
		readConfig = viper.ReadConfig(bytes.NewBuffer(config.AppJsonConfig))
	} else {
		viper.SetConfigName("app")
		viper.AddConfigPath(".")
		readConfig = viper.ReadInConfig()
	}

	if err := readConfig; err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println(errors.New("config file not found"))
		} else {
			fmt.Println(errors.New("config file was found but another error was produced"))
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}
}

// 首次启动自动开启爬虫
func firstSpider() {

	hasHome := utils.RedisDB.Exists("detail_links:id:1").Val()
	fmt.Println("hasHome", hasHome)
	// 不存在首页的key 则认为是第一次启动
	if hasHome == 0 {
		spider.Create().Start()
	}
}

func main() {

	// 初始化 redis 连接
	utils.InitRedisDB()
	defer utils.CloseRedisDB()

	port := viper.GetString(`app.port`)
	mod := viper.GetString(`app.spider_mod`)
	fmt.Println("监听端口", "http://127.0.0.1"+port)
	fmt.Println("spider_mod：" + mod)

	firstSpider()

	// 启动定时爬虫任务 全量
	utils.TimingSpider(func() {
		spider.Create().Start()
		return
	})

	// 爬虫 只爬取最近有更新的资源
	utils.RecentUpdate(func() {
		spider.Create().DoRecentUpdate()
		return
	})
	// 注册所有路由
	router.NewServer()
}
