#!/bin/bash

function start-postgres() {
    echo "Starting postgres"
    podman run --name ns-test-postgres -d --rm --replace \
            -e POSTGRES_USER=quarkus -e POSTGRES_PASSWORD=quarkus -e POSTGRES_DB=nameserver \
            -p 5432:5432 docker.io/library/postgres:17
}

function start-keycloak() {
    if ! podman image exists localhost/ns-test-keycloak:0; then
        echo "Creating keycloak test image"
        podman build -t ns-test-keycloak:0 -f Dockerfile.keycloak .
    fi

    echo "Starting keycloak"
    podman run --name ns-test-keycloak -d --rm --replace -e KEYCLOAK_ADMIN=admin \
        -e KEYCLOAK_ADMIN_PASSWORD=admin -p 8180:8080 localhost/ns-test-keycloak:0    
}

start-postgres
start-keycloak
echo "All services started"