package transport

import (
	"fmt"
	"github.com/adel-hadadi/load-balancer/internal/lb"
	"net/http"
)

type Server struct {
	balancer lb.Balancer
}

func New(balancer lb.Balancer) *Server {
	return &Server{
		balancer: balancer,
	}
}

func (s *Server) Serve(port int) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.routes(),
	}

	return srv.ListenAndServe()
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sv := s.balancer.Next(r)
		if sv == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		sv.Serve(w, r)
	})

	return mux
}
