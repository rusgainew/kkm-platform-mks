# Nginx Reverse Proxy

Nginx reverse proxy server for the API Gateway. Provides request routing, load balancing, caching, rate limiting, and SSL/TLS termination.

## Features

- **Reverse Proxy**: Routes traffic to API Gateway (8080)
- **Rate Limiting**:
  - General API: 200 req/s burst 20
  - Metrics: 100 req/s (restricted to Docker network)
- **Gzip Compression**: Enabled for text, JSON, and media types
- **Health Check**: `/health` endpoint for monitoring
- **Keep-Alive**: Connection pooling and optimization
- **Security Headers**: X-Frame-Options, X-Content-Type-Options, HSTS, etc.
- **SSL/TLS Termination**: Full HTTPS support with HTTP/2 (enabled)
  - Self-signed certificates included for development
  - Ready for Let's Encrypt or custom certificates

## Ports

- **HTTP**: Port 80 (redirects to HTTPS)
- **HTTPS**: Port 443 (enabled with self-signed cert for dev)

## Quick Start

### HTTP → HTTPS Redirect (automatic)

All HTTP requests are automatically redirected to HTTPS:

```bash
# Automatically redirects to https://localhost/api/v1/health
curl -L -k http://localhost/api/v1/health

# Check redirect response
curl -I http://localhost/api/v1/health
# HTTP/1.1 301 Moved Permanently
# Location: https://localhost/api/v1/health
```

### HTTPS (port 443)

```bash
# Direct HTTPS access (ignore self-signed cert in dev)
curl -k https://localhost/api/v1/health
```

### Health Check (HTTP only, no redirect)

The `/health` endpoint remains available on HTTP for internal monitoring:

```bash
curl http://localhost/health
# healthy
```

## SSL/TLS Configuration

### Current Setup (Development)

Self-signed certificate included in `certs/`:

- **Certificate**: `server.crt` (valid for 365 days)
- **Key**: `server.key`
- **Domain**: localhost (self-signed)

## Configuration

### Main Settings (nginx.conf)

```nginx
# Worker processes
worker_processes auto;

# Performance
client_max_body_size 20M;
keepalive_timeout 65s;

# Gzip
gzip on;
gzip_comp_level 6;
gzip_types text/plain application/json...

# Rate Limiting
limit_req_zone $binary_remote_addr zone=general_limit:10m rate=200r/s;
```

## Enabling HTTPS

1. Obtain SSL certificates:

   ```bash
   certbot certonly --standalone -d your-domain.com
   ```

2. Copy certificates to the container:

   ```bash
   mkdir -p certs
   cp /etc/letsencrypt/live/your-domain.com/fullchain.pem certs/cert.pem
   cp /etc/letsencrypt/live/your-domain.com/privkey.pem certs/key.pem
   ```

3. Update docker-compose.prod.yml:

   ```yaml
   nginx-proxy:
     volumes:
       - ./certs:/etc/nginx/certs:ro
   ```

4. Uncomment HTTPS configuration in nginx.conf (already enabled!)

## HTTPS Status

✅ **HTTPS is ENABLED and READY to use**

Current configuration:

- **Protocol**: TLS 1.2 & 1.3
- **HTTP/2**: Enabled
- **Certificate**: Self-signed (development) in `certs/server.crt`
- **Key**: `certs/server.key`
- **HSTS**: Enabled (max-age=31536000)
- **Security Headers**: All enabled

### Test HTTPS

```bash
# Development (ignore self-signed cert warning)
curl -k https://localhost/api/v1/health

# Check certificate
openssl s_client -connect localhost:443 -showcerts 2>&1 | grep -E "subject|issuer"

# Verify HTTP/2
curl -I --http2 https://localhost/api/v1/health | grep HTTP
# Output: HTTP/2 200 (or similar)
```

### Replace Certificates for Production

```bash
# Backup old certs (if needed)
mv certs certs.backup

# Create new certs directory
mkdir -p certs

# Option 1: Copy from Let's Encrypt
cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem certs/server.crt
cp /etc/letsencrypt/live/yourdomain.com/privkey.pem certs/server.key

# Option 2: Generate new self-signed (for testing)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout certs/server.key \
  -out certs/server.crt \
  -subj "/C=RU/ST=Moscow/L=Moscow/O=KKM/CN=yourdomain.com"

# Rebuild docker image
docker-compose -f docker-compose.prod.yml up -d --build nginx-proxy
```

## Health Monitoring

