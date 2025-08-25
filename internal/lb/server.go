package lb

import (
	"fmt"
	"github.com/adel-hadadi/load-balancer/internal/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ServerRegistry struct{}

func NewServerRegistry() *ServerRegistry {
	return &ServerRegistry{}
}

func (r *ServerRegistry) GetServers() []*Server {
	cfg, _ := config.Get()

	servers := make([]*Server, len(cfg.Servers))

	for k, raw := range cfg.Servers {
		u, err := url.Parse(raw)
		if err != nil {
			panic(err)
		}

		servers[k] = &Server{
			url:       raw,
			isHealthy: true,
			rp:        httputil.NewSingleHostReverseProxy(u),
		}
	}

	return servers
}

func (r *ServerRegistry) HealthyServers() []*Server {
	var healthy []*Server
	for _, server := range r.GetServers() {
		if server.isHealthy {
			healthy = append(healthy, server)
		}
	}

	return healthy
}

type Server struct {
	url       string
	isHealthy bool
	rp        *httputil.ReverseProxy
}

func (s *Server) CheckHealth() {
	cfg, _ := config.Get()

	response, err := http.Get(s.url + cfg.Healthcheck.Api)
	if err != nil || response.StatusCode != 200 {
		s.isHealthy = false

		log.Printf("%s is unhealthy", s.url)

		return
	}

	s.isHealthy = true
}

func NewServer(u string) (*Server, error) {
	target, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	rp := httputil.NewSingleHostReverseProxy(target)

	return &Server{
		url:       u,
		isHealthy: true,
		rp:        rp,
	}, nil
}

func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {
	s.rp.ServeHTTP(w, r)
}
