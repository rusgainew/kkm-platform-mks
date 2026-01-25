# Nginx Reverse Proxy for API Gateway
# Multi-stage build: minimal runtime image

FROM nginx:1.25-alpine

# Remove default nginx config
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom nginx configuration
COPY services/nginx-proxy/nginx.conf /etc/nginx/nginx.conf

# Copy SSL certificates (for development/testing)
RUN mkdir -p /etc/nginx/certs
COPY services/nginx-proxy/certs/ /etc/nginx/certs/

# Create log directories
RUN mkdir -p /var/log/nginx && \
  chown -R nginx:nginx /var/log/nginx && \
  chown -R nginx:nginx /var/cache/nginx && \
  chown -R nginx:nginx /etc/nginx && \
  chmod 700 /etc/nginx/certs && \
  chmod 600 /etc/nginx/certs/*

# Create non-root user (nginx user already exists in image)
# Ensure proper permissions
RUN touch /var/run/nginx.pid && \
  chown nginx:nginx /var/run/nginx.pid

# Expose HTTP and HTTPS ports
EXPOSE 80 443

# Health check
HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost/health || exit 1

# Switch to non-root user
USER nginx

# Start nginx in foreground
CMD ["nginx", "-g", "daemon off;"]
