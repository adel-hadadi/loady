package lb

import (
	"context"
	"time"
)

func (l *RoundRobinBalancer) healthcheck(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				for {
					time.Sleep(5 * time.Second)

					for _, sv := range l.Servers {
						sv.CheckHealth()
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
