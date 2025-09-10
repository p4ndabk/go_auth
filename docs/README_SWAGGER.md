# API Documentation - Go Auth

## Swagger/OpenAPI Documentation

Este projeto agora possui documentação completa da API usando Swagger/OpenAPI.

### Acessando a Documentação

Com o servidor rodando, acesse:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON Schema**: http://localhost:8080/swagger/doc.json
- **YAML Schema**: http://localhost:8080/swagger/swagger.yaml

### Principais Endpoints Documentados

#### Autenticação Pública
- `POST /register` - Registrar novo usuário
- `POST /login` - Login do usuário
- `GET /me` - Informações do usuário autenticado (requer token)

#### Admin - Applications (requer autenticação)
- `POST /admin/applications` - Criar aplicação
- `GET /admin/applications` - Listar aplicações
- `GET /admin/applications/{id}` - Obter aplicação específica
- `PUT /admin/applications/{id}` - Atualizar aplicação
- `DELETE /admin/applications/{id}` - Deletar aplicação

#### Admin - Roles (requer autenticação)
- `POST /admin/roles` - Criar role
- `GET /admin/roles` - Listar roles (filtro por application_id)
- `GET /admin/roles/{id}` - Obter role específica
- `PUT /admin/roles/{id}` - Atualizar role
- `DELETE /admin/roles/{id}` - Deletar role

#### Admin - Permissions (requer autenticação)
- `POST /admin/permissions` - Criar permissão
- `GET /admin/permissions` - Listar permissões (filtro por application_id)
- `GET /admin/permissions/{id}` - Obter permissão específica
- `PUT /admin/permissions/{id}` - Atualizar permissão
- `DELETE /admin/permissions/{id}` - Deletar permissão

#### Admin - Relacionamentos Role-Permission (requer autenticação)
- `POST /admin/roles/{id}/permissions` - Atribuir permissão a role
- `DELETE /admin/roles/{id}/permissions/{permissionId}` - Remover permissão de role
- `GET /admin/roles/{id}/permissions` - Listar permissões de uma role

#### Admin - Relacionamentos User-Role (requer autenticação)
- `POST /admin/users/{id}/roles` - Atribuir role a usuário
- `DELETE /admin/users/{id}/roles/{roleId}` - Remover role de usuário

#### Admin - Relacionamentos User-Application (requer autenticação)
- `POST /admin/users/{id}/applications` - Atribuir aplicação a usuário
- `DELETE /admin/users/{id}/applications/{applicationId}` - Remover aplicação de usuário
- `GET /admin/users/{id}/applications` - Listar aplicações do usuário

#### Admin - Relacionamentos Application-User (requer autenticação)
- `GET /admin/applications/{id}/users` - Listar usuários da aplicação

### Autenticação

Para endpoints protegidos, use o token JWT retornado pelo login:

1. Faça login em `/login` com email e senha
2. Copie o token retornado
3. No Swagger UI, clique no botão "Authorize"
4. Digite: `Bearer SEU_TOKEN_AQUI`
5. Agora você pode testar os endpoints protegidos

### Exemplo de Uso

1. **Registrar usuário**:
   ```json
   POST /register
   {
     "username": "admin",
     "email": "admin@example.com",
     "password": "123456"
   }
   ```

2. **Fazer login**:
   ```json
   POST /login
   {
     "email": "admin@example.com",
     "password": "123456"
   }
   ```

3. **Criar aplicação** (com token):
   ```json
   POST /admin/applications
   {
     "slug": "my-app",
     "name": "Minha Aplicação",
     "description": "Descrição da aplicação"
   }
   ```

4. **Atribuir aplicação ao usuário** (com token):
   ```json
   POST /admin/users/{user_id}/applications
   {
     "application_id": 1
   }
   ```

5. **Listar aplicações do usuário** (com token):
   ```
   GET /admin/users/{user_id}/applications
   ```

### Regenerando a Documentação

Se você modificar as anotações Swagger nos handlers, regenere a documentação:

```bash
~/go/bin/swag init --dir ./ --parseDependency --parseInternal
```

### Estrutura dos Arquivos Swagger

- `docs/docs.go` - Código Go gerado automaticamente
- `docs/swagger.json` - Schema JSON da API
- `docs/swagger.yaml` - Schema YAML da API

**⚠️ Importante**: Os arquivos na pasta `docs/` são gerados automaticamente. Não edite-os manualmente.
