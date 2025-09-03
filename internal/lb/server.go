package lb

import (
	"fmt"
	"github.com/adel-hadadi/load-balancer/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

var HealthyServers = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "healthy_servers_count",
	Help: "Number of servers currently marked as healthy",
})

type ServerRegistry struct {
	mu      sync.RWMutex
	urls    []string
	servers []*Server
}

func NewServerRegistry(servers []string) (*ServerRegistry, error) {
	r := &ServerRegistry{
		servers: make([]*Server, 0),
	}

	err := r.UpdateServers(servers)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *ServerRegistry) UpdateServers(urls []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	added, removed := utils.DiffSlices(r.urls, urls)

	for _, u := range added {
		s, err := NewServer(u)
		if err != nil {
			return fmt.Errorf("failed to add server %s: %v", u, err)
		}

		r.servers = append(r.servers, s)
		r.urls = append(r.urls, u)
	}

	for _, u := range removed {
		for k, server := range r.servers {
			if server.url == u {
				r.servers = append(r.servers[:k], r.servers[k+1:]...)
				r.urls = append(r.urls[:k], r.urls[k+1:]...)
			}
		}
	}

	return nil
}

func (r *ServerRegistry) All() []*Server {
	return r.servers
}

func (r *ServerRegistry) HealthyServers() []*Server {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var healthy []*Server
	for _, server := range r.servers {
		if server.IsHealthy() {
			healthy = append(healthy, server)
		}
	}

	return healthy
}

type Server struct {
	url               string
	activeConnections int64
	healthy           atomic.Bool
	rp                *httputil.ReverseProxy
}

func (s *Server) Name() string {
	return s.url
}

func (s *Server) IsHealthy() bool {
	return s.healthy.Load()
}

func (s *Server) SetHealthy(v bool) {
	if s.IsHealthy() && !v {
		HealthyServers.Dec()
	} else if !s.IsHealthy() && v {
		HealthyServers.Inc()
	}

	s.healthy.Store(v)
}

func (s *Server) CheckHealth(path string) {
	response, err := http.Get(s.url + path)
	if err != nil || response.StatusCode != 200 {
		s.SetHealthy(false)

		log.Printf("%s is unhealthy", s.url)

		return
	}

	s.SetHealthy(true)
}

func NewServer(u string) (*Server, error) {
	target, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("invalid urls: %w", err)
	}

	rp := httputil.NewSingleHostReverseProxy(target)

	s := &Server{
		url: u,
		rp:  rp,
	}

	s.healthy.Store(false)

	return s, nil
}

func (s *Server) Serve(w http.ResponseWriter, r *http.Request) http.Handler {
	atomic.AddInt64(&s.activeConnections, 1)
	defer atomic.AddInt64(&s.activeConnections, -1)

	return s.rp
}
