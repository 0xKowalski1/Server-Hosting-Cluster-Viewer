package handlers

import (
	"0xKowalski1/cluster-web-viewer/services"
	"0xKowalski1/cluster-web-viewer/templates"
	"fmt"

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
		fmt.Errorf("Error: %v", err)
	}

	return Render(c, 200, templates.ContainersPage(containers))
}
