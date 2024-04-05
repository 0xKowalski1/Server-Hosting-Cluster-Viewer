package apiwrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// API base URL. Adjust as needed.
const BaseURL = "http://localhost:8080"

// Client represents the API client
type Client struct {
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{},
	}
}

type Container struct {
	ID            string   `json:"id"`
	DesiredStatus string   `json:"desiredStatus"` // "running" or "stopped"
	Status        string   `json:"status"`
	NamespaceID   string   `json:"namespaceID"`
	NodeID        string   `json:"nodeID"`
	Image         string   `json:"image"`
	Env           []string `json:"env"`
	StopTimeout   int      `json:"stopTimeout"`
}

type ContainerListResponse struct {
	Containers []Container `json:"containers"`
}

// CreateContainerRequest represents the request to create a container
type CreateContainerRequest struct {
	ID          string   `json:"id"`
	Image       string   `json:"image"`
	Env         []string `json:"env"`
	StopTimeout int      `json:"stopTimeout"`
}

// CreateContainer creates a new container in the specified namespace
func (c *Client) CreateContainer(namespace string, req CreateContainerRequest) (*Container, error) {
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/namespaces/%s/containers", BaseURL, namespace)
	response, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	type ContainerResponse struct {
		Container Container `json:"container"`
	}

	var containerResponse ContainerResponse
	if err := json.NewDecoder(response.Body).Decode(&containerResponse); err != nil {
		return nil, err
	}

	fmt.Printf("Container Response: %+v\n", containerResponse.Container)

	return &containerResponse.Container, nil
}

func (c *Client) ListContainers(namespace string) ([]Container, error) {
	url := fmt.Sprintf("%s/namespaces/%s/containers", BaseURL, namespace)
	response, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	var resp ContainerListResponse // Adjusted to use the new wrapper struct
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return resp.Containers, nil // Return the slice of containers
}

func (c *Client) DeleteContainer(namespace string, containerID string) (string, error) {
	url := fmt.Sprintf("%s/namespaces/%s/containers/%s", BaseURL, namespace, containerID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return containerID, err
	}

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return containerID, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return containerID, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	return containerID, nil
}

func (c *Client) StartContainer(namespace string, containerID string) (string, error) {
	url := fmt.Sprintf("%s/namespaces/%s/containers/%s/start", BaseURL, namespace, containerID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return containerID, err
	}

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return containerID, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return containerID, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	return containerID, nil
}

func (c *Client) StopContainer(namespace string, containerID string) (string, error) {
	url := fmt.Sprintf("%s/namespaces/%s/containers/%s/stop", BaseURL, namespace, containerID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return containerID, err
	}

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return containerID, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return containerID, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	return containerID, nil
}
