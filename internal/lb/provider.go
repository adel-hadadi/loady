package lb

import (
	"errors"
	"net/http"
)

type BalancerProvider struct {
	registry *ServerRegistry
	balancer Balancer
}

func New(registry *ServerRegistry, algo string) (*BalancerProvider, error) {
	p := &BalancerProvider{
		registry: registry,
	}

	err := p.SetupAlgorithm(algo)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *BalancerProvider) Next(r *http.Request, servers []*Server) (*Server, error) {
	return p.balancer.Next(r, servers), nil
}

func (p *BalancerProvider) SetupAlgorithm(algo string) error {
	switch algo {
	case "round-robin":
		p.balancer = NewRoundRobinBalancer()
	case "iphash":
		p.balancer = NewIPHashBalancer()
	default:
		return errors.New("unknown algorithm, available balancers: round-robin, iphash)")
	}

	return nil
}
