package spider

import (
	"github.com/spf13/viper"
	"goBlog/utils"
	"goBlog/utils/spider/tian_kong"
	"goBlog/utils/spider/tian_kong/tian_kong_sync"
)

// 定义 mod 的映射关系
var spiderModMap = map[string]utils.SpiderTask{
	"async": &tian_kong.SpiderApi{},      // 异步 goroutine 并行
	"sync":  &tian_kong_sync.SpiderApi{}, // 同步 按顺序执行 串行
}

func Create() utils.SpiderTask {
	mod := viper.GetString(`app.spider_mod`)
	return spiderModMap[mod]
}
