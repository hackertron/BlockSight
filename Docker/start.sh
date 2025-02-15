#!/bin/bash

# Change to the Docker directory
cd "$(dirname "$0")"

# Create required directories if they don't exist
mkdir -p prometheus/data
mkdir -p grafana/provisioning/datasources
mkdir -p grafana/provisioning/dashboards
mkdir -p grafana/dashboards

# Copy configuration files if they don't exist
if [ ! -f "prometheus/prometheus.yml" ]; then
    cp prometheus.yml.example prometheus/prometheus.yml
fi

if [ ! -f "grafana/provisioning/datasources/prometheus.yaml" ]; then
    cp grafana/provisioning/datasources/prometheus.yaml.example grafana/provisioning/datasources/prometheus.yaml
fi

# Copy dashboard provider configuration if it doesn't exist
if [ ! -f "grafana/provisioning/dashboards/default.yaml" ]; then
    cp grafana/provisioning/dashboards/default.yaml.example grafana/provisioning/dashboards/default.yaml
fi

# Copy dashboard configuration if it doesn't exist
if [ ! -f "grafana/dashboards/blockchain_metrics.json" ]; then
    cp grafana/dashboards/blockchain_metrics.json.example grafana/dashboards/blockchain_metrics.json
fi

# Check if .env file exists
if [ ! -f ../.env ]; then
    echo "Error: .env file not found in project root"
    echo "Please create .env file with required environment variables"
    echo "Example:"
    echo "ALCHEMY_API_KEY=your_api_key"
    echo "ALCHEMY_NETWORK=eth-mainnet"
    exit 1
fi

# Stop any running containers
echo "Stopping existing containers..."
docker-compose down -v

# Build and start containers
echo "Building and starting containers..."
docker-compose up --build -d

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 30

# Check service status
echo "Checking service status..."
docker-compose ps

# Show logs of the blocksight service
echo "Showing blocksight logs..."
docker-compose logs -f blocksight