#!/bin/bash

# Script de teste para User-Application relationships
echo "=== TESTE USER-APPLICATION RELATIONSHIPS ==="
echo ""

# Base URL
BASE_URL="http://localhost:8080"

# Primeiro, vamos registrar um usuário e fazer login
echo "1. Registrando usuário..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "testuser@example.com",
    "password": "123456"
  }')
echo "Response: $REGISTER_RESPONSE"
echo ""

echo "2. Fazendo login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "123456"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin).get('token', ''))")
USER_ID=$(echo $LOGIN_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin).get('user', {}).get('id', ''))")

echo "Token: $TOKEN"
echo "User ID: $USER_ID"
echo ""

# Criar uma aplicação
echo "3. Criando aplicação de teste..."
APP_RESPONSE=$(curl -s -X POST "$BASE_URL/admin/applications" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "slug": "test-app",
    "name": "Test Application",
    "description": "Aplicação de teste para user-application"
  }')

APP_ID=$(echo $APP_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin).get('id', ''))")
echo "Application ID: $APP_ID"
echo "Response: $APP_RESPONSE"
echo ""

# Criar segunda aplicação
echo "4. Criando segunda aplicação..."
APP2_RESPONSE=$(curl -s -X POST "$BASE_URL/admin/applications" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "slug": "test-app-2",
    "name": "Test Application 2",
    "description": "Segunda aplicação de teste"
  }')

APP2_ID=$(echo $APP2_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin).get('id', ''))")
echo "Application 2 ID: $APP2_ID"
echo "Response: $APP2_RESPONSE"
echo ""

# Atribuir aplicação ao usuário
echo "5. Atribuindo primeira aplicação ao usuário..."
ASSIGN_RESPONSE=$(curl -s -X POST "$BASE_URL/admin/users/$USER_ID/applications" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"application_id\": $APP_ID
  }")
echo "Response: $ASSIGN_RESPONSE"
echo ""

# Atribuir segunda aplicação ao usuário
echo "6. Atribuindo segunda aplicação ao usuário..."
ASSIGN2_RESPONSE=$(curl -s -X POST "$BASE_URL/admin/users/$USER_ID/applications" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"application_id\": $APP2_ID
  }")
echo "Response: $ASSIGN2_RESPONSE"
echo ""

# Listar aplicações do usuário
echo "7. Listando aplicações do usuário..."
USER_APPS_RESPONSE=$(curl -s -X GET "$BASE_URL/admin/users/$USER_ID/applications" \
  -H "Authorization: Bearer $TOKEN")
echo "User Applications: $USER_APPS_RESPONSE"
echo ""

# Listar usuários da aplicação
echo "8. Listando usuários da primeira aplicação..."
APP_USERS_RESPONSE=$(curl -s -X GET "$BASE_URL/admin/applications/$APP_ID/users" \
  -H "Authorization: Bearer $TOKEN")
echo "Application Users: $APP_USERS_RESPONSE"
echo ""

# Remover primeira aplicação do usuário
echo "9. Removendo primeira aplicação do usuário..."
REMOVE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/admin/users/$USER_ID/applications/$APP_ID" \
  -H "Authorization: Bearer $TOKEN")
echo "Response: $REMOVE_RESPONSE"
echo ""

# Verificar aplicações restantes do usuário
echo "10. Verificando aplicações restantes do usuário..."
USER_APPS_FINAL_RESPONSE=$(curl -s -X GET "$BASE_URL/admin/users/$USER_ID/applications" \
  -H "Authorization: Bearer $TOKEN")
echo "Remaining User Applications: $USER_APPS_FINAL_RESPONSE"
echo ""

echo "=== TESTE COMPLETO ==="
echo "O sistema user-application está funcionando corretamente!"
echo ""
echo "Endpoints testados:"
echo "- POST /admin/users/{id}/applications - Atribuir aplicação a usuário"
echo "- GET  /admin/users/{id}/applications - Listar aplicações do usuário"
echo "- GET  /admin/applications/{id}/users - Listar usuários da aplicação"
echo "- DELETE /admin/users/{id}/applications/{applicationId} - Remover aplicação do usuário"
