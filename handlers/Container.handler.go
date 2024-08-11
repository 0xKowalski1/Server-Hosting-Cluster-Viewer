package handlers

import (
	"0xKowalski1/cluster-web-viewer/services"
	"0xKowalski1/cluster-web-viewer/templates"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"0xKowalski1/container-orchestrator/models"

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

func (handler *ContainerHandler) NewContainer(c echo.Context) error {
	formData := models.CreateContainerRequest{
		ID:    "minecraft-server-1",
		Image: "ghcr.io/0xkowalski1/minecraft-server:latest",
		Env: []string{
			"EULA=TRUE",
			"MEMORY=4",
		},
		StopTimeout:  5,
		MemoryLimit:  2,
		CpuLimit:     1,
		StorageLimit: 5,
		Ports: []models.Port{
			{
				HostPort:      30001,
				ContainerPort: 25565,
				Protocol:      "TCP",
			},
		},
	}

	return Render(c, 200, templates.NewContainersPage(formData))
}

func ConvertPortsToString(ports []models.Port) string {
	var portsStr []string
	for _, port := range ports {
		portsStr = append(portsStr, strconv.Itoa(port.HostPort)+":"+strconv.Itoa(port.ContainerPort)+"/"+port.Protocol)
	}
	return strings.Join(portsStr, ",")
}

func ParsePorts(portsStr string) ([]models.Port, error) {
	var ports []models.Port
	portsList := strings.Split(portsStr, ",")
	for _, portStr := range portsList {
		parts := strings.Split(portStr, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid port format")
		}
		hostPort, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		containerParts := strings.Split(parts[1], "/")
		if len(containerParts) != 2 {
			return nil, fmt.Errorf("invalid port format")
		}
		containerPort, err := strconv.Atoi(containerParts[0])
		if err != nil {
			return nil, err
		}
		protocol := containerParts[1]
		ports = append(ports, models.Port{
			HostPort:      hostPort,
			ContainerPort: containerPort,
			Protocol:      protocol,
		})
	}
	return ports, nil
}

func (handler *ContainerHandler) CreateContainer(c echo.Context) error {
	id := c.FormValue("id")
	image := c.FormValue("image")
	env := strings.Split(c.FormValue("env"), ",")
	stopTimeout, err := strconv.Atoi(c.FormValue("stopTimeout"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid stop timeout")
	}
	memoryLimit, err := strconv.Atoi(c.FormValue("memoryLimit"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid memory limit")
	}
	cpuLimit, err := strconv.Atoi(c.FormValue("cpuLimit"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid CPU limit")
	}
	storageLimit, err := strconv.Atoi(c.FormValue("storageLimit"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid storage limit")
	}
	ports, err := ParsePorts(c.FormValue("ports"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ports format")
	}

	newContainerRequest := models.CreateContainerRequest{
		ID:           id,
		Image:        image,
		Env:          env,
		StopTimeout:  stopTimeout,
		MemoryLimit:  memoryLimit,
		CpuLimit:     cpuLimit,
		StorageLimit: storageLimit,
		Ports:        ports,
	}

	_, err = handler.containerService.CreateContainer(newContainerRequest)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	c.Response().Header().Set("HX-Replace-Url", "/containers")
	return handler.GetContainers(c)
}

func (handler *ContainerHandler) DeleteContainer(c echo.Context) error {
	err := handler.containerService.DeleteContainer(c.Param("containerID"))

	if err != nil {
		log.Printf("Error: %v", err)
	}

	c.Response().Header().Set("HX-Replace-Url", "/containers")
	return handler.GetContainers(c)
}

func (handler *ContainerHandler) GetContainerLogs(c echo.Context) error {
	return Render(c, 200, templates.ContainerLogsPage(c.Param("containerID")))
}

func (handler *ContainerHandler) StreamContainerLogs(c echo.Context) error {
	containerID := c.Param("containerID")
	log.Println("Streaming logs for container:", containerID)

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	clientGone := c.Request().Context().Done()

	handleData := func(logLine string) {
		select {
		case <-clientGone:
			log.Println("Client disconnected, stopping log stream")
			return
		default:
			message := "<div>" + logLine + "</div>\n"
			_, err := c.Response().Write([]byte("data: " + message + "\n\n"))
			if err != nil {
				log.Println("Error writing to response:", err)
				return
			}
			flusher.Flush()
		}
	}

	err := handler.containerService.StreamContainerLogs(containerID, handleData)
	if err != nil {
		log.Println("Error streaming container logs:", err)
		return err
	}

	return nil
}
