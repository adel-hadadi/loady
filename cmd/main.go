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
	// TODO: read config path from params
	var configPath string
	flag.StringVar(&configPath, "config", "/etc/loady/config.yml", "path to config file")

	flag.Parse()
	// Read the params and if it not setted read config from /etc/loady/config.yml

	cfg, err := config.Get(configPath)
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
