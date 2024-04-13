package main

import (
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"0xKowalski1/container-orchestrator/api"
	// "0xKowalski1/container-orchestrator/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	tmpl := template.New("")
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			_, err = tmpl.ParseFiles(path)
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}
	return &Template{tmpl: tmpl}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

func nodesPageHandler(c echo.Context) error {
	// Handle page refresh
	if c.Request().Header.Get("HX-Request") == "" {
		return c.Render(200, "index", nil)
	}

	apiClient := api.NewApiWrapper("development")
	nodes, err := apiClient.ListNodes()
	if err != nil {
		// Handle Errors
		log.Printf("Error listing nodes: %v", err)
	}

	return c.Render(200, "nodes-page", nodes)
}

func containersPageHandler(c echo.Context) error {
	// Handle page refresh
	if c.Request().Header.Get("HX-Request") == "" {
		return c.Render(200, "index", nil)
	}

	apiClient := api.NewApiWrapper("development")
	containers, err := apiClient.ListContainers()
	if err != nil {
		// Handle Errors
	}

	return c.Render(200, "containers-page", containers)
}

func homePageHandler(c echo.Context) error {
	// Handle page refresh
	if c.Request().Header.Get("HX-Request") == "" {
		return c.Render(200, "index", nil)
	}

	return c.Render(200, "home-page", nil)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	e.GET("/", homePageHandler)

	e.GET("/containers", containersPageHandler)

	e.GET("/nodes", nodesPageHandler)

	e.Logger.Fatal(e.Start(":3000"))
}
