package handlers

import (
	"0xKowalski1/cluster-web-viewer/services"
	"0xKowalski1/cluster-web-viewer/templates"
	"log"

	"github.com/labstack/echo/v4"
)

type ContainerHandler struct {
	containerService *services.ContainerService
}

func NewContainerHandler(containerService *services.ContainerService) *ContainerHandler {
	return &ContainerHandler{
		containerService: containerService,
	}
}

func (handler *ContainerHandler) GetContainers(c echo.Context) error {
	containers, err := handler.containerService.GetContainers()

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return Render(c, 200, templates.ContainersPage(containers))
}

func (handler *ContainerHandler) GetContainer(c echo.Context) error {
	container, err := handler.containerService.GetContainer(c.Param("containerID"))

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return Render(c, 200, templates.ContainerPage(*container))
}

func (handler *ContainerHandler) DeleteContainer(c echo.Context) error {
	err := handler.containerService.DeleteContainer(c.Param("containerID"))

	if err != nil {
		log.Printf("Error: %v", err)
	}

	c.Response().Header().Set("HX-Replace-Url", "/containers")
	return handler.GetContainers(c)
}
