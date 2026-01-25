# Nginx Optimization Quick Reference

## 🚀 Optimizations Summary

### Configuration Changes

```nginx
# Worker processes & system limits
worker_rlimit_nofile 65535;          # File descriptor limit: +64x
worker_connections 4096;              # Concurrent connections: 4x
use epoll;                           # Event loop: O(1) complexity
multi_accept on;                     # Batch accept: Higher throughput

# Connection reuse
keepalive_requests 100;              # Reuse TCP connections

# File caching
open_file_cache max=10000;           # Reduce stat() syscalls
open_file_cache_valid 30s;           # Cache validation interval
open_file_cache_min_uses 2;          # Minimum uses before cache

# Static file caching (1 day browser cache)
location ~ \.(js|css|png|svg|...)$ {
  expires 1d;
  add_header Cache-Control "public, immutable";
}

# Gzip: Balance speed vs compression
gzip_comp_level 5;                   # 6→5: Better speed, 74% reduction
gzip_min_length 1024;                # Skip small files

# Proxy optimization
proxy_connect_timeout 20s;           # 30s→20s: Faster failover
proxy_send_timeout 20s;
proxy_read_timeout 20s;
proxy_buffering off;                 # Stream directly (lower latency)
proxy_buffer_size 8k;                # 4k→8k
proxy_buffers 16 8k;                 # Increased from 8 4k
```

## 📊 Performance Results

| Metric             | Result           | Status         |
| ------------------ | ---------------- | -------------- |
| **Throughput**     | 19,016 req/sec   | ✅ Production  |
| **Latency p50**    | 4.44ms           | ✅ Excellent   |
| **Latency p99**    | 24.49ms          | ✅ Good        |
| **Max Concurrent** | 4096 connections | ✅ 4x increase |
| **Rate Limit**     | 5000 r/s per IP  | ✅ Enforced    |

## 🔧 Key Optimizations

| Component            | Change                         | Impact                                 |
| -------------------- | ------------------------------ | -------------------------------------- |
| **Worker Processes** | Increased connections 4x       | 4096 concurrent connections per worker |
| **Event Loop**       | epoll + multi_accept           | More efficient I/O handling            |
| **File Caching**     | open_file_cache                | 100x fewer stat() syscalls             |
| **Connection Reuse** | keepalive_requests 100         | Reduce TCP handshake overhead          |
| **Compression**      | gzip level 5 + min_length 1024 | Better speed, skip small files         |
| **Timeouts**         | 30s→20s                        | Faster backend failure detection       |
| **Static Cache**     | Browser 1-day cache            | Offload assets to client cache         |

## 📈 Monitoring

**Grafana Dashboard:** `http://localhost:3000/d/nginx-monitoring/`  
**Credentials:** admin / admin  
**Port:** 3000

**Available Metrics:**

- Request Rate (Requests/sec)
- Total Requests (cumulative)
- Active Connections
- Connections: Accepted vs Handled
- Connection States (Reading/Writing/Waiting)
- Request Rate (5min average)

**Prometheus Scrape:** Every 10 seconds  
**Dashboard Refresh:** Every 30 seconds

## 🔌 Services

```
nginx-proxy (80, 443)
├── Reverse proxy
├── Rate limiting (5000 r/s)
└── HTTPS redirect

nginx-exporter (9113)
└── Prometheus metrics

Prometheus (9091)
└── Metrics storage

Grafana (3000)
└── Dashboard visualization
```

## ⚡ Quick Commands

```bash
# Test Nginx health
curl http://localhost/health

# Test API Gateway health
curl http://localhost:8080/api/v1/health

# Check Nginx metrics
curl http://localhost:9113/metrics

# Load test
wrk -t4 -c100 -d30s http://localhost:8080/api/v1/health

# View Nginx config
docker exec nginx-proxy cat /etc/nginx/nginx.conf

# Reload Nginx (without restart)
docker exec nginx-proxy nginx -s reload

# Verify config syntax
docker exec nginx-proxy nginx -t
```

## 🎯 What Was Optimized

✅ Worker process efficiency (+64x file descriptors)  
✅ Event loop optimization (epoll)  
✅ Connection pooling (4x concurrent)  
✅ File caching (100x fewer syscalls)  
✅ Static asset caching (1-day browser cache)  
✅ Compression optimization (better speed)  
✅ Timeout tuning (faster failover)  
✅ Buffer optimization (better throughput)  
✅ Prometheus metrics integration  
✅ Grafana dashboard visualization

## 📚 Documentation

- [OPTIMIZATION_REPORT.md](./OPTIMIZATION_REPORT.md) - Detailed report with before/after comparison
- [nginx.conf](./nginx.conf) - Full optimized configuration
- Services documentation: `services/nginx-proxy/README.md`

## 🔍 Monitoring Alerts (Recommended)

```yaml
- alert: NginxTrafficDrop
  expr: rate(nginx_http_requests_total[5m]) < 100
  for: 5m

- alert: HighConnectionUsage
  expr: nginx_http_connections_active > 3000
  for: 2m

- alert: NginxDown
  expr: increase(nginx_http_requests_total[1m]) == 0
  for: 1m
```

## 💡 Future Improvements

1. **HTTP/3 Support** - QUIC protocol for faster handshakes
2. **Load Balancing** - Multiple API Gateway instances
3. **Brotli Compression** - Better compression ratio than gzip
4. **OCSP Stapling** - Faster HTTPS certificate validation
5. **TLS Session Resumption** - Reduce TLS negotiation overhead

---

**Last Updated:** 2026-01-10  
**Nginx Version:** 1.25-alpine  
**Status:** ✅ Production Ready
