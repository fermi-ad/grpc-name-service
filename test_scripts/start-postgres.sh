#!/bin/bash

echo "Starting postgres"
podman run -d --name ns-test-postgres --rm --replace \
  -e POSTGRES_USER=quarkus \
  -e POSTGRES_PASSWORD=quarkus \
  -e POSTGRES_DB=nameserver \
  -p 5432:5432 \
  docker.io/library/postgres:17
