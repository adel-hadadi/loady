package lb

import (
	"net/http"
	"testing"
)

// helper: create sample servers
func mockServers(n int) []*Server {
	servers := make([]*Server, n)
	for i := range servers {
		servers[i], _ = NewServer("http://localhost")
	}
	return servers
}

// helper: create dummy request
func mockRequest() *http.Request {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.RemoteAddr = "192.168.1.100:1234"
	return req
}

// ─────────────────────────────────────────────
// Benchmarks
// ─────────────────────────────────────────────

func BenchmarkRoundRobinNext(b *testing.B) {
	balancer := NewRoundRobinBalancer()
	servers := mockServers(10)
	req := mockRequest()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		balancer.Next(req, servers)
	}
}

func BenchmarkIPHashNext(b *testing.B) {
	balancer := NewIPHashBalancer()
	servers := mockServers(10)
	req := mockRequest()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		balancer.Next(req, servers)
	}
}

func BenchmarkLeastConnectionNext(b *testing.B) {
	balancer := NewLeastConnectionBalancer()
	servers := mockServers(10)
	req := mockRequest()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		balancer.Next(req, servers)
	}
}
