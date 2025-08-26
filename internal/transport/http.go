package transport

import (
	"fmt"
	"github.com/adel-hadadi/load-balancer/internal/lb"
	"net/http"
)

type Server struct {
	controller *lb.Controller
}

func New(controller *lb.Controller) *Server {
	return &Server{
		controller: controller,
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
		sv, err := s.controller.Next(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
