# Quick Azure Container Apps Deployment Commands

## Prerequisites
```bash
# Login to Azure
az login

# Set subscription (if needed)
az account set --subscription "your-subscription-id"

# Install Container Apps extension
az extension add --name containerapp --upgrade
```

## Build and Push Image
```bash
# Navigate to Backend directory
cd Backend

# Login to Azure Container Registry
az acr login --name your-registry-name

# Build and push image
docker build -t your-registry.azurecr.io/clubs-backend:latest .
docker push your-registry.azurecr.io/clubs-backend:latest
```

## Create Container App (First Time)
```bash
# Create Container Apps environment (if not exists)
az containerapp env create \
  --name clubs-env \
  --resource-group your-resource-group \
  --location "East US"

# Create the container app
az containerapp create \
  --name clubs-backend \
  --resource-group your-resource-group \
  --environment clubs-env \
  --image your-registry.azurecr.io/clubs-backend:latest \
  --target-port 8080 \
  --ingress external \
  --min-replicas 1 \
  --max-replicas 5 \
  --cpu 0.5 \
  --memory 1Gi
```

## Update Existing Container App
```bash
# Update with new image
az containerapp update \
  --name clubs-backend \
  --resource-group your-resource-group \
  --image your-registry.azurecr.io/clubs-backend:latest
```

## Configure Environment Variables
```bash
# Set environment variables
az containerapp update \
  --name clubs-backend \
  --resource-group your-resource-group \
  --set-env-vars \
    DATABASE_URL=your-postgres-server.postgres.database.azure.com \
    DATABASE_PORT=5432 \
    DATABASE_USER=clubs \
    FRONTEND_URL=https://your-frontend-url.com
```

## Configure Secrets
```bash
# Add secrets for sensitive data
az containerapp secret set \
  --name clubs-backend \
  --resource-group your-resource-group \
  --secrets \
    db-password=your-secure-password \
    azure-tenant-id=your-tenant-id \
    azure-client-id=your-client-id \
    azure-client-secret=your-client-secret

# Update app to use secrets
az containerapp update \
  --name clubs-backend \
  --resource-group your-resource-group \
  --set-env-vars \
    DATABASE_USER_PASSWORD=secretref:db-password \
    AZURE_TENANT_ID=secretref:azure-tenant-id \
    AZURE_CLIENT_ID=secretref:azure-client-id \
    AZURE_CLIENT_SECRET=secretref:azure-client-secret
```

## Configure Health Probes
```bash
az containerapp update \
  --name clubs-backend \
  --resource-group your-resource-group \
  --readiness-probe-path "/health" \
  --readiness-probe-initial-delay 5 \
  --readiness-probe-period 10 \
  --liveness-probe-path "/health" \
  --liveness-probe-initial-delay 30 \
  --liveness-probe-period 30
```

## View Logs
```bash
# View container logs
az containerapp logs show \
  --name clubs-backend \
  --resource-group your-resource-group \
  --follow

# View specific revision logs
az containerapp revision list \
  --name clubs-backend \
  --resource-group your-resource-group

az containerapp logs show \
  --name clubs-backend \
  --resource-group your-resource-group \
  --revision your-revision-name
```

## Monitor Application
```bash
# Get app URL
az containerapp show \
  --name clubs-backend \
  --resource-group your-resource-group \
  --query "properties.configuration.ingress.fqdn" \
  --output tsv

# Test health endpoint
curl https://your-app-url.azurecontainerapps.io/health

# Check app status
az containerapp show \
  --name clubs-backend \
  --resource-group your-resource-group \
  --query "properties.runningStatus"
```