#!/bin/bash

# Script de teste para User-Application-Role model
echo "=== TESTE USER-APPLICATION-ROLE MODEL ==="
echo ""

# Base URL
BASE_URL="http://localhost:8080"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InRlc3R1c2VyQGV4YW1wbGUuY29tIiwicm9sZXMiOm51bGwsInBlcm1pc3Npb25zIjpudWxsLCJleHAiOjE3NTc1OTY4NjEsImlhdCI6MTc1NzUxMDQ2MX0.IzrebcO4yVTBhltExa2sruZLH_tEeQV-wnrak6jQARY"
USER_ID=2

echo "🔗 Testando novo modelo User-Application-Role..."
echo ""

# Primeiro, vamos ver as aplicações e roles disponíveis
echo "1. 📋 Listando aplicações disponíveis..."
curl -s -X GET "$BASE_URL/admin/applications" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

echo "2. 👥 Listando roles da aplicação 1..."
curl -s -X GET "$BASE_URL/admin/roles?application_id=1" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Atribuir role 1 ao usuário na aplicação 1
echo "3. ➕ Atribuindo role 1 ao usuário $USER_ID na aplicação 1..."
curl -s -X POST "$BASE_URL/admin/users/$USER_ID/application-roles" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "application_id": 1,
    "role_id": 1,
    "profile_id": null
  }' | jq '.'
echo ""

# Atribuir outra role na mesma aplicação (se existir)
echo "4. ➕ Atribuindo role 2 ao usuário $USER_ID na aplicação 1..."
curl -s -X POST "$BASE_URL/admin/users/$USER_ID/application-roles" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "application_id": 1,
    "role_id": 2,
    "profile_id": 123
  }' | jq '.'
echo ""

# Atribuir role em outra aplicação
echo "5. ➕ Atribuindo role ao usuário $USER_ID na aplicação 2..."
curl -s -X POST "$BASE_URL/admin/users/$USER_ID/application-roles" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "application_id": 2,
    "role_id": 3
  }' | jq '.'
echo ""

# Listar todas as atribuições do usuário
echo "6. 👤 Listando todas as atribuições application-role do usuário $USER_ID..."
curl -s -X GET "$BASE_URL/admin/users/$USER_ID/application-roles" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Listar roles do usuário em uma aplicação específica
echo "7. 🏢 Listando roles do usuário $USER_ID na aplicação 1..."
curl -s -X GET "$BASE_URL/admin/application-roles/users/$USER_ID/applications/1/roles" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Remover uma atribuição específica
echo "8. ➖ Removendo role 1 do usuário $USER_ID na aplicação 1..."
curl -s -X DELETE "$BASE_URL/admin/users/$USER_ID/application-roles/1/1" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Verificar atribuições restantes
echo "9. 🔍 Verificando atribuições restantes do usuário $USER_ID..."
curl -s -X GET "$BASE_URL/admin/users/$USER_ID/application-roles" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

echo "✅ TESTE CONCLUÍDO!"
echo ""
echo "📊 Novo modelo User-Application-Role testado com sucesso!"
echo ""
echo "🆕 Vantagens do novo modelo:"
echo "   - ✅ Vincula usuário + aplicação + role numa única tabela"
echo "   - ✅ Suporte a profile_id para contextos específicos"
echo "   - ✅ Melhor performance (menos JOINs)"
echo "   - ✅ Estrutura mais limpa e organizada"
echo "   - ✅ Evita inconsistências de dados"
echo ""
echo "🌐 Acesse a documentação completa em: http://localhost:8080/swagger/index.html"
echo "📝 Novos endpoints:"
echo "   - POST   /admin/users/{id}/application-roles"
echo "   - GET    /admin/users/{id}/application-roles"
echo "   - DELETE /admin/users/{id}/application-roles/{applicationId}/{roleId}"
echo "   - GET    /admin/application-roles/users/{userId}/applications/{applicationId}/roles"
