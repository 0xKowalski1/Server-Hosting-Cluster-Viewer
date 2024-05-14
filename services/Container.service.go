package services

import (
	Orchestrator "0xKowalski1/container-orchestrator/api"
	"0xKowalski1/container-orchestrator/models"
)

type ContainerService struct {
	orchestratorWrapper *Orchestrator.WrapperClient
}

func NewContainerService(orchestratorWrapper *Orchestrator.WrapperClient) *ContainerService {
	return &ContainerService{orchestratorWrapper: orchestratorWrapper}
}

func (service *ContainerService) GetContainers() ([]models.Container, error) {
	return service.orchestratorWrapper.ListContainers()
}

func (service *ContainerService) GetContainer(containerID string) (*models.Container, error) {
	return service.orchestratorWrapper.GetContainer(containerID)
}

func (service *ContainerService) DeleteContainer(containerID string) error {
	return service.orchestratorWrapper.DeleteContainer(containerID)
}
