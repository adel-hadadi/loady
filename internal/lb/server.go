package lb

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
)

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
	// default scheme
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}

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
