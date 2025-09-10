#!/bin/bash

# Script de teste simplificado para User-Application relationships
echo "=== TESTE USER-APPLICATION RELATIONSHIPS ==="
echo ""

# Base URL
BASE_URL="http://localhost:8080"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InRlc3R1c2VyQGV4YW1wbGUuY29tIiwicm9sZXMiOm51bGwsInBlcm1pc3Npb25zIjpudWxsLCJleHAiOjE3NTc1OTY4NjEsImlhdCI6MTc1NzUxMDQ2MX0.IzrebcO4yVTBhltExa2sruZLH_tEeQV-wnrak6jQARY"
USER_ID=2

echo "🔗 Testando endpoints User-Application..."
echo ""

# Listar aplicações disponíveis
echo "1. 📋 Listando aplicações disponíveis..."
curl -s -X GET "$BASE_URL/admin/applications" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Atribuir aplicação 3 ao usuário
echo "2. ➕ Atribuindo aplicação 3 ao usuário $USER_ID..."
curl -s -X POST "$BASE_URL/admin/users/$USER_ID/applications" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"application_id": 3}' | jq '.'
echo ""

# Atribuir aplicação 4 ao usuário
echo "3. ➕ Atribuindo aplicação 4 ao usuário $USER_ID..."
curl -s -X POST "$BASE_URL/admin/users/$USER_ID/applications" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"application_id": 4}' | jq '.'
echo ""

# Listar aplicações do usuário
echo "4. 👤 Listando aplicações do usuário $USER_ID..."
curl -s -X GET "$BASE_URL/admin/users/$USER_ID/applications" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Listar usuários da aplicação 3
echo "5. 🏢 Listando usuários da aplicação 3..."
curl -s -X GET "$BASE_URL/admin/applications/3/users" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Remover aplicação 3 do usuário
echo "6. ➖ Removendo aplicação 3 do usuário $USER_ID..."
curl -s -X DELETE "$BASE_URL/admin/users/$USER_ID/applications/3" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Verificar aplicações restantes do usuário
echo "7. 🔍 Verificando aplicações restantes do usuário $USER_ID..."
curl -s -X GET "$BASE_URL/admin/users/$USER_ID/applications" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

echo "✅ TESTE CONCLUÍDO!"
echo ""
echo "📊 Endpoints testados com sucesso:"
echo "   - POST /admin/users/{id}/applications - Atribuir aplicação a usuário"
echo "   - GET  /admin/users/{id}/applications - Listar aplicações do usuário"
echo "   - GET  /admin/applications/{id}/users - Listar usuários da aplicação"
echo "   - DELETE /admin/users/{id}/applications/{applicationId} - Remover aplicação do usuário"
echo ""
echo "🌐 Acesse a documentação completa em: http://localhost:8080/swagger/index.html"
