package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
)

// Struct TODO
type Todo struct {
	Title        string
	Completed    bool
	OnHold       bool
	PlanToFinish bool
	Doing        bool
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

func main() {

	temp := template.Must(template.ParseFiles("index.html"))

	e := echo.New()

	e.Renderer = NewTemplateRenderer(temp)
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		data := Content{
			PageTitle: "My first go template",
			Items: []Todo{
				{Title: "Task 1", Completed: false, OnHold: false, PlanToFinish: false, Doing: false},
			},
		}
		return c.Render(http.StatusOK, "index.html", data)

	})

	e.Logger.Fatal(e.Start(":1313"))
}
