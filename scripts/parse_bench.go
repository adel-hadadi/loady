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

var benchRe = regexp.MustCompile(`^(Benchmark\S+)\s+(\d+)\s+([\d\.]+)\s+ns/op`)

func main() {
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
			if len(m) == 4 {
				name := m[1]
				iters, _ := strconv.ParseInt(m[2], 10, 64)
				nsPerOp, _ := strconv.ParseFloat(m[3], 64)

				fmt.Printf("Parsed: %s iters=%d ns/op=%.3f\n", name, iters, nsPerOp)

				point := influxdb2.NewPoint(
					"benchmark",
					map[string]string{"name": name},
					map[string]interface{}{
						"iterations": iters,
						"ns_per_op":  nsPerOp,
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
