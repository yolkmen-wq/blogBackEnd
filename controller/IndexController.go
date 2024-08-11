package controller

import (
	"github.com/labstack/echo/v4"
	"goBlog/services"
	"goBlog/utils"
	"goBlog/utils/spider/tian_kong"
	"net/http"
	"strconv"
	"strings"
)

func Index(c echo.Context) error {
	show := make(map[string]interface{})

	NewFilmKey := "detail_links:id:1"
	NewTVKey := "detail_links:id:3"
	NewCartoonKey := "detail_links:id:24"
	NewFilm := services.MovieListsRange(NewFilmKey, 0, 14)
	NewTV := services.MovieListsRange(NewTVKey, 0, 14)
	NewCartoon := services.MovieListsRange(NewCartoonKey, 0, 14)

	show["newFilm"] = NewFilm
	show["newTv"] = NewTV
	show["newCartoon"] = NewCartoon

	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)

	show["navFilm"] = getAssignTypeSubCategories(Categories, "film")
	show["navTv"] = getAssignTypeSubCategories(Categories, "tv")
	return c.Render(http.StatusOK, "index", show)
}

func Display(c echo.Context) error {
	path := c.Path()
	cate := c.QueryParam("cate")
	_start := c.QueryParam("start")
	_stop := c.QueryParam("stop")

	show := make(map[string]interface{})

	key := "detail_links:id:5" // 默认首页

	start := int64(0)
	stop := int64(41)

	if len(_start) > 0 {
		StartInt64, _ := strconv.ParseInt(_start, 10, 64)
		start = StartInt64
	}

	if len(_stop) > 0 {
		StopInt64, _ := strconv.ParseInt(_stop, 10, 64)
		stop = StopInt64
	}

	prev := path + "?start=" + strconv.FormatInt(start-42, 10) + "&stop=" + strconv.FormatInt(stop-42, 10)
	next := path + "?start=" + strconv.FormatInt(start+42, 10) + "&stop=" + strconv.FormatInt(stop+42, 10)

	prevStatus := "1"
	nextStatus := "1"

	cateStrId := services.TransformCategoryId(cate)
	cateIntId, _ := strconv.Atoi(cateStrId)

	if len(cate) > 0 {
		key = "detail_links:id:" + cateStrId
		prev = path + "?cate=" + cate + "&start=" + strconv.FormatInt(start-42, 10) + "&stop=" + strconv.FormatInt(stop-42, 10)
		next = path + "?cate=" + cate + "&start=" + strconv.FormatInt(start+42, 10) + "&stop=" + strconv.FormatInt(stop+42, 10)
	}

	if start > stop || stop-start > 42 || start < 0 {
		start = 0
		stop = 41
	}

	MovieLists := services.MovieListsRange(key, start, stop)

	LenMovieLists := len(MovieLists)

	if start-42 < 0 {
		prevStatus = "0"
	}

	if LenMovieLists < 42 || LenMovieLists == 0 {
		nextStatus = "0"
	}

	show["movieLists"] = MovieLists
	show["prev"] = prev
	show["next"] = next
	show["prev_status"] = prevStatus
	show["next_status"] = nextStatus

	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)

	// 根据不同类别显示不同 筛选类别项
	if utils.InArray(cateIntId, tian_kong.GetAssignCategoryIds("film")) || cateIntId == 1 {
		show["currentSubCate"] = getAssignTypeSubCategories(Categories, "film")
	}

	if utils.InArray(cateIntId, tian_kong.GetAssignCategoryIds("tv")) || cateIntId == 3 {
		show["currentSubCate"] = getAssignTypeSubCategories(Categories, "tv")
	}

	if utils.InArray(cateIntId, tian_kong.GetAssignCategoryIds("cartoon")) || cateIntId == 24 {
		show["currentSubCate"] = getAssignTypeSubCategories(Categories, "cartoon")
	}

	//tmpl.GoTpl.ExecuteTemplate(w, "display", show)
	return c.JSON(http.StatusOK, show)
}

func Movie(c echo.Context) error {
	link := c.QueryParam("link")
	if link == "" {
		return c.String(http.StatusNotFound, "资源未找到")
	}

	show := make(map[string]interface{})
	MovieDetail := services.MovieDetail(link)

	if len(MovieDetail["info"].(map[string]string)) == 0 {
		return c.Render(http.StatusNotFound, "404", show)
	}

	show["MovieDetail"] = MovieDetail
	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)
	return c.Render(http.StatusOK, "detail", show)

	//tmpl.GoTpl.ExecuteTemplate(w, "detail", show)
}

func Play(c echo.Context) error {
	show := make(map[string]interface{})
	PlayUrl := c.QueryParam("play_url")
	PlayType := "kuyun"
	if strings.Contains(PlayUrl, ".mp4") {
		PlayType = "mp4"
	} else if strings.Contains(PlayUrl, ".m3u8") {
		PlayType = "m3u8"
	}

	show["playUrl"] = PlayUrl
	show["playType"] = PlayType

	link := c.QueryParam("link")
	episode := c.QueryParam("episode")
	MovieDetail := services.MovieDetail(link)

	if len(MovieDetail["info"].(map[string]string)) == 0 {
		return c.Render(http.StatusNotFound, "404", show)
		return nil
	}

	show["MovieDetail"] = MovieDetail
	show["episode"] = episode
	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)
	return c.Render(http.StatusOK, "play", show)
}

func Search(c echo.Context) error {
	show := make(map[string]interface{})
	q := c.QueryParam("q")

	var MovieLists []services.MovieListStruct
	if strings.TrimSpace(q) != "" {
		MovieLists = services.SearchMovies(q)
	}

	show["movieLists"] = MovieLists
	show["q"] = q
	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)
	//tmpl.GoTpl.ExecuteTemplate(w, "search", show)
	return c.Render(http.StatusOK, "search", show)

}

func About(c echo.Context) error {
	show := make(map[string]interface{})
	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)
	//tmpl.GoTpl.ExecuteTemplate(w, "about", show)
	return c.Render(http.StatusOK, "about", show)

}
func getAllCategory(Categories []utils.Categories) []utils.Categories {
	var allC []utils.Categories
	for _, c := range Categories {
		allC = append(allC, c)
		for _, subC := range c.Sub {
			allC = append(allC, subC)
		}
	}
	return allC
}

func getAssignTypeSubCategories(Categories []utils.Categories, _type string) []utils.Categories {
	var cate []utils.Categories
	switch _type {
	case "film":
		cate = Categories[0].Sub
	case "tv":
		cate = Categories[1].Sub
	case "cartoon":
		cate = Categories[2].Sub
	}
	return cate
}
