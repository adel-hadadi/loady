package lb

import (
	"context"
	"github.com/adel-hadadi/load-balancer/internal/config"
	"time"
)

func (p *BalancerProvider) Healthcheck(ctx context.Context) {
	go func() {
		for {
			cfg, _ := config.Get()

			select {
			case <-time.After(cfg.Healthcheck.Duration):
				for _, sv := range p.registry.GetServers() {
					sv.CheckHealth()
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
