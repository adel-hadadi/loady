# Loady

A lightweight, fast, and extensible load balancer written in Go.

Loady helps you distribute traffic across multiple backend servers with health checks, hot-reloadable configuration, and
support for multiple balancing algorithms (round-robin, least-connections, IP hash, etc.).

---

## ✨ Features

* 🔄 **Multiple load balancing algorithms** (round-robin, least-connections, IP hash, …)
* ❤️ **Health checks** with automatic failover
* ⚙️ **Hot-reloadable configuration** (no restart required)
* 🔐 **Graceful shutdown & error handling**
* 📊 **Metrics & logging** for monitoring traffic
* 📦 **Cross-platform** binaries for Linux, macOS, and Windows

---

## 🚀 Installation

### Using Go

```bash
go install github.com/adel-hadadi/loady@latest
```

### From Release

Download a prebuilt binary from the [Releases](https://github.com/adel-hadadi/loady/releases) page.

Example for Linux:

```bash
curl -L https://github.com/adel-hadadi/loady/releases/latest/download/loady_linux_amd64 -o loady
chmod +x loady
sudo mv loady /usr/local/bin/
```

---

## ⚡ Usage

Run Loady with a simple config file:

```bash
loady --config ./config.yaml
```

> [!NOTE]
> Loady can also run without that `--config` flag but by default its looking for config file in `/etc/loady/config.yml`.

### Example `config.yaml`

```yaml
algorithm: round-robin
servers:
  - http://127.0.0.1:8081
  - http://127.0.0.1:8082
  - http://127.0.0.1:8083

healthcheck:
  interval: 5s
  timeout: 2s
  path: /health
```

---

## 🔀 Supported Algorithms

* **Round Robin** – evenly distributes requests
* **Least Connections** – sends new requests to the server with the fewest active connections
* **IP Hash** – same client IP always goes to the same backend

---

## 🩺 Health Checks

Loady continuously probes backend servers and automatically removes unhealthy nodes from the pool. When they recover,
they’re added back automatically.

Example log:

```
[INFO] server 127.0.0.1:8082 is unhealthy
[INFO] server 127.0.0.1:8082 is back online
```

---

## 📊 Observability with Prometheus

**Loady** now integrates with **Prometheus** to provide real-time metrics about your load balancer and the services it
manages. This allows you to monitor the performance and health of your applications easily.

### Metrics Exported

Some of the key metrics **Loady** exposes:

| Metric                            | Description                                    |
|-----------------------------------|------------------------------------------------|
| `http_request_duration_seconds`   | Duration of HTTP requests handled by Loady.    |
| `http_total_operations_processed` | Total number of requests processed by Loady.   |
| `healthy_servers_count`           | Number of servers currently marked as healthy. |

### How to Use

1. Run Loady with Prometheus metrics enabled (enabled by default on `/metrics` endpoint).
2. Point your Prometheus server to scrape Loady:

```yaml
scrape_configs:
  - job_name: 'loady'
    static_configs:
      - targets: [ '<LOADY_HOST>:<METRICS_PORT>' ]
```

3. Visualize metrics with Grafana or any other Prometheus-compatible dashboard.

> [!TIP]
> These metrics help you track response times, detect unhealthy servers quickly, and analyze traffic patterns for better
> load balancing decisions.

---

## 🛠 Development

Clone the repo and build:

```bash
git clone https://github.com/adel-hadadi/loady.git
cd loady
go build -o loady ./cmd
```

Run tests:

```bash
go test ./...
```

---

## 📦 Roadmap

* [ ] Add sticky sessions
* [ ] Prometheus metrics
* [ ] gRPC load balancing
* [ ] Advanced configuration via API
* [ ] More detailed logging

---

## 🤝 Contributing

Contributions are welcome! Please open an issue or PR to discuss changes.

---

## 📜 License

MIT License – see [LICENSE](./LICENSE) for details.
