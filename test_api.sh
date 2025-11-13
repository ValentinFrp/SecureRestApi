#!/bin/bash

set -e

BASE_URL="http://localhost:8080"
EMAIL="test_$(date +%s)@example.com"
PASSWORD="password123"
TOKEN=""

echo "🧪 Test de l'API Secure REST API"
echo "================================"
echo ""

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/health")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ Health check OK${NC}"
    echo "  Response: $body"
else
    echo -e "${RED}✗ Health check FAILED (HTTP $http_code)${NC}"
    exit 1
fi
echo ""

# Test 2: Inscription
echo -e "${YELLOW}Test 2: Inscription (POST /api/auth/register)${NC}"
echo "  Email: $EMAIL"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "201" ]; then
    echo -e "${GREEN}✓ Inscription réussie${NC}"
    TOKEN=$(echo "$body" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    echo "  Token reçu: ${TOKEN:0:50}..."
else
    echo -e "${RED}✗ Inscription FAILED (HTTP $http_code)${NC}"
    echo "  Response: $body"
    exit 1
fi
echo ""

# Test 3: Inscription avec email existant (doit échouer)
echo -e "${YELLOW}Test 3: Inscription avec email existant (doit échouer)${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "409" ]; then
    echo -e "${GREEN}✓ Rejet correct de l'email dupliqué${NC}"
    echo "  Response: $body"
else
    echo -e "${RED}✗ Test FAILED (attendu HTTP 409, reçu $http_code)${NC}"
    exit 1
fi
echo ""

# Test 4: Connexion avec mauvais mot de passe (doit échouer)
echo -e "${YELLOW}Test 4: Connexion avec mauvais mot de passe (doit échouer)${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"wrongpassword\"}")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "401" ]; then
    echo -e "${GREEN}✓ Rejet correct du mauvais mot de passe${NC}"
    echo "  Response: $body"
else
    echo -e "${RED}✗ Test FAILED (attendu HTTP 401, reçu $http_code)${NC}"
    exit 1
fi
echo ""

# Test 5: Connexion réussie
echo -e "${YELLOW}Test 5: Connexion (POST /api/auth/login)${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ Connexion réussie${NC}"
    TOKEN=$(echo "$body" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    echo "  Nouveau token: ${TOKEN:0:50}..."
else
    echo -e "${RED}✗ Connexion FAILED (HTTP $http_code)${NC}"
    echo "  Response: $body"
    exit 1
fi
echo ""

# Test 6: Accès route protégée sans token (doit échouer)
echo -e "${YELLOW}Test 6: Accès route protégée sans token (doit échouer)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/auth/me")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "401" ]; then
    echo -e "${GREEN}✓ Rejet correct sans token${NC}"
    echo "  Response: $body"
else
    echo -e "${RED}✗ Test FAILED (attendu HTTP 401, reçu $http_code)${NC}"
    exit 1
fi
echo ""

# Test 7: Accès route protégée avec token invalide (doit échouer)
echo -e "${YELLOW}Test 7: Accès route protégée avec token invalide (doit échouer)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/auth/me" \
    -H "Authorization: Bearer invalidtoken123")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "401" ]; then
    echo -e "${GREEN}✓ Rejet correct du token invalide${NC}"
    echo "  Response: $body"
else
    echo -e "${RED}✗ Test FAILED (attendu HTTP 401, reçu $http_code)${NC}"
    exit 1
fi
echo ""

# Test 8: Accès route protégée avec token valide
echo -e "${YELLOW}Test 8: Profil utilisateur (GET /api/auth/me) - Route Protégée${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/auth/me" \
    -H "Authorization: Bearer $TOKEN")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ Accès route protégée réussi${NC}"
    echo "  Profil: $body"
else
    echo -e "${RED}✗ Accès route protégée FAILED (HTTP $http_code)${NC}"
    echo "  Response: $body"
    exit 1
fi
echo ""

# Résumé
echo "================================"
echo -e "${GREEN}✅ Tous les tests ont réussi !${NC}"
echo ""
echo "📊 Résumé:"
echo "  - Health check: ✓"
echo "  - Inscription: ✓"
echo "  - Email dupliqué rejeté: ✓"
echo "  - Mauvais mot de passe rejeté: ✓"
echo "  - Connexion: ✓"
echo "  - Sans token rejeté: ✓"
echo "  - Token invalide rejeté: ✓"
echo "  - Route protégée avec token: ✓"
echo ""
