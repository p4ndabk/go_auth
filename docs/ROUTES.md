# Estrutura de Rotas

Este projeto segue o padrão da comunidade Go para organização de rotas, mantendo o `main.go` limpo e focado apenas na inicialização da aplicação.

## Estrutura Atual

### `/cmd/api/main.go`
- **Responsabilidade**: Ponto de entrada da aplicação
- **Conteúdo**: Configuração de banco, inicialização de serviços e handlers, configuração do servidor HTTP
- **Tamanho**: ~125 linhas (anteriormente ~200+ linhas)

### `/internal/routes/router.go`
- **Responsabilidade**: Configuração e organização de todas as rotas
- **Características**:
  - Funções modulares para cada grupo de rotas
  - Separação clara entre rotas públicas, protegidas e administrativas
  - Middleware de CORS centralizado
  - Configuração do Swagger

## Organização das Rotas

### Rotas Públicas (sem autenticação)
```go
POST /register
POST /login  
GET  /health
```

### Rotas Protegidas (com autenticação)
```go
GET /me
```

### Rotas Administrativas (com autenticação)
```go
# Applications
POST   /admin/applications
GET    /admin/applications
GET    /admin/applications/:id
PUT    /admin/applications/:id
DELETE /admin/applications/:id

# Roles
POST   /admin/roles
GET    /admin/roles
GET    /admin/roles/:id
PUT    /admin/roles/:id
DELETE /admin/roles/:id

# Permissions
POST   /admin/permissions
GET    /admin/permissions
GET    /admin/permissions/:id
PUT    /admin/permissions/:id
DELETE /admin/permissions/:id

# Role-Permission relationships
POST   /admin/roles/:id/permissions
DELETE /admin/roles/:id/permissions/:permissionId
GET    /admin/roles/:id/permissions

# User-Application-Role relationships
POST   /admin/users/:id/application-roles
DELETE /admin/users/:id/application-roles/:applicationId/:roleId
GET    /admin/users/:id/application-roles
GET    /admin/application-roles/users/:userId/applications/:applicationId/roles
```

### Documentação
```go
GET /swagger/*any
```

## Benefícios desta Estrutura

### 1. **Separação de Responsabilidades**
- `main.go` focado apenas na inicialização
- `router.go` focado apenas na configuração de rotas
- Handlers focados apenas na lógica de negócio

### 2. **Manutenibilidade**
- Fácil adição de novas rotas
- Grupos de rotas claramente organizados
- Middleware centralizado

### 3. **Testabilidade**
- Rotas podem ser testadas independentemente
- Setup de routes pode ser reutilizado em testes

### 4. **Padrão da Comunidade**
- Segue as convenções estabelecidas pela comunidade Go
- Facilita onboarding de novos desenvolvedores
- Estrutura familiar para desenvolvedores Go

## Como Adicionar Novas Rotas

1. **Rotas Públicas**: Adicione em `setupPublicRoutes()`
2. **Rotas Protegidas**: Adicione em `setupProtectedRoutes()`
3. **Rotas Admin**: Crie uma nova função `setupXxxRoutes()` e chame em `setupAdminRoutes()`

## Exemplo de Adição de Novas Rotas

```go
// setupUserProfileRoutes configura as rotas de perfil do usuário
func setupUserProfileRoutes(protected *gin.RouterGroup, h *handlers.Handler) {
    protected.GET("/profile", h.GetProfile)
    protected.PUT("/profile", h.UpdateProfile)
    protected.POST("/profile/avatar", h.UploadAvatar)
}
```

Então chame em `setupProtectedRoutes()`:
```go
func setupProtectedRoutes(router *gin.Engine, h *handlers.Handler, authService *auth.Service) {
    protected := router.Group("/")
    protected.Use(authService.AuthMiddleware())
    {
        protected.GET("/me", h.Me)
        
        // Adicione aqui
        setupUserProfileRoutes(protected, h)
    }
}
```
