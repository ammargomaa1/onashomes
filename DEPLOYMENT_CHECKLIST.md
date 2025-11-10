# Deployment Checklist

## 🚀 Pre-Deployment Checklist

### Environment Configuration
- [ ] Copy `.env.example` to `.env`
- [ ] Set strong `JWT_SECRET` (min 32 characters)
- [ ] Set strong `JWT_REFRESH_SECRET` (different from JWT_SECRET)
- [ ] Configure database credentials
- [ ] Set `GIN_MODE=release` for production
- [ ] Configure appropriate `JWT_EXPIRY` (recommended: 900-3600 seconds)
- [ ] Configure appropriate `JWT_REFRESH_EXPIRY` (recommended: 604800 seconds)
- [ ] Set appropriate `DEFAULT_PAGE_SIZE` and `MAX_PAGE_SIZE`

### Database Setup
- [ ] Create production database
- [ ] Verify database connection
- [ ] Run migrations (auto-migration on startup)
- [ ] Verify seed data created
- [ ] Create database backups schedule
- [ ] Set up database monitoring

### Security
- [ ] Change all default passwords
- [ ] Generate strong JWT secrets
- [ ] Enable HTTPS/TLS
- [ ] Configure CORS for specific origins (not *)
- [ ] Set up rate limiting (implement if needed)
- [ ] Enable database SSL mode (`DB_SSLMODE=require`)
- [ ] Review and restrict database user permissions
- [ ] Set up firewall rules
- [ ] Enable audit logging

