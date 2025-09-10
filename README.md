# Go Auth API

API de autenticação e autorização em Go com JWT, suportando múltiplas aplicações, roles e permissions.

## 🏗️ Estrutura do Projeto

```
go_auth/
├── cmd/                    # Aplicação principal
│   └── api/               # API server
├── internal/              # Código interno da aplicação
│   ├── auth/             # Serviços de autenticação
│   ├── database/         # Conexão e migrations de banco
│   ├── handlers/         # HTTP handlers
│   └── config/           # Configurações
├── pkg/                   # Pacotes reutilizáveis
│   ├── logs/             # Sistema de logs
│   └── response/         # Padronização de respostas
├── migrations/            # Scripts de banco de dados
│   ├── 001_initial_schema.sql
│   └── 002_seed_data.sql
├── scripts/               # Scripts utilitários
│   ├── seed_db.sh        # População do banco
│   └── test_*.sh         # Scripts de teste
├── docs/                  # Documentação
│   ├── swagger.json      # Documentação da API (auto-gerada)
│   ├── swagger.yaml      # Documentação da API (auto-gerada)
│   ├── README_SWAGGER.md # Guia do Swagger
│   └── database/         # Documentação do banco
├── deployments/           # Arquivos de deployment
│   ├── Dockerfile        # Container Docker
│   └── Makefile          # Comandos de build
├── build/                 # Binários compilados
├── .env.example          # Exemplo de variáveis de ambiente
├── go.mod                # Dependências do módulo
├── go.sum                # Checksums das dependências
└── main.go               # Entrada da aplicação
```

## 🚀 Quick Start

### 1. Instalação
```bash
git clone <repository>
cd go_auth
cp .env.example .env
```

### 2. Executar aplicação
```bash
go run main.go
```

### 3. Popular banco com dados de teste
```bash
# Automático (detecta o tipo de banco)
./scripts/seed_db_auto.sh

# Ou específico para cada tipo
./scripts/seed_db.sh       # SQLite
./scripts/seed_db_mysql.sh # MySQL
```

### 4. Testar API
```bash
./scripts/test_api.sh
```

## 📖 Documentação

- **API Swagger**: http://localhost:8080/swagger/index.html
- **Guia Swagger**: [docs/README_SWAGGER.md](docs/README_SWAGGER.md)
- **Modelo User-Application-Role**: [docs/README_USER_APP_ROLE.md](docs/README_USER_APP_ROLE.md)
- **API Admin**: [docs/ADMIN_API.md](docs/ADMIN_API.md)

## 🗄️ Banco de Dados

O projeto suporta **SQLite** (padrão) e **MySQL** através de configuração de ambiente.

### Configuração do Banco

**SQLite (Padrão):**
```env
DATABASE_TYPE=sqlite
DATABASE_URL=./auth.db
```

**MySQL:**
```env
DATABASE_TYPE=mysql
DATABASE_URL=username:password@tcp(localhost:3306)/auth_db?charset=utf8mb4&parseTime=True&loc=Local
```

### Estrutura Principal

- **users** - Usuários do sistema
- **applications** - Aplicações/sistemas
- **roles** - Papéis/funções por aplicação
- **permissions** - Permissões específicas
- **user_application_role** - Relacionamento unificado usuário + aplicação + role

### Migrations

**SQLite:**
```bash
# Aplicar schema inicial
sqlite3 auth.db < migrations/001_initial_schema.sql

# Aplicar dados de exemplo
sqlite3 auth.db < migrations/002_seed_data.sql
```

**MySQL:**
```bash
# Aplicar schema inicial
mysql -u root -p < migrations/001_initial_schema_mysql.sql

# Aplicar dados de exemplo
mysql -u root -p < migrations/002_seed_data_mysql.sql
```

## 🧪 Testes

```bash
# Populate automaticamente baseado no DATABASE_TYPE
./scripts/seed_db_auto.sh

# Seed específico para SQLite
./scripts/seed_db.sh

# Seed específico para MySQL
./scripts/seed_db_mysql.sh

# Teste geral da API
./scripts/test_api.sh

# Teste da API admin
./scripts/test_admin_api.sh

# Teste do modelo user-application-role
./scripts/test_user_app_role.sh
```

## 🔧 Configuração

Variáveis de ambiente suportadas:

```env
# Servidor
PORT=8080
GIN_MODE=debug

# JWT
JWT_SECRET=your-secret-key

# Banco de Dados
DATABASE_TYPE=sqlite          # sqlite ou mysql
DATABASE_URL=./auth.db        # Para SQLite
# DATABASE_URL=user:pass@tcp(localhost:3306)/db?charset=utf8mb4&parseTime=True&loc=Local  # Para MySQL

# MySQL (apenas para scripts)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=auth_db
```

## 📊 Endpoints Principais

### Autenticação
- `POST /register` - Registrar usuário
- `POST /login` - Login
- `GET /me` - Informações do usuário

### Admin (requer autenticação)
- `GET /admin/applications` - Listar aplicações
- `GET /admin/roles` - Listar roles
- `GET /admin/permissions` - Listar permissions
- `POST /admin/users/{id}/application-roles` - Atribuir role ao usuário

## 🏷️ Tecnologias

- **Go 1.20+**
- **Gin Web Framework**
- **SQLite** / **MySQL** (configurável)
- **JWT (golang-jwt)**
- **Swagger/OpenAPI**
- **bcrypt**

## 📝 Licença

MIT License
