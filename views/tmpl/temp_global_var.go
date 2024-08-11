package tmpl

import (
	"embed"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"time"
)

// main 函数中定义的全局变量 在其他包中无法调用 ，重新定义新包实现全局变量
var (
	//go:embed *.html
	embedTmpl embed.FS
	//
	//// 自定义的函数必须在调用ParseFiles() ParseFS()之前创建。
	funcMap = template.FuncMap{
		"add": func(k1, k2 int) int {
			return k1 + k2
		},
		"year": func() int {
			return time.Now().Year()
		},
	}
	//GoTpl = template.Must(
	//	template.New("").
	//		Funcs(funcMap).
	//		ParseFS(embedTmpl, "*.html")) // 利用 air 监听文件变动 实时重新加载。修改html无须手动重启服务
	//GoTpl = template.Must(template.ParseGlob("./views/tmpl/*.html"))

)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func Temp() *Template {
	tpl := &Template{templates: template.Must(template.New("").Funcs(funcMap).ParseFS(embedTmpl, "*.html"))}
	return tpl
}
