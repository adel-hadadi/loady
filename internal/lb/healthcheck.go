package lb

import (
	"context"
	"time"
)

type HealthChecker struct {
}

func (h *HealthChecker) Start(ctx context.Context, reg *ServerRegistry, interval time.Duration, path string) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, s := range reg.All() {
					go s.CheckHealth(path)
				}
			}
		}
	}()
}
