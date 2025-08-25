package transport

import (
	"fmt"
	"github.com/adel-hadadi/load-balancer/internal/lb"
	"log"
	"net/http"
)

type Server struct {
	balancer *lb.BalancerProvider
}

func New(balancer *lb.BalancerProvider) *Server {
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
		sv, err := s.balancer.Next(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		if sv == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		sv.Serve(w, r)
	})

	return mux
}
