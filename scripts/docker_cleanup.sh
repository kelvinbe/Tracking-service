#!/bin/bash

# Get the list of running container IDs from image 'divvly-tracking-service'
container_ids=$(sudo docker ps -a -q -f ancestor=divvly-tracking-service)

# Check if there are any containers from 'divvly-tracking-service'
if [ -n "$container_ids" ]; then
    echo "Stopping and removing containers from image 'divvly-tracking-service'..."
    sudo docker stop $container_ids
    sudo docker rm $container_ids
else
    echo "No containers from image 'divvly-tracking-service' found."
fi

# Get the image ID of 'divvly-tracking-service'
image_id=$(sudo docker images -q divvly-tracking-service)

# Check if the image 'divvly-tracking-service' exists
if [ -n "$image_id" ]; then
    echo "Removing image 'divvly-tracking-service'..."
    sudo docker rmi $image_id
else
    echo "Image 'divvly-tracking-service' not found."
fi
