#!/bin/bash

# Build and Push Script for RHACM Global Hub Monitor
# Images will be pushed to quay.io/bzhai/

set -e

# Configuration
REGISTRY="quay.io/bzhai"
BACKEND_IMAGE="${REGISTRY}/rhacm-monitor-backend"
FRONTEND_IMAGE="${REGISTRY}/rhacm-monitor-frontend"
VERSION="${VERSION:-latest}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}RHACM Global Hub Monitor Build Script${NC}"
echo -e "${GREEN}=====================================${NC}"

# Check if podman or docker is available
if command -v podman &> /dev/null; then
    CONTAINER_CLI="podman"
elif command -v docker &> /dev/null; then
    CONTAINER_CLI="docker"
else
    echo -e "${RED}Error: Neither podman nor docker found${NC}"
    exit 1
fi

echo -e "${YELLOW}Using container CLI: ${CONTAINER_CLI}${NC}"

# Login to registry (optional, comment out if already logged in)
# echo -e "${YELLOW}Logging in to ${REGISTRY}...${NC}"
# ${CONTAINER_CLI} login ${REGISTRY}

# Build Backend Image
echo -e "${YELLOW}Building backend image...${NC}"
${CONTAINER_CLI} build \
    -f deployment/docker/Dockerfile.backend \
    -t ${BACKEND_IMAGE}:${VERSION} \
    -t ${BACKEND_IMAGE}:latest \
    .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Backend image built successfully${NC}"
else
    echo -e "${RED}✗ Backend build failed${NC}"
    exit 1
fi

# Build Frontend Image
echo -e "${YELLOW}Building frontend image...${NC}"
${CONTAINER_CLI} build \
    -f deployment/docker/Dockerfile.frontend \
    -t ${FRONTEND_IMAGE}:${VERSION} \
    -t ${FRONTEND_IMAGE}:latest \
    .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Frontend image built successfully${NC}"
else
    echo -e "${RED}✗ Frontend build failed${NC}"
    exit 1
fi

# Push images
echo -e "${YELLOW}Pushing images to ${REGISTRY}...${NC}"

echo -e "${YELLOW}Pushing backend image...${NC}"
${CONTAINER_CLI} push ${BACKEND_IMAGE}:${VERSION}
${CONTAINER_CLI} push ${BACKEND_IMAGE}:latest

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Backend image pushed successfully${NC}"
else
    echo -e "${RED}✗ Backend push failed${NC}"
    exit 1
fi

echo -e "${YELLOW}Pushing frontend image...${NC}"
${CONTAINER_CLI} push ${FRONTEND_IMAGE}:${VERSION}
${CONTAINER_CLI} push ${FRONTEND_IMAGE}:latest

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Frontend image pushed successfully${NC}"
else
    echo -e "${RED}✗ Frontend push failed${NC}"
    exit 1
fi

echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}Build and Push Completed Successfully!${NC}"
echo -e "${GREEN}=====================================${NC}"
echo ""
echo "Backend Image:  ${BACKEND_IMAGE}:${VERSION}"
echo "Frontend Image: ${FRONTEND_IMAGE}:${VERSION}"
echo ""
echo "To deploy, update the image references in:"
echo "  - deployment/k8s/backend-deployment.yaml"
echo "  - deployment/k8s/frontend-deployment.yaml"

