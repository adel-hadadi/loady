package main

import (
	"fmt"
	"github.com/adel-hadadi/load-balancer/internal/lb"
	"github.com/adel-hadadi/load-balancer/internal/transport"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: loady [...servers]")
		return
	}

	balancer := lb.NewIPHashBalancer()

	for _, v := range args {
		sv, err := lb.NewServer(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		balancer.Servers = append(balancer.Servers, sv)
	}

	sv := transport.New(balancer)

	log.Fatal(sv.Serve(80))
}