Check proxy health:

```bash
# HTTP
curl http://localhost/health

# HTTPS
curl -k https://localhost/health

# Both return: healthy

# From Docker Compose
docker compose -f docker-compose.prod.yml ps nginx-proxy
```

## Metrics Access

Nginx logs are available in the container at:

- Access logs: `/var/log/nginx/access.log`
- Error logs: `/var/log/nginx/error.log`

View logs:

```bash
docker compose -f docker-compose.prod.yml logs nginx-proxy
```

## Performance Tuning

### Connection Pooling

```nginx
upstream api_gateway {
    least_conn;
    server api-gateway:8080;
    keepalive 32;
}
```

### Buffering Settings

```nginx
proxy_buffer_size 4k;
proxy_buffers 8 4k;
proxy_busy_buffers_size 8k;
```

### Timeouts

- Connect: 30s
- Send: 30s
- Read: 30s

## Testing

### Local Testing

```bash
# Health check
curl http://localhost/health

# API request
curl -X POST http://localhost/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test"}'

# Check metrics (restricted)
curl http://localhost/metrics
# Should return 403 Forbidden if accessing from outside Docker network
```

### Load Testing

```bash
# Using Apache Bench
ab -n 1000 -c 100 http://localhost/

# Using wrk
wrk -t4 -c100 -d30s http://localhost/
```

## Security Considerations

1. **Rate Limiting**: Prevents DDoS attacks
2. **Non-root User**: Nginx runs as nginx user (not root)
3. **CORS**: API Gateway handles CORS headers
4. **Input Validation**: API Gateway validates requests
5. **TLS/SSL**: Enable for production

## Docker Compose Integration

```yaml
nginx-proxy:
  build:
    context: .
    dockerfile: services/nginx-proxy/Dockerfile
  ports:
    - "80:80"
    - "443:443"
  depends_on:
    - api-gateway
  healthcheck:
    test:
      [
        "CMD",
        "wget",
        "--quiet",
        "--tries=1",
        "--spider",
        "http://localhost/health",
      ]
    interval: 10s
    timeout: 3s
    retries: 3
```

## Production Recommendations

1. Enable HTTPS with valid certificates
2. Configure domain-specific server block
3. Set up WAF (Web Application Firewall)
4. Monitor access logs for suspicious activity
5. Implement security headers (HSTS, CSP)
6. Enable HTTP/2 and HTTP/3 support
7. Configure fail-over upstream servers

## Troubleshooting

### Connection Refused

```bash
# Check if API Gateway is running
docker compose ps api-gateway

# Check logs
docker compose logs nginx-proxy api-gateway
```

### High CPU Usage

- Check worker_processes setting
- Review rate limiting zones
- Monitor upstream server health

### Slow Responses

- Check API Gateway performance
- Verify network connectivity
- Review Gzip compression settings

## References

- [Nginx Official Documentation](https://nginx.org/en/docs/)
- [Nginx Best Practices](https://nginx.org/en/docs/http/ngx_http_upstream_module.html)
- [Reverse Proxy Configuration](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)

## Performance Optimization (January 2026)

This Nginx instance has been optimized for high throughput and low latency:

### Optimizations Applied

- **Worker Processes**: Auto-detection with 65,535 file descriptor limit
- **Event Loop**: epoll with multi_accept for efficient I/O
- **Connection Limits**: 4,096 concurrent connections per worker (4x baseline)
- **File Caching**: open_file_cache with 10,000 max entries
- **Static Assets**: Browser cache (1 day) for Swagger UI and static files
- **Compression**: Gzip level 5 with minimum 1,024 byte threshold
- **Timeouts**: Optimized to 20s for faster backend failure detection
- **Connection Pooling**: keepalive_requests 100 for upstream reuse

### Performance Metrics

- **Throughput**: 19,016+ req/sec
- **Latency (p50)**: 4.44ms
- **Latency (p99)**: 24.49ms
- **Max Connections**: 4,096 per worker

### Monitoring

Real-time Nginx metrics available via:

- **Prometheus**: http://localhost:9091 (scrapes every 10s)
- **Grafana Dashboard**: http://localhost:3000/d/nginx-monitoring/
  - Credentials: admin / admin
  - Visualizes: Request rate, active connections, connection states

### Additional Documentation

- [OPTIMIZATION_REPORT.md](./OPTIMIZATION_REPORT.md) - Detailed before/after analysis with 12 sections
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Quick reference guide with key metrics
