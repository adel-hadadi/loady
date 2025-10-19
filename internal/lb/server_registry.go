package lb

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"sync"
)

var HealthyServers = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "healthy_servers_count",
	Help: "Number of servers currently marked as healthy",
})

type ServerRegistry struct {
	mu      sync.RWMutex
	servers map[string]*Server
}

func NewServerRegistry() *ServerRegistry {
	return &ServerRegistry{
		servers: make(map[string]*Server),
	}
}

func (r *ServerRegistry) HandleEvent(evt ServerEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	srv, err := NewServer(evt.Address)
	if err != nil {
		return err
	}

	switch evt.Type {
	case ServerAdded:
		r.servers[evt.ID] = srv
		log.Printf("server added: %s", evt.ID)
	case ServerRemoved:
		delete(r.servers, evt.ID)
		log.Printf("server removed: %v", evt.ID)
	}

	return nil
}

func (r *ServerRegistry) All() map[string]*Server {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.servers
}

func (r *ServerRegistry) HealthyServers() []*Server {
	r.mu.RLock()
	defer r.mu.RUnlock()

	servers := make([]*Server, 0)
	for _, server := range r.servers {
		if server.IsHealthy() {
			servers = append(servers, server)
		}
	}

	return servers
}
