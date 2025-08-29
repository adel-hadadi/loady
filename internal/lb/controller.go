package lb

import (
	"context"
	"github.com/adel-hadadi/load-balancer/internal/config"
	"log"
	"net/http"
)

type Controller struct {
	cfg                *config.Config
	registry           *ServerRegistry
	checker            *HealthChecker
	balancer           Balancer
	cancelHealth       context.CancelFunc
	provider           *BalancerProvider
	healthCheckRunning bool
}

func NewController(cfg *config.Config) (*Controller, error) {
	reg, err := NewServerRegistry(cfg.Servers)
	if err != nil {
		return nil, err
	}

	provider, err := New(reg, cfg.Algorithm)
	if err != nil {
		return nil, err
	}

	return &Controller{
		registry:           reg,
		provider:           provider,
		checker:            &HealthChecker{},
		healthCheckRunning: false,
		cfg:                cfg,
	}, nil
}

func (c *Controller) Run(ctx context.Context) error {
	c.runHealthCheck()

	go func() {
		ch := c.cfg.Watch()
		for {
			select {
			case <-ctx.Done():
				return
			case newCfg := <-ch:
				log.Println("Config change detect")

				err := c.applyConfig(newCfg)
				if err != nil {
					log.Printf("failed to apply configuration: %v", err)
					return
				}
			}
		}

	}()

	return nil
}

func (c *Controller) applyConfig(cfg *config.Config) error {
	err := c.registry.UpdateServers(cfg.Servers)
	if err != nil {
		return err
	}

	c.runHealthCheck()

	err = c.provider.SetupAlgorithm(cfg.Algorithm)
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

	c.checker.Start(ctx, c.registry, c.cfg.Healthcheck.Duration, c.cfg.Healthcheck.Api)
}

func (c *Controller) Next(req *http.Request) (*Server, error) {
	healthy := c.registry.HealthyServers()
	if len(healthy) == 0 {
		return nil, nil
	}

	return c.provider.Next(req, healthy)
}
