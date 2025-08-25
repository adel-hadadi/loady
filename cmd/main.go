package main

import (
	"context"
	"github.com/adel-hadadi/load-balancer/internal/config"
	"github.com/adel-hadadi/load-balancer/internal/lb"
	"github.com/adel-hadadi/load-balancer/internal/transport"
	"log"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	registry := lb.NewServerRegistry()

	balancer := lb.New(registry)

	balancer.Healthcheck(context.Background())

	sv := transport.New(balancer)

	log.Println("load balancer served on port:", cfg.Port)
	log.Fatal(sv.Serve(cfg.Port))
}
