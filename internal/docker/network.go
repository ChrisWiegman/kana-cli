package docker

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
)

type ExposedPorts struct {
	Port     string
	Protocol string
}

type portConfig struct {
	PortBindings nat.PortMap
	PortSet      nat.PortSet
}

func (d *Client) EnsureNetwork(name string) (created bool, network types.NetworkResource, err error) {
	hasNetwork, network, err := findNetworkByName(name, d.apiClient)

	if err != nil {
		return false, types.NetworkResource{}, err
	}

	if hasNetwork {
		return false, network, nil
	}

	networkCreateResults, err := d.apiClient.NetworkCreate(context.Background(), name, types.NetworkCreate{
		Driver: "bridge",
	})

	if err != nil {
		return false, types.NetworkResource{}, err
	}

	hasNetwork, network, err = findNetworkByID(networkCreateResults.ID, d.apiClient)

	if err != nil {
		return false, types.NetworkResource{}, err
	}

	if hasNetwork {
		return true, network, nil
	}

	return false, types.NetworkResource{}, fmt.Errorf("could not create network")
}

func (d *Client) RemoveNetwork(name string) (removed bool, err error) {
	hasNetwork, network, err := findNetworkByName(name, d.apiClient)

	if err != nil {
		return false, err
	}

	if !hasNetwork {
		return false, nil
	}

	return true, d.apiClient.NetworkRemove(context.Background(), network.ID)
}

func findNetworkByID(id string, apiClient APIClient) (found bool, network types.NetworkResource, err error) {
	networks, err := apiClient.NetworkList(context.Background(), types.NetworkListOptions{})

	if err != nil {
		return false, types.NetworkResource{}, err
	}

	for i := range networks {
		if networks[i].ID == id {
			return true, networks[i], nil
		}
	}

	return false, types.NetworkResource{}, nil
}

func findNetworkByName(name string, apiClient APIClient) (found bool, network types.NetworkResource, err error) {
	networks, err := apiClient.NetworkList(context.Background(), types.NetworkListOptions{})

	if err != nil {
		return false, types.NetworkResource{}, err
	}

	for i := range networks {
		if networks[i].Name == name {
			return true, networks[i], nil
		}
	}

	return false, types.NetworkResource{}, nil
}

func getNetworkConfig(ports []ExposedPorts, randomPorts bool) (portConfig, error) {
	portBindings := make(nat.PortMap)
	portSet := make(nat.PortSet)

	for _, port := range ports {
		portName, err := nat.NewPort(port.Protocol, port.Port)
		if err != nil {
			return portConfig{}, err
		}

		hostPort := port.Port

		if randomPorts {
			port, err := getRandomPort()
			if err != nil {
				return portConfig{}, err
			}

			hostPort = port
		}

		portBindings[portName] = []nat.PortBinding{
			{
				HostPort: hostPort,
			},
		}

		portSet[portName] = struct{}{}
	}

	return portConfig{
		PortBindings: portBindings,
		PortSet:      portSet,
	}, nil
}

// getRandomPort Returns an open, ephemeral port for mapping a container.
func getRandomPort() (string, error) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	server := httptest.NewServer(handler)
	urlParts, err := url.ParseRequestURI(server.URL)

	server.Close()

	if err != nil {
		return "", err
	}

	return urlParts.Port(), nil
}
