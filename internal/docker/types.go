package docker

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	networkTypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// APIClient is an interface that clients that talk with a docker server must implement.
type APIClient interface {
	ContainerAPIClient
	ImageAPIClient
	NetworkAPIClient
}

// Ensure that Client always implements APIClient.
var _ APIClient = &client.Client{}

// ContainerAPIClient defines API client methods for the containers.
type ContainerAPIClient interface {
	ContainerCreate(
		ctx context.Context,
		config *container.Config,
		hostConfig *container.HostConfig,
		networkingConfig *networkTypes.NetworkingConfig,
		platform *specs.Platform,
		containerName string) (container.CreateResponse, error)
	ContainerAttach(ctx context.Context, container string, options container.AttachOptions) (types.HijackedResponse, error)
	ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error)
	ContainerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error)
	ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error)
	ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error)
	ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	ContainerLogs(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error)
	ContainerRemove(ctx context.Context, container string, options container.RemoveOptions) error
	ContainerStart(ctx context.Context, container string, options container.StartOptions) error
	ContainerStop(ctx context.Context, name string, options container.StopOptions) error
	ContainerWait(
		ctx context.Context,
		container string,
		condition container.WaitCondition) (<-chan container.WaitResponse, <-chan error)
}

// ImageAPIClient defines API client methods for the images.
type ImageAPIClient interface {
	ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error)
	ImageRemove(ctx context.Context, image string, options image.RemoveOptions) ([]image.DeleteResponse, error)
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
}

// NetworkAPIClient defines API client methods for the networks.
type NetworkAPIClient interface {
	NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error)
	NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error)
	NetworkRemove(ctx context.Context, network string) error
}

// ViperClient defines a mock Viper client for testing.
type ViperClient interface {
	SetConfigName(in string)
	SetConfigType(in string)
	AddConfigPath(in string)
	ReadInConfig() error
	SafeWriteConfig() error
	GetTime(key string) time.Time
	Set(key string, value interface{})
	WriteConfig() error
}
