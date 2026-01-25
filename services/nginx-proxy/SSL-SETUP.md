# SSL/TLS Certificate Configuration

## Obtaining Free Certificates with Let's Encrypt

### Using Certbot

1. **Install Certbot** (on host machine):

   ```bash
   sudo apt-get install certbot python3-certbot-nginx
   ```

2. **Generate Certificates**:

   ```bash
   sudo certbot certonly --standalone \
     -d your-domain.com \
     -d www.your-domain.com \
     --email admin@your-domain.com \
     --agree-tos \
     --non-interactive
   ```

3. **Copy to Project**:
   ```bash
   mkdir -p certs
   sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem certs/cert.pem
   sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem certs/key.pem
   sudo chown $(whoami):$(whoami) certs/*
   ```

### Using Docker Certbot

```bash
# Generate certificates using Docker
docker run -it --rm -v $(pwd)/certs:/etc/letsencrypt certbot/certbot certonly \
  --standalone \
  -d your-domain.com \
  -d www.your-domain.com \
  --email admin@your-domain.com \
  --agree-tos \
  --non-interactive
```

## Self-Signed Certificate (Development Only)

```bash
# Generate self-signed certificate
mkdir -p certs
openssl req -x509 -newkey rsa:4096 -keyout certs/key.pem -out certs/cert.pem -days 365 -nodes \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
```

## Docker Compose Volume Mapping

Update `docker-compose.prod.yml`:

```yaml
nginx-proxy:
  volumes:
    - ./certs:/etc/nginx/certs:ro
```

## Certificate Renewal

### Automatic Renewal (Host)

1. **Create Renewal Script** (`renew-certs.sh`):

   ```bash
   #!/bin/bash
   certbot renew --quiet
   cp /etc/letsencrypt/live/your-domain.com/fullchain.pem ~/project/certs/cert.pem
   cp /etc/letsencrypt/live/your-domain.com/privkey.pem ~/project/certs/key.pem
   docker compose -f docker-compose.prod.yml restart nginx-proxy
   ```

2. **Add to Crontab**:
   ```bash
   crontab -e
   # Add line:
   0 0 1 * * /path/to/renew-certs.sh
   ```

### Manual Renewal

```bash
# Check renewal status
sudo certbot renew --dry-run

# Renew certificates
sudo certbot renew

# Restart Nginx
docker compose -f docker-compose.prod.yml restart nginx-proxy
```

## Enabling HTTPS in nginx.conf

Uncomment and configure the HTTPS server block:

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com www.your-domain.com;

    ssl_certificate /etc/nginx/certs/cert.pem;
    ssl_certificate_key /etc/nginx/certs/key.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    location / {
        proxy_pass http://api_gateway;
        # ... proxy configuration
    }
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;
    return 301 https://$host$request_uri;
}
```

## Testing SSL Configuration

```bash
# Test certificate validity
openssl x509 -in certs/cert.pem -text -noout

# Test with curl
curl -v --cacert certs/cert.pem https://localhost/

# Check certificate expiration
openssl x509 -enddate -noout -in certs/cert.pem

# Test SSL configuration
openssl s_client -connect localhost:443 -showcerts
```

## Monitoring Certificate Expiration

### Using SSL Certificate Checker

```bash
# Check certificate expiration date
echo | openssl s_client -servername your-domain.com -connect localhost:443 2>/dev/null | openssl x509 -noout -dates

# Days until expiration
echo | openssl s_client -servername your-domain.com -connect localhost:443 2>/dev/null | openssl x509 -noout -dates | grep notAfter
```

### Monitoring Dashboard

Add certificate expiration check to Prometheus:

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "ssl_cert"
    metrics_path: "/probe"
    static_configs:
      - targets: ["https://your-domain.com"]
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: ssl-exporter:9219
```

## Common Issues

### Certificate Permission Denied

```bash
# Fix permissions
chmod 644 certs/cert.pem
chmod 600 certs/key.pem
chown nginx:nginx certs/*
```

### ACME Challenge Failed

Ensure port 80 is accessible:

```bash
# Check if port 80 is open
sudo netstat -tlnp | grep :80
```

### Certificate Already in Use

```bash
# Revoke old certificate
sudo certbot revoke --cert-path /etc/letsencrypt/live/your-domain.com/cert.pem

# Delete old certificate
sudo rm -rf /etc/letsencrypt/live/your-domain.com/
```
