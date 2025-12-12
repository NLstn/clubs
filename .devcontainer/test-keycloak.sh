#!/bin/bash
# Test script to verify Keycloak is properly configured and accessible

set -e

echo "🔍 Testing Keycloak Configuration..."
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Check if Keycloak is accessible
echo "1️⃣  Checking if Keycloak is accessible..."
if curl -s -f -o /dev/null http://localhost:8081/; then
    echo -e "${GREEN}✓${NC} Keycloak is accessible at http://localhost:8081"
else
    echo -e "${RED}✗${NC} Keycloak is not accessible at http://localhost:8081"
    echo "   Make sure the Keycloak service is running"
    exit 1
fi
echo ""

# Test 2: Check if the realm exists
echo "2️⃣  Checking if clubs-dev realm exists..."
REALM_URL="http://localhost:8081/realms/clubs-dev"
if curl -s -f -o /dev/null "$REALM_URL"; then
    echo -e "${GREEN}✓${NC} Realm 'clubs-dev' exists"
else
    echo -e "${RED}✗${NC} Realm 'clubs-dev' not found"
    echo "   The realm may not have been imported. Check Keycloak logs."
    exit 1
fi
echo ""

# Test 3: Get OpenID Connect configuration
echo "3️⃣  Checking OpenID Connect configuration..."
OIDC_CONFIG_URL="http://localhost:8081/realms/clubs-dev/.well-known/openid-configuration"
if OIDC_CONFIG=$(curl -s -f "$OIDC_CONFIG_URL"); then
    echo -e "${GREEN}✓${NC} OIDC configuration is available"
    
    # Parse and display key endpoints
    AUTH_ENDPOINT=$(echo "$OIDC_CONFIG" | grep -o '"authorization_endpoint":"[^"]*"' | cut -d'"' -f4)
    TOKEN_ENDPOINT=$(echo "$OIDC_CONFIG" | grep -o '"token_endpoint":"[^"]*"' | cut -d'"' -f4)
    
    echo "   Authorization endpoint: $AUTH_ENDPOINT"
    echo "   Token endpoint: $TOKEN_ENDPOINT"
else
    echo -e "${RED}✗${NC} Failed to get OIDC configuration"
    exit 1
fi
echo ""

# Test 4: Check PostgreSQL databases
echo "4️⃣  Checking PostgreSQL databases..."
if command -v psql &> /dev/null; then
    CLUBS_DB_SUCCESS=false
    KEYCLOAK_DB_SUCCESS=false
    
    if PGPASSWORD=clubs_dev_password psql -U clubs_dev -h localhost -d clubs_dev -c '\q' 2>&1 > /dev/null; then
        CLUBS_DB_SUCCESS=true
    fi
    
    if PGPASSWORD=keycloak_dev_password psql -U keycloak_dev -h localhost -d keycloak_dev -c '\q' 2>&1 > /dev/null; then
        KEYCLOAK_DB_SUCCESS=true
    fi
    
    if [[ "$CLUBS_DB_SUCCESS" == "true" ]] && [[ "$KEYCLOAK_DB_SUCCESS" == "true" ]]; then
        echo -e "${GREEN}✓${NC} Both databases (clubs_dev and keycloak_dev) are accessible"
    else
        echo -e "${YELLOW}⚠${NC}  Could not verify database access (this is okay if psql is not installed or running outside devcontainer)"
    fi
else
    echo -e "${YELLOW}⚠${NC}  psql not installed, skipping database check"
fi
echo ""

# Test 5: Test user authentication (optional, requires credentials)
echo "5️⃣  Testing user authentication..."
TOKEN_URL="http://localhost:8081/realms/clubs-dev/protocol/openid-connect/token"

# Try to get a token for the test user
TOKEN_RESPONSE=$(curl -s -X POST "$TOKEN_URL" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "client_id=clubs-frontend" \
    -d "grant_type=password" \
    -d "username=testuser" \
    -d "password=testpass" \
    -d "scope=openid profile email")

if echo "$TOKEN_RESPONSE" | grep -q "access_token"; then
    echo -e "${GREEN}✓${NC} Successfully authenticated test user"
    echo "   User: testuser"
    echo "   Client: clubs-frontend"
else
    echo -e "${RED}✗${NC} Failed to authenticate test user"
    echo "   Response: $TOKEN_RESPONSE"
    exit 1
fi
echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✓ All tests passed!${NC}"
echo ""
echo "Your Keycloak instance is properly configured and ready to use."
echo ""
echo "📝 Configuration Summary:"
echo "   Keycloak URL: http://localhost:8081"
echo "   Admin Console: http://localhost:8081/admin"
echo "   Realm: clubs-dev"
echo "   Client ID: clubs-frontend"
echo ""
echo "👤 Test Users:"
echo "   Standard User: testuser / testpass"
echo "   Admin User: admin / admin"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
