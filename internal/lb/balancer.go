package lb

import (
	"hash/fnv"
	"net/http"
	"sort"
	"sync"

	"github.com/adel-hadadi/load-balancer/utils"
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

const defaultReplicas = 150

type IPHashBalancer struct {
	ring     []ringEntry
	servers  []*Server
	mu       sync.RWMutex
	replicas int
}

type ringEntry struct {
	hash   uint32
	server *Server
}

func NewIPHashBalancer() *IPHashBalancer {
	return &IPHashBalancer{
		replicas: defaultReplicas,
	}
}

func (b *IPHashBalancer) Next(r *http.Request, servers []*Server) *Server {
	if len(servers) == 0 {
		return nil
	}

	b.mu.Lock()

	if !b.serversEqual(servers) {
		b.buildRing(servers)
	}
	b.mu.Unlock()

	ip := utils.ClientIP(r)
	hash := b.hashKey(ip)

	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.ring) == 0 {
		return nil
	}

	idx := sort.Search(len(b.ring), func(i int) bool {
		return b.ring[i].hash >= hash
	})

	if idx >= len(b.ring) {
		idx = 0
	}

	return b.ring[idx].server
}

func (b *IPHashBalancer) buildRing(servers []*Server) {
	b.ring = make([]ringEntry, 0, len(servers)*b.replicas)
	b.servers = make([]*Server, len(servers))
	copy(b.servers, servers)

	h := fnv.New32a()
	for _, server := range servers {
		serverKey := server.Name()
		for i := range b.replicas {
			h.Reset()
			_, _ = h.Write([]byte(serverKey))
			_, _ = h.Write([]byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
			b.ring = append(b.ring, ringEntry{
				hash:   h.Sum32(),
				server: server,
			})
		}
	}

	sort.Slice(b.ring, func(i, j int) bool {
		return b.ring[i].hash < b.ring[j].hash
	})
}

func (b *IPHashBalancer) hashKey(key string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return h.Sum32()
}

func (b *IPHashBalancer) serversEqual(servers []*Server) bool {
	if len(b.servers) != len(servers) {
		return false
	}
	serverMap := make(map[string]bool)
	for _, s := range b.servers {
		serverMap[s.Name()] = true
	}
	for _, s := range servers {
		if !serverMap[s.Name()] {
			return false
		}
	}
	return true
}

type LeastConnectionBalancer struct {
}

func NewLeastConnectionBalancer() *LeastConnectionBalancer {
	return &LeastConnectionBalancer{}
}

func (c *LeastConnectionBalancer) Next(r *http.Request, servers []*Server) *Server {
	var chosen *Server
	for _, s := range servers {
		if chosen == nil {
			chosen = s
			continue
		}

		if chosen.activeConnections > s.activeConnections {
			chosen = s
		}
	}

	return chosen
}
