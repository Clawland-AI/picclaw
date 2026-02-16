# Health Check Endpoint

PicoClaw provides an HTTP health check endpoint for monitoring and orchestration.

## Endpoint

```
GET /healthz
```

## Configuration

Add to `config.json`:

```json
{
  "gateway": {
    "host": "0.0.0.0",
    "port": 18790,
    "health_port": 9090
  }
}
```

Or via environment variable:

```bash
export PICOCLAW_GATEWAY_HEALTH_PORT=9090
```

## Response

```json
{
  "status": "ok",
  "uptime": "1h23m45s",
  "version": "0.1.0",
  "agent_name": "picclaw",
  "go_version": "go1.24.0",
  "build_time": "2026-02-16T08:00:00Z"
}
```

**HTTP 200 OK** - Agent is healthy and running

## Usage

### Local Development

```bash
# Start the gateway
picoclaw gateway

# Check health
curl http://localhost:9090/healthz
```

### Docker

```bash
# Build
docker build -t picclaw .

# Run with health check
docker run -d \
  --name picclaw \
  --health-cmd "curl -f http://localhost:9090/healthz || exit 1" \
  --health-interval 30s \
  --health-timeout 5s \
  --health-retries 3 \
  picclaw
```

### Kubernetes

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: picclaw
spec:
  containers:
  - name: picclaw
    image: picclaw:latest
    ports:
    - containerPort: 9090
      name: health
    livenessProbe:
      httpGet:
        path: /healthz
        port: 9090
      initialDelaySeconds: 10
      periodSeconds: 30
    readinessProbe:
      httpGet:
        path: /healthz
        port: 9090
      initialDelaySeconds: 5
      periodSeconds: 10
```

### systemd

```ini
[Unit]
Description=PicoClaw Agent
After=network.target

[Service]
ExecStart=/usr/local/bin/picoclaw gateway
ExecStartPost=/bin/sleep 2
ExecStartPost=/usr/bin/curl -f http://localhost:9090/healthz
Restart=always

[Install]
WantedBy=multi-user.target
```

## Testing

```bash
# Run tests
go test -v ./pkg/health

# With coverage
go test -v -cover ./pkg/health

# Expected output:
# === RUN   TestNewServer
# --- PASS: TestNewServer (0.00s)
# === RUN   TestHealthzEndpoint
# --- PASS: TestHealthzEndpoint (0.00s)
# ...
# PASS
# coverage: 100.0% of statements
```

## Monitoring

### Prometheus

```yaml
scrape_configs:
  - job_name: 'picclaw'
    metrics_path: '/healthz'
    static_configs:
      - targets: ['localhost:9090']
```

### Uptime Monitoring

Services like UptimeRobot, Pingdom, or Datadog can use `/healthz` for availability checks.

## Troubleshooting

### Port already in use

```bash
# Change port in config.json
{
  "gateway": {
    "health_port": 9091
  }
}
```

### Health check fails

1. Check if gateway is running: `ps aux | grep picoclaw`
2. Check port binding: `netstat -tuln | grep 9090`
3. Check logs: `tail -f ~/.picoclaw/picoclaw.log`
4. Verify config: `cat ~/.picoclaw/config.json`

### Method not allowed (405)

Only `GET` requests are supported:

```bash
# ✅ Correct
curl -X GET http://localhost:9090/healthz

# ❌ Wrong
curl -X POST http://localhost:9090/healthz
```
