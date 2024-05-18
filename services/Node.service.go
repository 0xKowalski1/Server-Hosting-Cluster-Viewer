package services

import (
	Orchestrator "0xKowalski1/container-orchestrator/api-wrapper"
	"0xKowalski1/container-orchestrator/models"
)

type NodeService struct {
	orchestratorWrapper *Orchestrator.WrapperClient
}

func NewNodeService(orchestratorWrapper *Orchestrator.WrapperClient) *NodeService {
	return &NodeService{orchestratorWrapper: orchestratorWrapper}
}

func (service *NodeService) GetNodes() ([]models.Node, error) {
	return service.orchestratorWrapper.ListNodes()
}
