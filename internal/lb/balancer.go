package lb

import (
	"github.com/adel-hadadi/load-balancer/utils"
	"hash/fnv"
	"net/http"
	"sync"
)

type Balancer interface {
	Next(r *http.Request) *Server
}

type RoundRobinBalancer struct {
	Servers []*Server
	Current int
	mu      sync.RWMutex
}

func NewLoadBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{
		Servers: make([]*Server, 0),
	}
}

func (l *RoundRobinBalancer) Next(r *http.Request) *Server {
	if len(l.Servers) == 0 {
		return nil
	}

	l.mu.Lock()

	server := l.Servers[l.Current]

	l.Current = (l.Current + 1) % len(l.Servers)

	l.mu.Unlock()

	if !server.isHealthy {
		return l.Next(r)
	}

	return server
}

type IPHashBalancer struct {
	Servers []*Server
}

func NewIPHashBalancer() *IPHashBalancer {
	return &IPHashBalancer{
		Servers: make([]*Server, 0),
	}
}

func (b *IPHashBalancer) Next(r *http.Request) *Server {
	if len(b.Servers) == 0 {
		return nil
	}

	ip := utils.ClientIP(r)

	h := fnv.New32a()
	_, _ = h.Write([]byte(ip))

	index := int(h.Sum32()) % len(b.Servers)

	// What if server is unhealthy?

	return b.Servers[index]
}
