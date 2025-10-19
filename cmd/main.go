package main

import (
	"context"
	"flag"
	"github.com/adel-hadadi/load-balancer/internal/config"
	"github.com/adel-hadadi/load-balancer/internal/lb"
	"github.com/adel-hadadi/load-balancer/internal/transport"
	"log"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "/etc/loady/config.yml", "path to config file")

	flag.Parse()

	cfg, err := config.Get(configPath)
	if err != nil {
		log.Fatal(err)
	}

	srvRegistry := lb.NewServerRegistry()

	controller, err := lb.NewController(cfg, srvRegistry)
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
