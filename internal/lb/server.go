package lb

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	url       string
	isHealthy bool
	rp        *httputil.ReverseProxy
}

func (s *Server) CheckHealth() {
	response, err := http.Get(s.url + "/health")
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
