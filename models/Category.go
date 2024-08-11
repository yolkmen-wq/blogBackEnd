package models

import (
	"goBlog/utils"
)

func AllCategory() string {
	return utils.RedisDB.Get(utils.CategoriesKey).Val()
}
