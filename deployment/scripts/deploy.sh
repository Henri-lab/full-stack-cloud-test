#!/bin/bash

# Deployment script for FullStack application
# Usage: ./deploy.sh [environment]
# Example: ./deploy.sh production

set -e

ENVIRONMENT=${1:-production}
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "========================================="
echo "Deploying FullStack Application"
echo "Environment: $ENVIRONMENT"
echo "========================================="

# Load environment variables
if [ -f "$PROJECT_ROOT/.env.$ENVIRONMENT" ]; then
    echo "Loading environment variables from .env.$ENVIRONMENT"
    export $(cat "$PROJECT_ROOT/.env.$ENVIRONMENT" | grep -v '^#' | xargs)
else
    echo "Warning: .env.$ENVIRONMENT file not found"
    echo "Using .env.example as template"
fi

# Navigate to deployment directory
cd "$PROJECT_ROOT/deployment"

# Pull latest changes
echo "Pulling latest changes from git..."
cd "$PROJECT_ROOT"
git pull origin main || echo "Warning: Git pull failed or not a git repository"

# Build and deploy with Docker Compose
cd "$PROJECT_ROOT/deployment"

if [ "$ENVIRONMENT" = "production" ]; then
    echo "Building and starting production containers..."
    docker-compose -f docker-compose.prod.yml down
    docker-compose -f docker-compose.prod.yml build --no-cache
    docker-compose -f docker-compose.prod.yml up -d
else
    echo "Building and starting development containers..."
    docker-compose -f docker-compose.yml down
    docker-compose -f docker-compose.yml build --no-cache
    docker-compose -f docker-compose.yml up -d
fi

# Wait for services to be healthy
echo "Waiting for services to be healthy..."
sleep 10

# Check service status
echo "Checking service status..."
if [ "$ENVIRONMENT" = "production" ]; then
    docker-compose -f docker-compose.prod.yml ps
else
    docker-compose -f docker-compose.yml ps
fi

# Run database migrations (if needed)
echo "Running database migrations..."
# Add migration commands here if needed

echo "========================================="
echo "Deployment completed successfully!"
echo "========================================="
echo ""
echo "Services:"
echo "  - Frontend: http://localhost"
echo "  - Backend API: http://localhost:8080"
echo "  - Grafana: http://localhost:3000"
echo "  - Prometheus: http://localhost:9090"
echo ""
echo "To view logs:"
if [ "$ENVIRONMENT" = "production" ]; then
    echo "  docker-compose -f docker-compose.prod.yml logs -f"
else
    echo "  docker-compose -f docker-compose.yml logs -f"
fi
