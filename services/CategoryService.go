package services

import (
	"goBlog/models"
	"goBlog/utils"
)

func AllCategoryData() []utils.Categories {
	categories := models.AllCategory()

	var nav []utils.Categories
	utils.Json.Unmarshal([]byte(categories), &nav)

	return nav
}
