package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Todo struct {
	Id        string
	Title     string
	Detail    string
	Completed bool
}

type Content struct {
	PageTitle string
	Items     []Todo
}

type TemplateRenderer struct {
	tmpl *template.Template
}

func NewTemplateRenderer(tmpls *template.Template) TemplateRenderer {
	return TemplateRenderer{tmpls}
}

func (t TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

type ProjectPath struct {
	ID string `param:"ID"`
}

func main() {
	const (
		indexPath  string = "./public/views/*.html"
		staticPath string = "./public/static"
	)

	temp := template.New("index.html")

	ItemsTodo := []Todo{
		{
			Id:        "id1",
			Title:     "Task 1",
			Detail:    "Buy a new phone",
			Completed: false,
		},
		{
			Id:        "id2",
			Title:     "Task 2",
			Detail:    "Buy a new pc",
			Completed: false,
		},
	}

	if _, err := temp.ParseGlob(indexPath); err != nil {
		log.Fatalf("unable to parse glob %e", err)
	}

	e := echo.New()

	e.Static("public/static/", staticPath)

	e.Renderer = NewTemplateRenderer(temp)
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())

	e.GET("/", func(c echo.Context) error {
		data := Content{
			PageTitle: "List of todos",
			Items:     ItemsTodo,
		}

		return c.Render(http.StatusOK, "index.html", data)

	})

	e.GET("/detail/:Id", func(c echo.Context) error {

		path := new(ProjectPath)

		if err := (&echo.DefaultBinder{}).BindPathParams(c, path); err != nil {
			return err
		}

		projectID := path.ID
		Item := []Todo{}

		for Key, _ := range ItemsTodo {

			i := ItemsTodo[Key]

			if i.Id == projectID {

				Item = []Todo{{
					Id:        i.Id,
					Title:     i.Title,
					Detail:    i.Detail,
					Completed: i.Completed,
				}}

			}

		}

		SingleData := Content{
			PageTitle: "details Todo",
			Items:     Item,
		}

		return c.Render(http.StatusOK, "detail.html", SingleData)

	})

	e.Logger.Fatal(e.Start(":1313"))

}
