#!/bin/bash

# Server initialization script for Debian/Ubuntu
# This script sets up a fresh server for the FullStack application
# Run as root: sudo bash init-server.sh

set -e

echo "========================================="
echo "Initializing Server for FullStack App"
echo "========================================="

# Update system
echo "Updating system packages..."
apt update && apt upgrade -y

# Install essential packages
echo "Installing essential packages..."
apt install -y \
    curl \
    wget \
    git \
    ufw \
    fail2ban \
    htop \
    vim \
    ca-certificates \
    gnupg \
    lsb-release

# Install Docker
echo "Installing Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh

    # Add current user to docker group
    usermod -aG docker $USER || true

    echo "Docker installed successfully"
else
    echo "Docker is already installed"
fi

# Install Docker Compose
echo "Installing Docker Compose..."
if ! command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep 'tag_name' | cut -d\" -f4)
    curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    echo "Docker Compose installed successfully"
else
    echo "Docker Compose is already installed"
fi

# Configure firewall
echo "Configuring firewall..."
ufw --force enable
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw reload

echo "Firewall configured"

# Configure fail2ban
echo "Configuring fail2ban..."
systemctl enable fail2ban
systemctl start fail2ban

# Create application directory
echo "Creating application directory..."
mkdir -p /opt/fullstack
chown -R $USER:$USER /opt/fullstack

# Install Certbot for SSL (Let's Encrypt)
echo "Installing Certbot..."
apt install -y certbot python3-certbot-nginx

# Set up automatic security updates
echo "Setting up automatic security updates..."
apt install -y unattended-upgrades
dpkg-reconfigure -plow unattended-upgrades

# Optimize system settings
echo "Optimizing system settings..."
cat >> /etc/sysctl.conf <<EOF

# FullStack App Optimizations
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 2048
net.ipv4.ip_local_port_range = 10000 65535
vm.swappiness = 10
EOF

sysctl -p

# Create swap file if not exists
if [ ! -f /swapfile ]; then
    echo "Creating swap file..."
    fallocate -l 2G /swapfile
    chmod 600 /swapfile
    mkswap /swapfile
    swapon /swapfile
    echo '/swapfile none swap sw 0 0' >> /etc/fstab
fi

# Set up log rotation
echo "Setting up log rotation..."
cat > /etc/logrotate.d/fullstack <<EOF
/opt/fullstack/logs/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 $USER $USER
    sharedscripts
}
EOF

echo "========================================="
echo "Server initialization completed!"
echo "========================================="
echo ""
echo "Next steps:"
echo "1. Clone your repository to /opt/fullstack"
echo "2. Configure environment variables"
echo "3. Run the deployment script"
echo ""
echo "Important:"
echo "- Reboot the server to apply all changes"
echo "- Configure SSL certificates with: certbot --nginx -d yourdomain.com"
echo "- Review and customize fail2ban settings in /etc/fail2ban/"
