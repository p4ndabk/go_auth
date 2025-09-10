# Configuração Multi-Banco de Dados

Este documento detalha como configurar e usar diferentes bancos de dados no Go Auth API.

## 🎯 Bancos Suportados

- **SQLite** - Banco padrão, arquivo local
- **MySQL** - Banco relacional em servidor

## 🔧 Configuração

### SQLite (Padrão)

**Vantagens:**
- Sem configuração de servidor
- Arquivo local simples
- Ideal para desenvolvimento
- Zero configuração

**Configuração (.env):**
```env
DATABASE_TYPE=sqlite
DATABASE_URL=./auth.db
```

**Características:**
- Arquivo de banco: `auth.db`
- Schema: `migrations/001_initial_schema.sql`
- Seed: `migrations/002_seed_data.sql`

### MySQL

**Vantagens:**
- Banco robusto para produção
- Suporte a múltiplas conexões
- Ferramentas de administração
- Melhor performance

**Configuração (.env):**
```env
DATABASE_TYPE=mysql
DATABASE_URL=username:password@tcp(localhost:3306)/auth_db?charset=utf8mb4&parseTime=True&loc=Local

# Para scripts de seed
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=auth_db
```

**Características:**
- Requer servidor MySQL rodando
- Schema: `migrations/001_initial_schema_mysql.sql`
- Seed: `migrations/002_seed_data_mysql.sql`

## 🚀 Como Alternar Entre Bancos

### 1. Usando Variáveis de Ambiente

**Para SQLite:**
```bash
export DATABASE_TYPE=sqlite
export DATABASE_URL=./auth.db
go run main.go
```

**Para MySQL:**
```bash
export DATABASE_TYPE=mysql
export DATABASE_URL="user:pass@tcp(localhost:3306)/auth_db?charset=utf8mb4&parseTime=True&loc=Local"
go run main.go
```

### 2. Usando Arquivo .env

**Crie .env para SQLite:**
```env
DATABASE_TYPE=sqlite
DATABASE_URL=./auth.db
```

**Ou para MySQL:**
```env
DATABASE_TYPE=mysql
DATABASE_URL=user:pass@tcp(localhost:3306)/auth_db?charset=utf8mb4&parseTime=True&loc=Local
```

### 3. Usando Scripts de Seed Inteligentes

**Automático (detecta tipo):**
```bash
./scripts/seed_db_auto.sh
```

**Específico por tipo:**
```bash
./scripts/seed_db.sh       # SQLite
./scripts/seed_db_mysql.sh # MySQL
```

## 📋 Diferenças de Schema

### SQLite Features
- `INTEGER PRIMARY KEY` para auto-increment
- `TEXT` para strings longas
- `BOOLEAN` como INTEGER (0/1)
- Constraints básicas

### MySQL Features  
- `INT AUTO_INCREMENT PRIMARY KEY`
- `VARCHAR(n)` com tamanhos específicos
- `BOOLEAN` nativo
- Indexes otimizados
- Foreign keys com CASCADE
- Charset UTF8MB4
- Unique constraints compostas

## 🛠️ Setup MySQL Local

### 1. Instalar MySQL
```bash
# macOS
brew install mysql
brew services start mysql

# Ubuntu/Debian
sudo apt install mysql-server
sudo systemctl start mysql

# CentOS/RHEL
sudo yum install mysql-server
sudo systemctl start mysqld
```

### 2. Configurar Database
```sql
CREATE DATABASE auth_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'auth_user'@'localhost' IDENTIFIED BY 'auth_password';
GRANT ALL PRIVILEGES ON auth_db.* TO 'auth_user'@'localhost';
FLUSH PRIVILEGES;
```

### 3. Testar Conexão
```bash
mysql -u auth_user -p auth_db -e "SELECT 1;"
```

## 🧪 Testando Configurações

### Teste SQLite
```bash
# Configure
echo "DATABASE_TYPE=sqlite" > .env
echo "DATABASE_URL=./test_auth.db" >> .env

# Execute
go run main.go &
curl http://localhost:8080/swagger/
```

### Teste MySQL
```bash
# Configure  
echo "DATABASE_TYPE=mysql" > .env
echo "DATABASE_URL=user:pass@tcp(localhost:3306)/auth_db?charset=utf8mb4&parseTime=True&loc=Local" >> .env

# Execute
go run main.go &
curl http://localhost:8080/swagger/
```

## 🔍 Debugging

### Logs de Conexão
A aplicação mostra qual banco está sendo usado:
```
2025/09/10 11:06:11 Connected to sqlite database successfully
```
ou
```
2025/09/10 11:06:11 Connected to mysql database successfully
```

### Verificar Schema
**SQLite:**
```bash
sqlite3 auth.db ".schema users"
```

**MySQL:**
```bash
mysql -u user -p -e "DESCRIBE auth_db.users;"
```

### Verificar Dados
**SQLite:**
```bash
sqlite3 auth.db "SELECT * FROM applications;"
```

**MySQL:**
```bash
mysql -u user -p auth_db -e "SELECT * FROM applications;"
```

## 📝 Migration Tips

### SQLite → MySQL
1. Use script de migração específico
2. Ajuste tipos de dados (TEXT → VARCHAR)
3. Configure charset UTF8MB4
4. Teste foreign keys

### MySQL → SQLite  
1. Simplifique constraints
2. Converta AUTO_INCREMENT → INTEGER PRIMARY KEY
3. Remova charset específicos
4. Teste com dados menores

## 🎯 Produção

### Recomendações

**Desenvolvimento:** SQLite
- Rápido setup
- Sem configuração
- Portável

**Produção:** MySQL
- Mais robusto
- Melhor performance
- Backup/restore profissional
- Monitoramento

### Configuração Docker

**docker-compose.yml:**
```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpass
      MYSQL_DATABASE: auth_db
      MYSQL_USER: auth_user
      MYSQL_PASSWORD: auth_pass
    ports:
      - "3306:3306"
      
  app:
    build: .
    environment:
      DATABASE_TYPE: mysql
      DATABASE_URL: "auth_user:auth_pass@tcp(mysql:3306)/auth_db?charset=utf8mb4&parseTime=True&loc=Local"
    depends_on:
      - mysql
```

## 🚨 Troubleshooting

### Erro: "driver not found"
- Verifique se o driver está importado no main.go
- SQLite: `_ "github.com/mattn/go-sqlite3"`
- MySQL: `_ "github.com/go-sql-driver/mysql"`

### Erro: "connection refused"
- Verifique se o servidor MySQL está rodando
- Teste conexão: `mysql -u user -p`
- Verifique host/porta na DATABASE_URL

### Erro: "table not found"
- Execute migrations apropriadas
- SQLite: `sqlite3 auth.db < migrations/001_initial_schema.sql`
- MySQL: `mysql -u user -p < migrations/001_initial_schema_mysql.sql`

---

**Configuração Multi-Banco implementada com sucesso! 🎉**
