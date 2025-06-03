# Use the official Microsoft devcontainers Go image as base
FROM mcr.microsoft.com/devcontainers/go:1.24

# Install Node.js (LTS version)
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g npm@latest

# Install additional tools for development
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        postgresql-client \
        build-essential \
        git \
        curl \
        wget \
        unzip \
        ca-certificates \
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Install Go tools for development
RUN go install github.com/air-verse/air@latest \
    && go install golang.org/x/tools/cmd/goimports@latest \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Set the default shell to bash for better compatibility
SHELL ["/bin/bash", "-c"]

# Create workspace directory
RUN mkdir -p /workspace
WORKDIR /workspace