### Application
- [ ] Build production binary: `make build-linux`
- [ ] Test binary execution
- [ ] Set up process manager (systemd, supervisor, PM2)
- [ ] Configure log rotation
- [ ] Set up health check monitoring
- [ ] Configure reverse proxy (nginx, Apache)
- [ ] Set up SSL certificates (Let's Encrypt)

### Testing
- [ ] Test user registration
- [ ] Test user login
- [ ] Test admin login
- [ ] Test permission system
- [ ] Test pagination
- [ ] Test token refresh
- [ ] Load testing
- [ ] Security testing

## 📋 Production Environment Setup

### 1. System Requirements
```bash
# Minimum requirements
- CPU: 2 cores
- RAM: 2GB
- Disk: 20GB SSD
- OS: Ubuntu 20.04+ or similar

# Recommended
- CPU: 4+ cores
- RAM: 4GB+
- Disk: 50GB+ SSD
- OS: Ubuntu 22.04 LTS
```

### 2. Install Dependencies
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y

# Install Go (if building on server)
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install nginx (reverse proxy)
sudo apt install nginx -y

# Install certbot (SSL)
sudo apt install certbot python3-certbot-nginx -y
```

### 3. Database Setup
```bash
# Create database user
sudo -u postgres psql
CREATE USER ecommerce_user WITH PASSWORD 'your_secure_password';
CREATE DATABASE ecommerce_db OWNER ecommerce_user;
GRANT ALL PRIVILEGES ON DATABASE ecommerce_db TO ecommerce_user;
\q

# Configure PostgreSQL for remote connections (if needed)
sudo nano /etc/postgresql/14/main/postgresql.conf
# Set: listen_addresses = 'localhost'

sudo nano /etc/postgresql/14/main/pg_hba.conf
# Add: host ecommerce_db ecommerce_user 127.0.0.1/32 md5

sudo systemctl restart postgresql
```

### 4. Application Deployment
```bash
# Create application directory
sudo mkdir -p /opt/ecommerce-api
sudo chown $USER:$USER /opt/ecommerce-api

# Copy files
cd /opt/ecommerce-api
# Upload your binary and .env file

# Make binary executable
chmod +x ecommerce-api-linux-amd64

# Test run
./ecommerce-api-linux-amd64
```

### 5. Systemd Service Setup
```bash
# Create service file
sudo nano /etc/systemd/system/ecommerce-api.service
```

```ini
[Unit]
Description=E-Commerce API Service
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/ecommerce-api
ExecStart=/opt/ecommerce-api/ecommerce-api-linux-amd64
Restart=always
RestartSec=5
StandardOutput=append:/var/log/ecommerce-api/output.log
StandardError=append:/var/log/ecommerce-api/error.log

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/ecommerce-api

[Install]
WantedBy=multi-user.target
```

```bash
# Create log directory
sudo mkdir -p /var/log/ecommerce-api
sudo chown www-data:www-data /var/log/ecommerce-api

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable ecommerce-api
sudo systemctl start ecommerce-api
sudo systemctl status ecommerce-api
```

### 6. Nginx Configuration
```bash
sudo nano /etc/nginx/sites-available/ecommerce-api
```

```nginx
upstream api_backend {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name api.yourdomain.com;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
    limit_req zone=api_limit burst=20 nodelay;

    location / {
        proxy_pass http://api_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://api_backend/health;
        access_log off;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/ecommerce-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# Setup SSL with Let's Encrypt
sudo certbot --nginx -d api.yourdomain.com
```

### 7. Firewall Configuration
```bash
# UFW setup
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 80/tcp      # HTTP
sudo ufw allow 443/tcp     # HTTPS
sudo ufw enable
sudo ufw status
```

## 🔍 Monitoring Setup

### 1. Health Check Monitoring
```bash
# Create monitoring script
sudo nano /opt/ecommerce-api/health-check.sh
```

```bash
#!/bin/bash
HEALTH_URL="http://localhost:8080/health"
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $HEALTH_URL)

if [ $RESPONSE -ne 200 ]; then
    echo "API health check failed with status: $RESPONSE"
    # Send alert (email, Slack, etc.)
    systemctl restart ecommerce-api
fi
```

```bash
chmod +x /opt/ecommerce-api/health-check.sh

# Add to crontab (every 5 minutes)
crontab -e
*/5 * * * * /opt/ecommerce-api/health-check.sh
```

### 2. Log Monitoring
```bash
# View logs
sudo journalctl -u ecommerce-api -f

# View application logs
tail -f /var/log/ecommerce-api/output.log
tail -f /var/log/ecommerce-api/error.log

# Setup log rotation
sudo nano /etc/logrotate.d/ecommerce-api
```

```
/var/log/ecommerce-api/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 www-data www-data
    sharedscripts
    postrotate
        systemctl reload ecommerce-api > /dev/null 2>&1 || true
    endscript
}
```

## 🔄 Deployment Process

### Initial Deployment
```bash
1. Prepare server environment
2. Install dependencies
3. Setup database
4. Upload application files
5. Configure environment variables
6. Setup systemd service
7. Configure nginx
8. Setup SSL
9. Start services
10. Verify deployment
```

### Updates/Redeployment
```bash
# 1. Build new version
make build-linux

# 2. Backup current version
sudo cp /opt/ecommerce-api/ecommerce-api-linux-amd64 \
        /opt/ecommerce-api/ecommerce-api-linux-amd64.backup

# 3. Upload new binary
scp bin/ecommerce-api-linux-amd64 user@server:/opt/ecommerce-api/

# 4. Restart service
sudo systemctl restart ecommerce-api

# 5. Verify
curl http://localhost:8080/health

# 6. Rollback if needed
sudo cp /opt/ecommerce-api/ecommerce-api-linux-amd64.backup \
        /opt/ecommerce-api/ecommerce-api-linux-amd64
sudo systemctl restart ecommerce-api
```

## 🔐 Security Hardening

### Application Level
- [ ] Use environment variables for secrets
- [ ] Implement rate limiting
- [ ] Add request logging
- [ ] Implement IP whitelisting (if needed)
- [ ] Add request size limits
- [ ] Implement CSRF protection (if needed)
- [ ] Add API versioning

### Server Level
- [ ] Disable root SSH login
- [ ] Use SSH keys only
- [ ] Keep system updated
- [ ] Install fail2ban
- [ ] Configure firewall
- [ ] Enable SELinux/AppArmor
- [ ] Regular security audits

### Database Level
- [ ] Use strong passwords
- [ ] Limit network access
- [ ] Enable SSL connections
- [ ] Regular backups
- [ ] Audit logging
- [ ] Principle of least privilege

## 📊 Performance Optimization

### Application
- [ ] Enable Gin release mode
- [ ] Configure connection pooling
- [ ] Implement caching (Redis)
- [ ] Optimize database queries
- [ ] Add database indexes
- [ ] Use CDN for static assets

### Server
- [ ] Tune PostgreSQL settings
- [ ] Configure nginx caching
- [ ] Enable gzip compression
- [ ] Optimize kernel parameters
- [ ] Monitor resource usage

## 🆘 Troubleshooting

### Service won't start
```bash
# Check service status
sudo systemctl status ecommerce-api

# Check logs
sudo journalctl -u ecommerce-api -n 50

# Check binary permissions
ls -la /opt/ecommerce-api/

# Test binary manually
cd /opt/ecommerce-api
./ecommerce-api-linux-amd64
```

### Database connection issues
```bash
# Test database connection
psql -h localhost -U ecommerce_user -d ecommerce_db

# Check PostgreSQL status
sudo systemctl status postgresql

# Check PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-14-main.log
```

### High CPU/Memory usage
```bash
# Monitor resources
htop
# or
top

# Check API processes
ps aux | grep ecommerce-api

# Check database connections
sudo -u postgres psql -c "SELECT * FROM pg_stat_activity;"
```

## ✅ Post-Deployment Verification

- [ ] Health endpoint responds: `curl https://api.yourdomain.com/health`
- [ ] User registration works
- [ ] User login works
- [ ] Admin login works
- [ ] Permission system works
- [ ] Pagination works
- [ ] Token refresh works
- [ ] SSL certificate valid
- [ ] Logs are being written
- [ ] Monitoring is active
- [ ] Backups are configured

## 📞 Emergency Contacts

- **DevOps Team**: [contact info]
- **Database Admin**: [contact info]
- **Security Team**: [contact info]

## 📚 Additional Resources

- [Go Production Best Practices](https://golang.org/doc/)
- [PostgreSQL Performance](https://www.postgresql.org/docs/current/performance-tips.html)
- [Nginx Optimization](https://nginx.org/en/docs/)
- [Let's Encrypt](https://letsencrypt.org/)

---

**Last Updated**: November 2024
**Version**: 1.0
