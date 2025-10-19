package lb

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/adel-hadadi/load-balancer/internal/config"
)

type Controller struct {
	cfg                *config.Config
	registry           *ServerRegistry
	checker            *HealthChecker
	cancelHealth       context.CancelFunc
	provider           *BalancerProvider
	healthCheckRunning bool
	serverProvider     ServerProvider
}

func NewController(cfg *config.Config, serverRegistry *ServerRegistry) (*Controller, error) {
	provider, err := New(serverRegistry, cfg.Algorithm)
	if err != nil {
		return nil, err
	}

	serverProvider, err := NewServerProvider(cfg.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create server provider: %w", err)
	}

	return &Controller{
		registry:           serverRegistry,
		provider:           provider,
		checker:            &HealthChecker{},
		healthCheckRunning: false,
		cfg:                cfg,
		serverProvider:     serverProvider,
	}, nil
}

func (c *Controller) Run(ctx context.Context) error {
	c.runHealthCheck()

	evts := make(chan ServerEvent)
	err := c.serverProvider.Watch(ctx, evts)
	if err != nil {
		return fmt.Errorf("failed to watch server events: %w", err)
	}

	go func() {
		for evt := range evts {
			err := c.registry.HandleEvent(evt)
			if err != nil {
				log.Printf("failed to handle event: %v", err)
				return
			}
		}
	}()

	go func() {
		e := c.cfg.Watch()
		for {
			select {
			case newCfg := <-e:
				log.Println("Config change detect")

				err := c.applyConfig(newCfg)
				if err != nil {
					log.Printf("failed to apply configuration: %v", err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *Controller) applyConfig(cfg *config.Config) error {
	c.runHealthCheck()

	err := c.provider.SetupAlgorithm(cfg.Algorithm)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) runHealthCheck() {
	if c.healthCheckRunning {
		c.cancelHealth()
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancelHealth = cancel
	c.healthCheckRunning = true

	c.checker.Start(ctx, c.registry, c.cfg.Healthcheck.Interval, c.cfg.Healthcheck.Path)
}

func (c *Controller) Next(req *http.Request) (*Server, error) {
	healthy := c.registry.HealthyServers()
	if len(healthy) == 0 {
		return nil, nil
	}

	return c.provider.Next(req, healthy)
}
