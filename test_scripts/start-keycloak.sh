#!/bin/bash

if podman image exists localhost/ns-test-keycloak:0; then
    echo "Keycloak test image exists"
else
    echo "Creating keycloak test image"
    podman build -t ns-test-keycloak:0 -f Dockerfile.keycloak .
fi

echo "Starting keycloak"
podman run --name ns-test-keycloak -d --rm --replace -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin -p 8180:8080 localhost/ns-test-keycloak:0
