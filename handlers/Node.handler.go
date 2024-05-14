package handlers

import (
	"0xKowalski1/cluster-web-viewer/services"
	"0xKowalski1/cluster-web-viewer/templates"
	"log"

	"github.com/labstack/echo/v4"
)

type NodeHandler struct {
	nodeService *services.NodeService
}

func NewNodeHandler(nodeService *services.NodeService) *NodeHandler {
	return &NodeHandler{
		nodeService: nodeService,
	}
}

func (handler *NodeHandler) GetNodes(c echo.Context) error {
	nodes, err := handler.nodeService.GetNodes()

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return Render(c, 200, templates.NodesPage(nodes))
}
