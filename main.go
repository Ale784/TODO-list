package main

import (
	"fmt"
	"html/template"
	"io"
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

	temp := template.New("Index.html")

	if _, err := temp.ParseFiles("./src/views/Index.html", "./src/views/Detail.html"); err != nil {
		fmt.Println("OPPPS")
	}

	e := echo.New()

	e.Renderer = NewTemplateRenderer(temp)
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())

	ItemsTodo := []Todo{
		{
			Id:        "myfuckingid",
			Title:     "Task 1",
			Detail:    "Buy a new phone",
			Completed: false,
		},
		{
			Id:        "myFUckingTODONUmbert2o",
			Title:     "Task 2",
			Detail:    "Buy a new pc",
			Completed: false,
		},
	}

	e.GET("/", func(c echo.Context) error {
		data := Content{
			PageTitle: "List of todos",
			Items:     ItemsTodo,
		}
		return c.Render(http.StatusOK, "Index.html", data)

	})

	e.GET("/Detail/:Id", func(c echo.Context) error {

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
					Completed: false,
				}}

			}

		}

		SingleData := Content{
			PageTitle: "Details Todo",
			Items:     Item,
		}

		return c.Render(http.StatusOK, "Detail.html", SingleData)

	})

	e.Logger.Fatal(e.Start(":1313"))
}
