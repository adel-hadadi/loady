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

	controller, err := lb.NewController(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = controller.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	sv := transport.New(controller)

	log.Println("load balancer served on port:", cfg.Port)
	log.Fatal(sv.Serve(cfg.Port))
}
