package lb

import (
	"github.com/adel-hadadi/load-balancer/utils"
	"hash/fnv"
	"net/http"
	"sync"
)

type Balancer interface {
	Next(r *http.Request, servers []*Server) *Server
}

type RoundRobinBalancer struct {
	Current int
	mu      sync.RWMutex
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{}
}

func (l *RoundRobinBalancer) Next(r *http.Request, servers []*Server) *Server {
	if len(servers) == 0 {
		return nil
	}

	server := servers[l.Current]

	l.Current = (l.Current + 1) % len(servers)

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

func (b *IPHashBalancer) Next(r *http.Request, servers []*Server) *Server {
	if len(servers) == 0 {
		return nil
	}

	ip := utils.ClientIP(r)

	h := fnv.New32a()
	_, _ = h.Write([]byte(ip))

	index := int(h.Sum32()) % len(servers)

	return servers[index]
}
