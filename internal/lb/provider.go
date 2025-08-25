package lb

import (
	"errors"
	"fmt"
	"github.com/adel-hadadi/load-balancer/internal/config"
	"net/http"
)

type BalancerProvider struct {
	registry *ServerRegistry
	balancer Balancer
	algo     string
}

func New(registry *ServerRegistry) *BalancerProvider {
	return &BalancerProvider{
		registry: registry,
	}
}

func (p *BalancerProvider) Next(r *http.Request) (*Server, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("error on reading config: %w", err)
	}

	if p.algo != cfg.Algorithm {
		if err := p.refreshAlgorithm(cfg.Algorithm); err != nil {
			return nil, err
		}
	}

	servers := p.registry.HealthyServers()

	return p.balancer.Next(r, servers), nil
}

func (p *BalancerProvider) refreshAlgorithm(algo string) error {
	switch algo {
	case "round-robin":
		p.balancer = NewLoadBalancer()
	case "iphash":
		p.balancer = NewIPHashBalancer()
	default:
		return errors.New("unknown algorithm, available balancers: round-robin, iphash)")
	}

	p.algo = algo

	return nil
}
