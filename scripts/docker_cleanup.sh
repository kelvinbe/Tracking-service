#!/bin/bash

# Get the list of running container IDs
container_ids=$(docker ps -q)

# Get list of any stopped containers
stopped_containers=$(docker ps -a -q)

# Check if there are any stopped containers
if [ -n "$stopped_containers" ]; then
    echo "Removing stopped containers..."
    docker rm $stopped_containers
else
    echo "No stopped containers found."
fi

# Check if there are any running containers
if [ -n "$container_ids" ]; then
    echo "Stopping running containers..."
    docker stop $container_ids
    docker rm $container_ids
else
    echo "No running containers found."
fi

# Check for any images
image_ids=$(docker images -q)

# Check if there are any images
if [ -n "$image_ids" ]; then
    echo "Removing images..."
    docker rmi $image_ids
else
    echo "No images found."
fi
