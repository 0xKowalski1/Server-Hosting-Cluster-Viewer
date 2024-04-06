package main

import (
	"html/template"
	"io"
	"log"

	"0xKowalski1/container-orchestrator/api"
	"0xKowalski1/container-orchestrator/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var namespace string = "development" // temp

var apiClient *api.WrapperClient // API client instance

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
	apiClient = api.NewApiWrapper(namespace) // Initialize the API client

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
	e.GET("/containers/:id/status", watchContainer)
	e.GET("/containers/:id/logs", streamContainerLogs)

	e.Logger.Fatal(e.Start(":3000"))
}

func listContainers(c echo.Context) error {
	containers, err := apiClient.ListContainers()
	if err != nil {
		return err
	}
	return c.Render(200, "containers", containers)
}

func createContainer(c echo.Context) error {

	req := models.CreateContainerRequest{
		ID:          c.FormValue("id"),     // Use a unique identifier
		Image:       c.FormValue("image"),  // Specify the container image
		Env:         []string{"EULA=TRUE"}, // Any environment variables
		StopTimeout: 5,
	}

	container, err := apiClient.CreateContainer(req)
	if err != nil {
		return err
	}
	return c.Render(200, "container", container)
}

func deleteContainer(c echo.Context) error {

	_, err := apiClient.DeleteContainer(c.Param("id"))

	if err != nil {
		return err
	}
	return nil
}

func startContainer(c echo.Context) error {

	_, err := apiClient.StartContainer(c.Param("id"))

	if err != nil {
		return err
	}
	return nil
}

func stopContainer(c echo.Context) error {

	_, err := apiClient.StopContainer(c.Param("id"))

	if err != nil {
		return err
	}
	return nil
}

func watchContainer(c echo.Context) error {
	containerID := c.Param("id") // Get container ID from the request path

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	c.Response().WriteHeader(200)

	// Stream updates from the orchestrator and send to the client
	err := apiClient.WatchContainer(containerID, func(data string) {
		// Construct an SSE formatted message
		_, writeErr := c.Response().Write([]byte(data + "\n\n"))
		if writeErr != nil {
			log.Printf("Error writing to client: %v", writeErr)

			return
		}
		c.Response().Flush()
	})

	if err != nil {
		// Handle the error
		log.Printf("Error streaming updates: %v", err)
		return err
	}

	select {
	case <-c.Request().Context().Done():
		// The client disconnected
		log.Println("Client disconnected")
	}

	return nil
}

func streamContainerLogs(c echo.Context) error {
	containerID := c.Param("id") // Get container ID from the request path

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	c.Response().WriteHeader(200)

	// Stream updates from the orchestrator and send to the client
	err := apiClient.StreamContainerLogs(containerID, func(data string) {
		// Construct an SSE formatted message
		_, writeErr := c.Response().Write([]byte("data: " + data + "\n\n"))
		if writeErr != nil {
			log.Printf("Error writing to client: %v", writeErr)

			return
		}
		c.Response().Flush()
	})

	if err != nil {
		// Handle the error
		log.Printf("Error streaming updates: %v", err)
		return err
	}

	select {
	case <-c.Request().Context().Done():
		// The client disconnected
		log.Println("Client disconnected")
	}

	return nil
}
