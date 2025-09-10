# Melhoria: Application ID na Tabela Role-Permissions

## 🎯 Objetivo

Adicionar `application_id` na tabela `role_permissions` para melhorar a normalização e permitir consultas mais eficientes por aplicação.

## 📋 Mudanças Implementadas

### 1. Schema Database

**SQLite (`migrations/001_initial_schema.sql`):**
```sql
-- ANTES
CREATE TABLE `role_permissions` (
  `id` integer PRIMARY KEY,
  `role_id` integer,
  `permission_id` integer
);

-- DEPOIS
CREATE TABLE `role_permissions` (
  `id` integer PRIMARY KEY,
  `role_id` integer,
  `permission_id` integer,
  `application_id` integer,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(role_id, permission_id, application_id)
);
```

**MySQL (`migrations/001_initial_schema_mysql.sql`):**
```sql
-- DEPOIS
CREATE TABLE IF NOT EXISTS role_permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    application_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    UNIQUE KEY unique_role_permission_app (role_id, permission_id, application_id),
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id),
    INDEX idx_application_id (application_id)
);
```

### 2. Estrutura Go

**Novo Struct:**
```go
type RolePermission struct {
    ID            int       `json:"id"`
    RoleID        int       `json:"role_id"`
    PermissionID  int       `json:"permission_id"`
    ApplicationID int       `json:"application_id"`
    CreatedAt     time.Time `json:"created_at"`
}
```

### 3. Funções Atualizadas

**AssignPermissionToRole:**
- ✅ Valida que role e permission pertencem à mesma aplicação
- ✅ Inclui `application_id` automaticamente
- ✅ Verifica duplicatas com contexto de aplicação

**Novas Funções:**
- `GetRolePermissionDetails()` - Retorna registros completos
- `GetRolePermissionsByApplication()` - Filtra por aplicação

### 4. Dados de Seed

**SQLite:**
```sql
INSERT INTO role_permissions (role_id, permission_id, application_id, created_at) VALUES
  (1, 1, 1, datetime('now')), -- Admin role, read_users permission, main-app
  (1, 2, 1, datetime('now')), -- Admin role, write_users permission, main-app
  -- ... etc
```

**MySQL:**
```sql
INSERT IGNORE INTO role_permissions (role_id, permission_id, application_id) VALUES
  (1, 1, 1), -- Admin role, read_users permission, main-app
  (4, 8, 2), -- Admin Panel Admin, admin_access permission, admin-app
  -- ... etc
```

## ✅ Benefícios Alcançados

### 1. **Normalização Melhorada**
- Elimina redundância de dados
- Garante consistência referencial
- Melhora integridade dos dados

### 2. **Consultas Mais Eficientes**
- Filtros por aplicação diretos
- Joins otimizados
- Índices específicos

### 3. **Validação Automática**
- Verifica compatibilidade role-permission-application
- Previne associações inválidas
- Garante contexto correto

### 4. **Flexibilidade**
- Suporte a múltiplas aplicações
- Permissões específicas por contexto
- Escalabilidade melhorada

## 🔍 Exemplos de Uso

### Atribuir Permissão (com validação)
```go
// Automaticamente obtém application_id e valida compatibilidade
err := dbService.AssignPermissionToRole(roleID, permissionID)
if err != nil {
    // Erro se role e permission não pertencem à mesma aplicação
}
```

### Consultar Permissões por Aplicação
```go
// Obter permissões de uma role em aplicação específica
permissions, err := dbService.GetRolePermissionsByApplication(roleID, applicationID)

// Obter detalhes completos dos relacionamentos
details, err := dbService.GetRolePermissionDetails(roleID)
```

### Estrutura dos Dados
```
role_permissions:
+----+---------+---------------+----------------+---------------------+
| id | role_id | permission_id | application_id | created_at          |
+----+---------+---------------+----------------+---------------------+
|  1 |       1 |             1 |              1 | 2025-09-10 14:30:21 |
|  2 |       1 |             2 |              1 | 2025-09-10 14:30:21 |
|  3 |       4 |             8 |              2 | 2025-09-10 14:30:21 |
+----+---------+---------------+----------------+---------------------+
```

## 🚨 Considerações de Migração

### Para Bancos Existentes

**SQLite:**
```sql
-- Backup da tabela atual
CREATE TABLE role_permissions_backup AS SELECT * FROM role_permissions;

-- Recriar tabela com nova estrutura
DROP TABLE role_permissions;
-- Execute novo schema...

-- Migrar dados (se necessário)
INSERT INTO role_permissions (role_id, permission_id, application_id)
SELECT rp.role_id, rp.permission_id, r.application_id 
FROM role_permissions_backup rp
JOIN roles r ON rp.role_id = r.id;
```

**MySQL:**
```sql
-- Adicionar coluna
ALTER TABLE role_permissions ADD COLUMN application_id INT NOT NULL;

-- Atualizar dados existentes
UPDATE role_permissions rp 
JOIN roles r ON rp.role_id = r.id 
SET rp.application_id = r.application_id;

-- Adicionar constraints
ALTER TABLE role_permissions 
ADD FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE;
```

## 📊 Impacto

### Performance
- ✅ **Melhor**: Consultas com filtro por aplicação
- ✅ **Melhor**: Índices específicos
- ✅ **Melhor**: Joins otimizados

### Manutenibilidade  
- ✅ **Melhor**: Estrutura mais clara
- ✅ **Melhor**: Validações automáticas
- ✅ **Melhor**: Menos possibilidade de erros

### Funcionalidade
- ✅ **Nova**: Permissões por contexto de aplicação
- ✅ **Nova**: Validação de compatibilidade
- ✅ **Nova**: Consultas granulares

---

**Implementação concluída com sucesso! 🎉**

A tabela `role_permissions` agora possui `application_id`, melhorando significativamente a estrutura de dados e proporcionando melhor performance e consistência.
