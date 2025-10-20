package lb

import (
	"context"
	"errors"
	"github.com/adel-hadadi/load-balancer/utils"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"log"
)

type ServerEventType int

const (
	ServerAdded ServerEventType = iota
	ServerRemoved
)

type ServerEvent struct {
	Type    ServerEventType
	ID      string
	Address string
}

type ServerProvider interface {
	Watch(ctx context.Context, events chan<- ServerEvent) error
}

type ServerProviderOptions func(any)

func NewServerProvider(p string, opts ...ServerProviderOptions) (ServerProvider, error) {
	switch p {
	case "docker":
		return NewDockerProvider(opts...), nil
	default:
		return nil, errors.New("invalid provider. available providers are: [docker]")
	}
}

type DockerProvider struct {
}

func NewDockerProvider(...ServerProviderOptions) *DockerProvider {
	return &DockerProvider{}
}

func (p *DockerProvider) Watch(ctx context.Context, e chan<- ServerEvent) error {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer apiClient.Close()

	apiClient.NegotiateAPIVersion(context.Background())

	chanEvents, chanErrs := apiClient.Events(ctx, events.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("event", "start"),
			filters.Arg("event", "stop"),
		),
	})

	go func() {
		for {
			select {
			case event := <-chanEvents:
				c, err := apiClient.ContainerInspect(ctx, event.Actor.ID)
				if err != nil {
					log.Printf("failed to inspect container %s: %v", event.Actor.ID, err)
					return
				}

				if c.Config.Labels["loady.enabled"] == "true" {
					addr := utils.GetContainerAddress(c)

					srvEvent := ServerEvent{
						ID:      c.ID,
						Address: addr,
					}
					switch event.Action {
					case "start":
						srvEvent.Type = ServerAdded
					case "stop":
						srvEvent.Type = ServerRemoved
					}

					e <- srvEvent
				}
			case err := <-chanErrs:
				if err != nil {
					log.Printf("error on listening to docker events: %v", err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
