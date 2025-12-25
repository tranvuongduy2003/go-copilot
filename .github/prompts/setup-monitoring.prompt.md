---
description: Set up a complete monitoring and observability stack with Prometheus, Grafana, Loki, and Tempo
---

# Setup Monitoring Stack

Set up a comprehensive observability stack for metrics, logs, and traces following best practices.

## Stack Components

**Metrics**: Prometheus
**Visualization**: Grafana
**Logging**: Loki + Promtail
**Tracing**: Tempo
**Alerting**: Alertmanager

## Requirements

### Metrics (Prometheus)

- [ ] Application metrics scraping
- [ ] Node exporter for host metrics
- [ ] cAdvisor for container metrics
- [ ] Custom alert rules
- [ ] Recording rules for performance

### Logging (Loki)

- [ ] Docker/Kubernetes log collection
- [ ] JSON log parsing
- [ ] Label extraction
- [ ] Log retention policy

### Tracing (Tempo)

- [ ] OpenTelemetry support
- [ ] Trace-to-log correlation
- [ ] Trace-to-metric correlation
- [ ] Service graph generation

### Alerting

- [ ] Critical alerts to PagerDuty
- [ ] Warning alerts to Slack
- [ ] Email for daily summaries
- [ ] Runbook links in alerts

### Dashboards

- [ ] Service overview dashboard
- [ ] Infrastructure dashboard
- [ ] SLO dashboard
- [ ] On-call dashboard

## Implementation

### 1. Docker Compose

Location: `monitoring/docker-compose.monitoring.yml`

```yaml
version: "3.9"

services:
  prometheus:
    image: prom/prometheus:v2.48.0
    container_name: prometheus
    volumes:
      - ./prometheus:/etc/prometheus:ro
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.retention.time=15d'
      - '--web.enable-lifecycle'
    ports:
      - "9090:9090"
    networks:
      - monitoring

  grafana:
    image: grafana/grafana:10.2.2
    container_name: grafana
    volumes:
      - ./grafana/provisioning:/etc/grafana/provisioning:ro
      - grafana_data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
    ports:
      - "3001:3000"
    networks:
      - monitoring

  alertmanager:
    image: prom/alertmanager:v0.26.0
    container_name: alertmanager
    volumes:
      - ./alertmanager:/etc/alertmanager:ro
    ports:
      - "9093:9093"
    networks:
      - monitoring

  loki:
    image: grafana/loki:2.9.2
    container_name: loki
    volumes:
      - ./loki:/etc/loki:ro
      - loki_data:/loki
    command: -config.file=/etc/loki/loki-config.yml
    ports:
      - "3100:3100"
    networks:
      - monitoring

  promtail:
    image: grafana/promtail:2.9.2
    container_name: promtail
    volumes:
      - ./promtail:/etc/promtail:ro
      - /var/log:/var/log:ro
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
    command: -config.file=/etc/promtail/promtail-config.yml
    networks:
      - monitoring

  tempo:
    image: grafana/tempo:2.3.1
    container_name: tempo
    volumes:
      - ./tempo:/etc/tempo:ro
      - tempo_data:/tmp/tempo
    command: -config.file=/etc/tempo/tempo-config.yml
    ports:
      - "3200:3200"
      - "4317:4317"
      - "4318:4318"
    networks:
      - monitoring

networks:
  monitoring:
    driver: bridge

volumes:
  prometheus_data:
  grafana_data:
  loki_data:
  tempo_data:
```

### 2. Prometheus Configuration

Location: `monitoring/prometheus/prometheus.yml`

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']

rule_files:
  - /etc/prometheus/alerts/*.yml

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'backend'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: '/metrics'

  - job_name: 'frontend'
    static_configs:
      - targets: ['frontend:3000']
    metrics_path: '/api/metrics'
```

### 3. Alert Rules

Location: `monitoring/prometheus/alerts/application.yml`

```yaml
groups:
  - name: application
    rules:
      - alert: HighErrorRate
        expr: |
          (sum(rate(http_requests_total{status=~"5.."}[5m]))
          / sum(rate(http_requests_total[5m]))) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          runbook_url: "https://runbooks.example.com/high-error-rate"

      - alert: HighLatency
        expr: |
          histogram_quantile(0.95,
            sum(rate(http_request_duration_seconds_bucket[5m])) by (le)
          ) > 2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"

      - alert: ServiceDown
        expr: up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service {{ $labels.job }} is down"
```

### 4. Alertmanager Configuration

Location: `monitoring/alertmanager/alertmanager.yml`

```yaml
global:
  slack_api_url: '${SLACK_WEBHOOK_URL}'

route:
  receiver: 'slack-notifications'
  group_by: ['alertname', 'severity']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h

  routes:
    - match:
        severity: critical
      receiver: 'pagerduty-critical'
    - match:
        severity: warning
      receiver: 'slack-warnings'

receivers:
  - name: 'slack-notifications'
    slack_configs:
      - channel: '#alerts'
        send_resolved: true

  - name: 'slack-warnings'
    slack_configs:
      - channel: '#alerts-warnings'
        send_resolved: true

  - name: 'pagerduty-critical'
    pagerduty_configs:
      - service_key: '${PAGERDUTY_KEY}'
```

### 5. Grafana Datasources

Location: `monitoring/grafana/provisioning/datasources/datasources.yml`

```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    url: http://prometheus:9090
    isDefault: true

  - name: Loki
    type: loki
    url: http://loki:3100
    jsonData:
      derivedFields:
        - datasourceUid: tempo
          matcherRegex: "traceId=(\\w+)"
          name: TraceID
          url: "$${__value.raw}"

  - name: Tempo
    type: tempo
    url: http://tempo:3200
    jsonData:
      tracesToLogs:
        datasourceUid: loki
```

## Application Integration

### Node.js Metrics

```javascript
const promClient = require('prom-client');

// Enable default metrics
promClient.collectDefaultMetrics();

// Custom metrics
const httpRequestDuration = new promClient.Histogram({
  name: 'http_request_duration_seconds',
  help: 'Duration of HTTP requests in seconds',
  labelNames: ['method', 'route', 'status'],
  buckets: [0.01, 0.05, 0.1, 0.5, 1, 2, 5]
});

// Metrics endpoint
app.get('/metrics', async (req, res) => {
  res.set('Content-Type', promClient.register.contentType);
  res.end(await promClient.register.metrics());
});
```

### OpenTelemetry Tracing

```javascript
const { NodeSDK } = require('@opentelemetry/sdk-node');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-http');

const sdk = new NodeSDK({
  traceExporter: new OTLPTraceExporter({
    url: 'http://tempo:4318/v1/traces'
  }),
  serviceName: 'backend'
});

sdk.start();
```

## Validation

After implementation:

```bash
# Start monitoring stack
make monitoring-up

# Verify services
curl http://localhost:9090/-/healthy  # Prometheus
curl http://localhost:3001/api/health  # Grafana
curl http://localhost:3100/ready       # Loki
curl http://localhost:3200/ready       # Tempo

# Check targets
curl http://localhost:9090/api/v1/targets

# Check alerts
curl http://localhost:9090/api/v1/alerts
```

## Access URLs

| Service | URL | Credentials |
|---------|-----|-------------|
| Grafana | http://localhost:3001 | admin/admin |
| Prometheus | http://localhost:9090 | - |
| Alertmanager | http://localhost:9093 | - |

## Output

Provide:
1. Docker Compose configuration
2. Prometheus configuration with scrape jobs
3. Alert rules for key metrics
4. Alertmanager routing configuration
5. Grafana datasource provisioning
6. Application integration code
7. Dashboard JSON (if applicable)
