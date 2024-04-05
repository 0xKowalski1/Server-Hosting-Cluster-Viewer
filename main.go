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
	e.DELETE("/containers/:id", deleteContainer)
	e.POST("/containers/:id/start", startContainer)
	e.POST("/containers/:id/stop", stopContainer)

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
	namespace := "example"

	req := apiwrapper.CreateContainerRequest{
		ID:          c.FormValue("id"),     // Use a unique identifier
		Image:       c.FormValue("image"),  // Specify the container image
		Env:         []string{"EULA=TRUE"}, // Any environment variables
		StopTimeout: 5,
	}

	container, err := apiClient.CreateContainer(namespace, req)
	if err != nil {
		return err
	}
	return c.Render(200, "container", container)
}

func deleteContainer(c echo.Context) error {
	namespace := "example"

	_, err := apiClient.DeleteContainer(namespace, c.Param("id"))

	if err != nil {
		return err
	}
	return nil
}

func startContainer(c echo.Context) error {
	namespace := "example"

	_, err := apiClient.StartContainer(namespace, c.Param("id"))

	if err != nil {
		return err
	}
	return nil
}

func stopContainer(c echo.Context) error {
	namespace := "example"

	_, err := apiClient.StopContainer(namespace, c.Param("id"))

	if err != nil {
		return err
	}
	return nil
}
