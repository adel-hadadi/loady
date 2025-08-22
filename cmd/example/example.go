package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: go run example.go [port]")
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Hello, from %s server", args[0])))
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	log.Fatal(http.ListenAndServe(":"+args[0], nil))
}
