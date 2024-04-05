package main

import (
	"html/template"
	"io"

	"github.com/0xKowalski1/server-hosting/apiwrapper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var apiClient *apiwrapper.Client // API client instance

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

func main() {
	apiClient = apiwrapper.NewClient() // Initialize the API client

	e := echo.New()

	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})
	e.GET("/containers", listContainers)
	e.POST("/create", createContainer)

	e.Logger.Fatal(e.Start(":3000"))
}

func listContainers(c echo.Context) error {
	containers, err := apiClient.ListContainers("example")
	if err != nil {
		return err
	}
	return c.Render(200, "containers", containers)
}

func createContainer(c echo.Context) error {
	namespace := "example" // Your namespace here

	req := apiwrapper.CreateContainerRequest{
		ID:    "minecraft-server",                       // Use a unique identifier
		Image: "docker.io/itzg/minecraft-server:latest", // Specify the container image
		Env:   []string{"EULA=TRUE"},                    // Any environment variables
	}

	_, err := apiClient.CreateContainer(namespace, req)
	if err != nil {
		return err
	}
	return nil
}
