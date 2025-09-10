# Go Auth API

API de autenticação e autorização em Go com JWT, suportando múltiplas aplicações com login unificado.

## Funcionalidades

- **Autenticação**: Login com usuário e senha
- **Autorização**: Sistema de roles e permissions
- **JWT**: Tokens seguros com roles e permissions incluídas
- **Multi-aplicações**: Suporte a várias aplicações com login unificado
- **SQLite**: Banco de dados simples e eficiente

## Estrutura do Projeto

```
.
├── main.go              # Ponto de entrada da aplicação
├── schema.sql           # Schema do banco de dados
├── auth/
│   └── auth.go         # Serviço de autenticação JWT
├── db/
│   └── db.go           # Queries SQL e operações de banco
├── handlers/
│   └── handlers.go     # Handlers HTTP (register, login, me)
└── config/
    └── config.go       # Configurações da aplicação
```

## Modelo de Dados

### Tabelas Principais
- **users**: Usuários do sistema
- **roles**: Roles (funções) disponíveis
- **permissions**: Permissões específicas
- **applications**: Aplicações que usam o sistema
- **user_roles**: Relacionamento usuário-role
- **role_permissions**: Relacionamento role-permissão

## Endpoints

### Endpoints Públicos

### POST /register
Registra um novo usuário.

**Request:**
```json
{
  "username": "joao",
  "email": "joao@email.com",
  "password": "123456"
}
```

**Response:**
```json
{
  "message": "User created successfully",
  "user": {
    "id": 1,
    "uuid": "123e4567-e89b-12d3-a456-426614174000",
    "username": "joao",
    "email": "joao@email.com"
  }
}
```

### POST /login
Autentica um usuário e retorna JWT.

**Request:**
```json
{
  "email": "joao@email.com",
  "password": "123456"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "uuid": "123e4567-e89b-12d3-a456-426614174000",
    "username": "joao",
    "email": "joao@email.com",
    "roles": ["admin", "user"],
    "permissions": ["read_users", "write_users", "delete_users"]
  }
}
```

### Endpoints Protegidos

### GET /me
Retorna dados do usuário autenticado.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "id": 1,
  "uuid": "123e4567-e89b-12d3-a456-426614174000",
  "username": "joao",
  "email": "joao@email.com",
  "roles": ["admin", "user"],
  "permissions": ["read_users", "write_users", "delete_users"]
}
```

### Endpoints Administrativos

Para gerenciar applications, roles, permissions e seus relacionamentos, consulte a [Documentação da API Administrativa](ADMIN_API.md).

**Principais rotas administrativas:**
- `POST /admin/applications` - Criar aplicação
- `POST /admin/roles` - Criar role
- `POST /admin/permissions` - Criar permission
- `POST /admin/roles/:id/permissions` - Associar permission à role
- `POST /admin/users/:id/roles` - Associar role ao usuário

> **Nota:** Todas as rotas administrativas requerem autenticação JWT.

## Instalação e Execução

1. **Clone o repositório:**
```bash
git clone <repo-url>
cd go_auth
```

2. **Configure as variáveis de ambiente:**
```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

3. **Execute a aplicação:**
```bash
make run
# ou diretamente:
go run main.go
```

4. **Popular com dados de exemplo (opcional):**
```bash
make seed
```

5. **Testar a API:**
```bash
make test-api        # Testa endpoints básicos
make test-admin      # Testa endpoints administrativos
```

A API estará disponível em `http://localhost:8080`

### Comandos Disponíveis

- `make build` - Compilar aplicação
- `make run` - Executar aplicação
- `make test` - Executar testes
- `make test-api` - Testar endpoints da API
- `make test-admin` - Testar endpoints administrativos
- `make clean` - Limpar arquivos de build
- `make reset-db` - Resetar banco de dados
- `make seed` - Popular banco com dados de exemplo
- `make help` - Mostrar ajuda

## Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|---------|
| `PORT` | Porta do servidor | `8080` |
| `DATABASE_URL` | URL do banco SQLite | `./auth.db` |
| `JWT_SECRET` | Chave secreta para JWT | `your-secret-key-change-this-in-production` |

## Estrutura JWT

O token JWT contém as seguintes informações:

```json
{
  "user_id": 1,
  "email": "joao@email.com",
  "roles": ["admin", "user"],
  "permissions": ["read_users", "write_users", "delete_users"],
  "exp": 1635724800,
  "iat": 1635638400
}
```

## Tecnologias Utilizadas

- **Go**: Linguagem principal
- **Gin**: Framework HTTP
- **SQLite**: Banco de dados
- **JWT**: Autenticação stateless
- **bcrypt**: Hash de senhas
- **UUID**: Identificadores únicos

## Arquitetura

A aplicação segue uma arquitetura simples em camadas:

1. **HTTP Layer** (`handlers/`): Recebe requisições HTTP e valida dados
2. **Business Layer** (`auth/`): Lógica de autenticação e autorização
3. **Data Layer** (`db/`): Operações de banco de dados
4. **Configuration** (`config/`): Gerenciamento de configurações

## Segurança

- Senhas são hasheadas com bcrypt
- Tokens JWT com expiração de 24 horas
- Validação de entrada nos endpoints
- Middleware de autenticação para rotas protegidas

## Desenvolvimento

Para adicionar novas funcionalidades:

1. **Novos endpoints**: Adicione handlers em `handlers/handlers.go`
2. **Novas queries**: Adicione métodos em `db/db.go`
3. **Nova lógica de auth**: Modifique `auth/auth.go`
4. **Novas configurações**: Atualize `config/config.go`

## Estrutura do banco de dados
[Modelo DB](https://dbdiagram.io/d/682b8f3a1227bdcb4e0742b9)
