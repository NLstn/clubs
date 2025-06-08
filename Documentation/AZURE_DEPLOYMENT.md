# Azure Container Apps Deployment Guide

This guide explains how to deploy the Clubs backend to Azure Container Apps for staging and production environments.

## Prerequisites

- Azure CLI installed and configured
- Docker installed
- Azure Container Registry (or other container registry)
- Azure Container Apps Environment

## Container Configuration

The backend is containerized using Docker with the following features:

### Docker Image Features
- **Multi-stage build**: Optimized for production with minimal runtime image
- **Non-root user**: Runs as `appuser` (UID 1001) for security
- **Health checks**: Built-in health endpoint at `/health`
- **Alpine Linux**: Minimal attack surface and smaller image size
- **Static binary**: No external dependencies in runtime

### Health Check Endpoint
The application provides a health check endpoint at `/health` that:
- Returns HTTP 200 for healthy status
- Returns HTTP 503 for unhealthy status
- Checks database connectivity
- Provides detailed service status in JSON format

Example healthy response:
```json
{
  "status": "healthy",
  "services": {
    "api": "healthy",
    "database": "healthy"
  }
}
```

## Environment Variables

### Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL host | `your-postgres-server.postgres.database.azure.com` |
| `DATABASE_PORT` | PostgreSQL port | `5432` |
| `DATABASE_USER` | Database username | `clubs` |
| `DATABASE_USER_PASSWORD` | Database password | `your-secure-password` |
| `AZURE_TENANT_ID` | Azure AD Tenant ID | `your-tenant-id` |
| `AZURE_CLIENT_ID` | Azure AD Client ID | `your-client-id` |
| `AZURE_CLIENT_SECRET` | Azure AD Client Secret | `your-client-secret` |

### Optional Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FRONTEND_URL` | Frontend URL for CORS | `http://localhost:5173` |
| `AZURE_ACS_ENDPOINT` | Azure Communication Services endpoint | |
| `AZURE_ACS_SENDER_ADDRESS` | Email sender address | |

## Deployment Steps

### 1. Build and Push Container Image

```bash
# Navigate to the Backend directory
cd Backend

# Build the Docker image
docker build -t your-registry.azurecr.io/clubs-backend:latest .

# Push to Azure Container Registry
az acr login --name your-registry
docker push your-registry.azurecr.io/clubs-backend:latest
```

### 2. Create Container App

```bash
# Create a Container App
az containerapp create \
  --name clubs-backend \
  --resource-group your-resource-group \
  --environment your-container-env \
  --image your-registry.azurecr.io/clubs-backend:latest \
  --target-port 8080 \
  --ingress external \
  --min-replicas 1 \
  --max-replicas 10 \
  --cpu 0.5 \
  --memory 1Gi \
  --env-vars \
    DATABASE_URL=your-postgres-server.postgres.database.azure.com \
    DATABASE_PORT=5432 \
    DATABASE_USER=clubs \
    DATABASE_USER_PASSWORD=secretref:db-password \
    AZURE_TENANT_ID=secretref:azure-tenant-id \
    AZURE_CLIENT_ID=secretref:azure-client-id \
    AZURE_CLIENT_SECRET=secretref:azure-client-secret \
    FRONTEND_URL=https://your-frontend-url.com
```

### 3. Configure Secrets

```bash
# Add secrets (recommended for sensitive data)
az containerapp secret set \
  --name clubs-backend \
  --resource-group your-resource-group \
  --secrets \
    db-password=your-secure-password \
    azure-tenant-id=your-tenant-id \
    azure-client-id=your-client-id \
    azure-client-secret=your-client-secret
```

### 4. Configure Health Probes

```bash
# Update the container app with health probes
az containerapp update \
  --name clubs-backend \
  --resource-group your-resource-group \
  --set-env-vars HEALTH_CHECK_ENABLED=true \
  --startup-probe-path /health \
  --startup-probe-initial-delay 10 \
  --startup-probe-period 30 \
  --startup-probe-timeout 5 \
  --liveness-probe-path /health \
  --liveness-probe-initial-delay 30 \
  --liveness-probe-period 30 \
  --liveness-probe-timeout 5 \
  --readiness-probe-path /health \
  --readiness-probe-initial-delay 5 \
  --readiness-probe-period 10 \
  --readiness-probe-timeout 3
```

## Azure Container Apps Configuration

### Recommended Resource Limits

```yaml
resources:
  cpu: "0.5"
  memory: "1Gi"
```

### Scaling Configuration

```yaml
scale:
  minReplicas: 1
  maxReplicas: 10
  rules:
  - name: "http-rule"
    http:
      metadata:
        concurrentRequests: "10"
```

### Ingress Configuration

```yaml
ingress:
  external: true
  targetPort: 8080
  allowInsecure: false
  traffic:
  - weight: 100
    latestRevision: true
```

## Database Setup

Ensure your PostgreSQL database is configured with:

1. **PostgreSQL Extensions**: The application requires the `pgcrypto` extension
2. **Network Access**: Allow connections from Azure Container Apps
3. **SSL/TLS**: Recommended for production environments
4. **Connection Limits**: Configure appropriate connection limits

Example database initialization:
```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE USER clubs WITH PASSWORD 'your-secure-password';
CREATE DATABASE clubs OWNER clubs;
```

## Monitoring and Logging

### Application Insights Integration

The application logs requests and responses. To integrate with Application Insights:

1. Add Application Insights instrumentation key as environment variable
2. Configure log forwarding in Container Apps
3. Set up alerts for health check failures

### Health Monitoring

The `/health` endpoint can be used for:
- Container orchestration health checks
- Load balancer health probes
- Monitoring system integration
- Automated alerting

### Log Analytics

Container Apps automatically forwards logs to Log Analytics. Monitor:
- Application logs
- Container health status
- Resource utilization
- Request/response metrics

## Security Considerations

1. **Use managed identity** for Azure services when possible
2. **Store secrets** in Azure Key Vault or Container Apps secrets
3. **Enable HTTPS only** for external ingress
4. **Configure CORS** appropriately for your frontend domain
5. **Use non-root container** (already configured)
6. **Regularly update** base images for security patches

## Troubleshooting

### Common Issues

1. **Health check failures**:
   - Check database connectivity
   - Verify environment variables
   - Review application logs

2. **Connection timeouts**:
   - Check firewall rules
   - Verify network security groups
   - Ensure database allows connections

3. **Authentication errors**:
   - Verify Azure AD configuration
   - Check client secrets expiration
   - Validate tenant ID

### Debugging Commands

```bash
# View container logs
az containerapp logs show \
  --name clubs-backend \
  --resource-group your-resource-group

# Check container status
az containerapp show \
  --name clubs-backend \
  --resource-group your-resource-group

# Test health endpoint
curl https://your-app-url.azurecontainerapps.io/health
```

## CI/CD Integration

For automated deployments, consider:

1. **GitHub Actions** or **Azure DevOps** for CI/CD
2. **Automatic image building** on code changes
3. **Progressive deployment** strategies
4. **Automated testing** before deployment
5. **Rollback mechanisms** for failed deployments

Example GitHub Actions workflow structure:
```yaml
name: Deploy to Azure Container Apps
on:
  push:
    branches: [main]
    paths: ['Backend/**']

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Build and push
      # Build and push container image
    - name: Deploy to Container Apps
      # Deploy using az containerapp up or update
```