package main

import (
	"fmt"

	"github.com/labstack/echo/v4"

	echomiddleware "github.com/labstack/echo/v4/middleware"

	"0xKowalski1/cluster-web-viewer/handlers"
	"0xKowalski1/cluster-web-viewer/services"
	Orchestrator "0xKowalski1/container-orchestrator/api-wrapper"
)

func main() {
	// Create a container orchestrator wrapper
	orchestratorWrapper := Orchestrator.NewApiWrapper("localhost") // Get me from env

	e := echo.New()

	// Static assets
	e.Static("/", "assets")

	// Services
	containerService := services.NewContainerService(orchestratorWrapper)
	nodeService := services.NewNodeService(orchestratorWrapper)

	// Handlers
	homeHandler := handlers.NewHomeHandler()
	containerHandler := handlers.NewContainerHandler(containerService)
	nodeHandler := handlers.NewNodeHandler(nodeService)

	// Middleware
	//e.Use(echomiddleware.Logger())

	// Configure CORS
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Routes
	e.GET("/", homeHandler.GetHome)

	// /containers
	e.GET("/containers", containerHandler.GetContainers)
	e.GET("/containers/:containerID", containerHandler.GetContainer)
	e.GET("/containers/:containerID/logs", containerHandler.GetContainerLogs)
	e.GET("/containers/:containerID/logs/stream", containerHandler.StreamContainerLogs)
	e.GET("/containers/new", containerHandler.NewContainer)
	e.POST("/containers", containerHandler.CreateContainer)
	e.DELETE("/containers/:containerID", containerHandler.DeleteContainer)

	// /nodes
	e.GET("/nodes", nodeHandler.GetNodes)

	fmt.Printf("Listening on :3001")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":3001")))
}
