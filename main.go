package main

import (
	"fmt"

	"github.com/labstack/echo/v4"

	echomiddleware "github.com/labstack/echo/v4/middleware"

	"0xKowalski1/cluster-web-viewer/handlers"
	"0xKowalski1/cluster-web-viewer/services"
	Orchestrator "0xKowalski1/container-orchestrator/api"
)

func main() {
	// Create a container orchestrator wrapper
	orchestratorWrapper := Orchestrator.NewApiWrapper("development", "localhost") // Get me from env

	e := echo.New()

	// Static assets
	e.Static("/", "assets")

	// Services
	containerService := services.NewContainerService(orchestratorWrapper)

	// Handlers
	containerHandler := handlers.NewContainerHandler(containerService)

	// Middleware
	e.Use(echomiddleware.Logger())

	// Configure CORS
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Routes
	e.GET("/containers", containerHandler.GetContainers)

	fmt.Printf("Listening on :3001")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":3001")))
}
