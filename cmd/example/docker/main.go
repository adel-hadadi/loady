package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

// getContainerIP returns the first non-loopback IPv4 address
func getContainerIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "unknown"
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && ip.To4() != nil {
				return ip.String()
			}
		}
	}
	return "unknown"
}

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	ip := getContainerIP()

	fmt.Fprintf(w, "===== Container Info =====\n")
	fmt.Fprintf(w, "Hostname: %s\n", hostname)
	fmt.Fprintf(w, "Container IP: %s\n", ip)

	fmt.Fprintf(w, "\n===== Request Info =====\n")
	fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "URL: %s\n", r.URL.String())
	fmt.Fprintf(w, "Protocol: %s\n", r.Proto)

	fmt.Fprintf(w, "\nHeaders:\n")
	for name, values := range r.Header {
		fmt.Fprintf(w, "  %s: %s\n", name, strings.Join(values, ", "))
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Starting server on :80 ...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
