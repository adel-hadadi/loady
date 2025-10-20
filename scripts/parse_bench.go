package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var benchRe = regexp.MustCompile(
	`^(Benchmark\S+)\s+(\d+)\s+([\d\.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op`,
)

func main() {
	commit := os.Getenv("GITHUB_SHA")
	branch := os.Getenv("GITHUB_REF_NAME")

	file, err := os.Open("bench.txt") // file with `go test -bench=. -benchmem` output
	if err != nil {
		panic(err)
	}
	defer file.Close()

	url := os.Getenv("INFLUX_URL")
	token := os.Getenv("INFLUX_TOKEN")
	org := os.Getenv("INFLUX_ORG")
	bucket := os.Getenv("INFLUX_BUCKET")

	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPIBlocking(org, bucket)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Benchmark") {
			m := benchRe.FindStringSubmatch(line)
			if len(m) == 6 {
				name := m[1]
				iters, _ := strconv.ParseInt(m[2], 10, 64)
				nsPerOp, _ := strconv.ParseFloat(m[3], 64)
				bPerOp, _ := strconv.ParseInt(m[4], 10, 64)
				allocsPerOp, _ := strconv.ParseInt(m[5], 10, 64)

				fmt.Printf("✅ Parsed: %-40s | iters=%d ns/op=%.3f B/op=%d allocs/op=%d\n",
					name, iters, nsPerOp, bPerOp, allocsPerOp)

				point := influxdb2.NewPoint(
					"benchmark",
					map[string]string{
						"name":   name,
						"branch": branch,
						"commit": commit,
					},
					map[string]interface{}{
						"iterations": iters,
						"time_ns":    nsPerOp,
						"bytes":      bPerOp,
						"allocs":     allocsPerOp,
					},
					time.Now(),
				)
				if err := writeAPI.WritePoint(context.Background(), point); err != nil {
					fmt.Println("write error:", err)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	client.Close()
	fmt.Println("✅ done writing benchmarks to InfluxDB")
}
