package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	file := os.Args[1]
	commit := os.Getenv("GITHUB_SHA")
	branch := os.Getenv("GITHUB_REF_NAME")

	client := influxdb2.NewClient(os.Getenv("INFLUX_URL"), os.Getenv("INFLUX_TOKEN"))
	writeAPI := client.WriteAPIBlocking(os.Getenv("INFLUX_ORG"), os.Getenv("INFLUX_BUCKET"))
	defer client.Close()

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Benchmark") {
			var name string
			var ns, bytes, allocs float64

			_, _ = fmt.Sscanf(line, "%s %*d %f ns/op %f B/op %f allocs/op", &name, &ns, &bytes, &allocs)

			p := influxdb2.NewPoint("benchmark",
				map[string]string{
					"branch": branch,
					"commit": commit,
					"name":   name,
				},
				map[string]interface{}{
					"time_ns": ns,
					"bytes":   bytes,
					"allocs":  allocs,
				},
				time.Now(),
			)

			if err := writeAPI.WritePoint(context.Background(), p); err != nil {
				panic(err)
			}
		}
	}
}
