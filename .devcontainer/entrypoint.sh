#!/bin/bash
set -e

# Install socat if not already installed
if ! command -v socat &> /dev/null; then
    apt-get update && apt-get install -y socat
fi

# Start socat in background to forward localhost:8081 to keycloak:8080
socat TCP-LISTEN:8081,fork,reuseaddr TCP:keycloak:8080 &

# Execute the original command or keep container running
exec "$@"
