#!/bin/bash

set -e

# Step 1: Build the Go binary for linux/amd64
echo "Building file-uploader-orch binary for linux/amd64"
GOARCH=amd64 GOOS=linux go build -o file-uploader-orch

# Step 2: Login to Docker Hub
echo "Login and pushing to Docker Hub"
docker login

# Step 3: Set up docker buildx for multi-platform support (if not already done)
docker buildx create --use  # Only need to do this once per session

# Step 4: Build and push the Docker image for linux/amd64
echo "Building and pushing image to Docker Hub"
docker buildx build --platform linux/amd64 -t kivuos1999/file-uploader-orch --push .

# Step 5: Cleanup the binary after build
echo "Cleanup"
rm file-uploader-orch

# Step 6: Success message
echo "Build succeeded and image pushed to Docker Hub"
