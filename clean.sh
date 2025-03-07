#!/bin/bash

# Stop and remove all running containers
echo "Stopping all running containers..."
docker stop $(docker ps -q)

echo "Removing all containers (running and stopped)..."
docker rm $(docker ps -aq)

# Remove all images
echo "Removing all images..."
docker rmi $(docker images -q)

# Remove all unused volumes
echo "Removing all volumes..."
docker volume prune -f

# Remove unused networks
echo "Removing all unused networks..."
docker network prune -f

echo "Docker cleanup complete."