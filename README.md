# Go Auth API

API de autenticação e autorização em Go com JWT, suportando múltiplas
aplicações, roles e permissions.

Inclui um **backoffice web** para gerenciar aplicações, roles e permissions.

Arquitetura, convenções e o "porquê" das decisões de design vivem em
[AGENT.md](./AGENT.md) — leia esse arquivo antes de mexer no código. As
convenções do frontend estão em [AGENT_FRONT.md](./AGENT_FRONT.md). Para
adicionar um novo domínio, copie [BASE_SPEC.md](./BASE_SPEC.md).

## Stack

- Go 1.25, Gin, GORM
- SQLite (`glebarez/sqlite`, pure Go) por padrão, MySQL como alternativa
- Migrations via `gormigrate` (`cmd/migrate`)
- JWT (`golang-jwt`) + bcrypt
- Swagger/OpenAPI (`swaggo`)

## Quick start

```bash
cp .env.example .env

go run ./cmd/migrate   # aplica o schema
go run ./cmd/seed      # opcional: popula com dados de exemplo (admin@admin.com / admin123)
go run ./cmd/api       # sobe a API em :8080
```

- Backoffice: http://localhost:8080/backoffice/
- Swagger UI: http://localhost:8080/api/docs/index.html

## Backoffice

Interface web para gerenciar aplicações, roles e permissions, construída com
[Tabler](https://tabler.io) e embutida no binário (`go:embed`) — não há build
de frontend, npm ou `node_modules` para instalar.

Faça login com um usuário existente (`admin@admin.com` / `admin123` se você
rodou o seed).

Há duas telas de gestão, uma para cada lado do vínculo:

**Detalhe da aplicação** (clique no nome da aplicação) — *"quem faz o quê
neste sistema"*:

- uma **matriz roles × permissions**, onde cada célula é um checkbox que
  concede ou revoga aquela permission para aquela role na hora;
- criar, editar e excluir roles e permissions da aplicação;
- vincular usuários e gerenciar as roles de cada um.

**Detalhe do usuário** (clique no nome do usuário) — *"o que esta pessoa
pode fazer"*:

- um card por aplicação que o usuário acessa, com as roles que ele tem ali e
  as permissions herdadas delas;
- conceder uma role (escolhida numa lista agrupada por aplicação) ou remover.

As telas de Roles e Permissions no menu servem para visão entre aplicações
(com filtro por aplicação).

### Como um usuário ganha permissões

O vínculo é sempre **usuário → (aplicação + role) → permissions**. Não existe
ligação direta entre usuário e permission: você dá uma role ao usuário dentro
de uma aplicação, e ele herda as permissions daquela role. É o que a tabela
`user_application_role` registra, e é o que aparece no JWT no login.

Por isso, nas duas telas, toda ação de "conceder" escolhe uma **role** — as
permissions aparecem como consequência, nunca como algo que se atribui
direto a alguém.

> **Atenção:** o middleware de autenticação valida o JWT mas **não** exige
> role de admin — qualquer usuário autenticado consegue chamar `/api/admin/*`
> direto. O backoffice não muda isso. Veja "Known gap" em
> [AGENT_FRONT.md](./AGENT_FRONT.md).

## Comandos

```bash
go run ./cmd/api             # sobe o servidor
go run ./cmd/migrate         # aplica migrations pendentes
go run ./cmd/migrate -rollback  # desfaz a última migration
go run ./cmd/seed            # popula dados de exemplo
go build ./...
go vet ./...
go test ./...
go tool swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

## Configuração

Variáveis de ambiente (ver `.env.example`):

| Variável | Padrão | Descrição |
|---|---|---|
| `PORT` | `8080` | Porta HTTP |
| `DB_DRIVER` | `sqlite` | `sqlite` ou `mysql` |
| `DB_PATH` | `data/go_auth.db` | Caminho do arquivo SQLite (quando `DB_DRIVER=sqlite`) |
| `DB_DSN` | — | DSN do MySQL (quando `DB_DRIVER=mysql`) |
| `JWT_SECRET` | — | Segredo usado para assinar os tokens JWT |
| `CORS_ALLOWED_ORIGINS` | `*` | `*` ou lista separada por vírgula |

## Modelo de dados

- **users** — usuários do sistema
- **applications** — aplicações/sistemas que usam este auth
- **roles** — papéis por aplicação
- **permissions** — permissões por aplicação
- **role_permissions** — quais permissions uma role concede
- **user_application_role** — qual role um usuário tem em cada aplicação

## Endpoints principais

Todas as rotas vivem sob `/api`. A documentação completa e sempre atualizada
é o Swagger (`/api/docs/index.html`), gerado a partir das anotações em cada
`handler.go`.

- `POST /api/register`, `POST /api/login`, `GET /api/me`
- `GET /api/admin/users`
- CRUD em `/api/admin/applications`, `/api/admin/roles`, `/api/admin/permissions`
- `/api/admin/roles/{id}/permissions` — atribuir/remover/listar permissions de uma role
- `/api/admin/users/{id}/application-roles` — atribuir/remover/listar roles de um usuário por aplicação

## Licença

MIT
