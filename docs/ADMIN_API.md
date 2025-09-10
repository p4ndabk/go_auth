# API Administrativa - Documentação

Esta documentação descreve as rotas administrativas para gerenciar applications, roles, permissions e seus relacionamentos.

## Autenticação

Todas as rotas administrativas requerem autenticação via JWT Bearer token:
```
Authorization: Bearer <jwt_token>
```

## Applications

### POST /admin/applications
Cria uma nova aplicação.

**Request:**
```json
{
  "slug": "my-app",
  "name": "My Application",
  "description": "Description of my application"
}
```

**Response:**
```json
{
  "message": "Application created successfully",
  "application": {
    "id": 1,
    "uuid": "123e4567-e89b-12d3-a456-426614174000",
    "slug": "my-app",
    "name": "My Application",
    "description": "Description of my application",
    "active": true
  }
}
```

### GET /admin/applications
Lista todas as aplicações.

**Response:**
```json
{
  "applications": [
    {
      "id": 1,
      "uuid": "123e4567-e89b-12d3-a456-426614174000",
      "slug": "my-app",
      "name": "My Application",
      "description": "Description of my application",
      "active": true
    }
  ]
}
```

### GET /admin/applications/:id
Obtém uma aplicação específica.

### PUT /admin/applications/:id
Atualiza uma aplicação.

**Request:**
```json
{
  "slug": "updated-app",
  "name": "Updated Application",
  "description": "Updated description",
  "active": true
}
```

### DELETE /admin/applications/:id
Deleta uma aplicação.

## Roles

### POST /admin/roles
Cria uma nova role.

**Request:**
```json
{
  "application_id": 1,
  "slug": "admin",
  "name": "Administrator",
  "description": "Full system administrator"
}
```

**Response:**
```json
{
  "message": "Role created successfully",
  "role": {
    "id": 1,
    "application_id": 1,
    "uuid": "123e4567-e89b-12d3-a456-426614174000",
    "slug": "admin",
    "name": "Administrator",
    "description": "Full system administrator",
    "active": true
  }
}
```

### GET /admin/roles
Lista todas as roles.

**Query Parameters:**
- `application_id` (opcional): Filtra roles por aplicação

**Examples:**
- `GET /admin/roles` - Lista todas as roles
- `GET /admin/roles?application_id=1` - Lista roles da aplicação 1

### GET /admin/roles/:id
Obtém uma role específica.

### PUT /admin/roles/:id
Atualiza uma role.

**Request:**
```json
{
  "slug": "super-admin",
  "name": "Super Administrator",
  "description": "Updated description",
  "active": true
}
```

### DELETE /admin/roles/:id
Deleta uma role.

## Permissions

### POST /admin/permissions
Cria uma nova permission.

**Request:**
```json
{
  "application_id": 1,
  "name": "Read Users",
  "slug": "read_users",
  "description": "Permission to read user data"
}
```

**Response:**
```json
{
  "message": "Permission created successfully",
  "permission": {
    "id": 1,
    "application_id": 1,
    "name": "Read Users",
    "slug": "read_users",
    "description": "Permission to read user data",
    "active": true
  }
}
```

### GET /admin/permissions
Lista todas as permissions.

**Query Parameters:**
- `application_id` (opcional): Filtra permissions por aplicação

**Examples:**
- `GET /admin/permissions` - Lista todas as permissions
- `GET /admin/permissions?application_id=1` - Lista permissions da aplicação 1

### GET /admin/permissions/:id
Obtém uma permission específica.

### PUT /admin/permissions/:id
Atualiza uma permission.

**Request:**
```json
{
  "name": "Read All Users",
  "slug": "read_all_users",
  "description": "Updated description",
  "active": true
}
```

### DELETE /admin/permissions/:id
Deleta uma permission.

## Role-Permission Relationships

### POST /admin/roles/:id/permissions
Associa uma permission a uma role.

**Request:**
```json
{
  "permission_id": 1
}
```

**Response:**
```json
{
  "message": "Permission assigned to role successfully"
}
```

### GET /admin/roles/:id/permissions
Lista todas as permissions de uma role.

**Response:**
```json
{
  "permissions": [
    {
      "id": 1,
      "application_id": 1,
      "name": "Read Users",
      "slug": "read_users",
      "description": "Permission to read user data",
      "active": true
    }
  ]
}
```

### DELETE /admin/roles/:id/permissions/:permissionId
Remove uma permission de uma role.

**Response:**
```json
{
  "message": "Permission removed from role successfully"
}
```

## User-Role Relationships

### POST /admin/users/:id/roles
Associa uma role a um usuário.

**Request:**
```json
{
  "role_id": 1
}
```

**Response:**
```json
{
  "message": "Role assigned to user successfully"
}
```

### DELETE /admin/users/:id/roles/:roleId
Remove uma role de um usuário.

**Response:**
```json
{
  "message": "Role removed from user successfully"
}
```

## Códigos de Erro

- `400` - Bad Request: Dados de entrada inválidos
- `401` - Unauthorized: Token JWT inválido ou ausente
- `404` - Not Found: Recurso não encontrado
- `409` - Conflict: Recurso já existe (ex: slug duplicado, associação já existe)
- `500` - Internal Server Error: Erro interno do servidor

## Exemplos de Uso

### Configurando uma nova aplicação completa

1. **Criar aplicação:**
```bash
curl -X POST http://localhost:8080/admin/applications \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"slug": "blog", "name": "Blog System", "description": "A blog application"}'
```

2. **Criar permissions:**
```bash
curl -X POST http://localhost:8080/admin/permissions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"application_id": 1, "name": "Write Posts", "slug": "write_posts", "description": "Can create and edit posts"}'
```

3. **Criar role:**
```bash
curl -X POST http://localhost:8080/admin/roles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"application_id": 1, "slug": "blogger", "name": "Blogger", "description": "Can write blog posts"}'
```

4. **Associar permission à role:**
```bash
curl -X POST http://localhost:8080/admin/roles/1/permissions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"permission_id": 1}'
```

5. **Associar role ao usuário:**
```bash
curl -X POST http://localhost:8080/admin/users/1/roles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"role_id": 1}'
```

Após essas operações, quando o usuário fizer login, receberá um JWT contendo a role "blogger" e a permission "write_posts".
