package main

import (
	"context"
	"html/template"
	"io"
	"log"
	"strconv"

	"0xKowalski1/container-orchestrator/api"
	"0xKowalski1/container-orchestrator/models"

	"strings"

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
	stopTimeout, err := strconv.Atoi(c.FormValue("stopTimeout"))
	if err != nil {
		stopTimeout = 5
	}

	memoryLimit, _ := strconv.Atoi(c.FormValue("memoryLimit"))

	cpuLimit, _ := strconv.Atoi(c.FormValue("cpuLimit"))

	// Splitting the comma-separated environment variables string into a slice of strings.
	env := strings.Split(c.FormValue("env"), ",")

	req := models.CreateContainerRequest{
		ID:          c.FormValue("id"),
		Image:       c.FormValue("image"),
		Env:         env,
		StopTimeout: stopTimeout,
		MemoryLimit: memoryLimit,
		CpuLimit:    cpuLimit}

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

	ctx, cancel := context.WithCancel(c.Request().Context()) // Create a cancellable context
	defer cancel()                                           // Ensure cancel is called to clean up the context

	// Stream updates from the orchestrator and send to the client
	err := apiClient.WatchContainer(containerID, func(data string) {
		select {
		case <-ctx.Done(): // Check if client has disconnected
			log.Println("Status client disconnected, stopping data stream.")
			return
		default:
			_, writeErr := c.Response().Write([]byte(data + "\n\n"))
			if writeErr != nil {
				log.Printf("Error writing to client: %v", writeErr)
				cancel() // Stop the data stream on error
				return
			}
			c.Response().Flush()
		}
	})

	if err != nil {
		log.Printf("Error streaming updates: %v", err)
		return err
	}

	<-ctx.Done() // Wait for the context to be cancelled
	log.Println("Finished streaming container updates.")
	return nil
}

func streamContainerLogs(c echo.Context) error {
	containerID := c.Param("id") // Get container ID from the request path

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(200)

	ctx, cancel := context.WithCancel(c.Request().Context()) // Create a cancellable context
	defer cancel()                                           // Ensure cancel is called to clean up the context

	// Stream updates from the orchestrator and send to the client
	err := apiClient.StreamContainerLogs(containerID, func(data string) {
		select {
		case <-ctx.Done(): // Check if client has disconnected
			log.Println("Stream client disconnected, stopping log stream.")
			return
		default:
			_, writeErr := c.Response().Write([]byte("data: " + data + "\n\n"))
			if writeErr != nil {
				log.Printf("Error writing to client: %v", writeErr)
				cancel() // Stop the log stream on error
				return
			}
			c.Response().Flush()
		}
	})

	if err != nil {
		log.Printf("Error streaming container logs: %v", err)
		return err
	}

	<-ctx.Done() // Wait for the context to be cancelled
	log.Println("Finished streaming container logs.")
	return nil
}
