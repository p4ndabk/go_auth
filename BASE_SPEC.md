# Spec: &lt;domain name&gt;

Copy this file to `internal/<domain>/SPEC.md` and fill it in before writing
any code for a new domain. Delete instructional comments as you fill in
each section.

## Objective

What is this domain for, in one or two sentences?

## Scope

What's in scope for this first version? What's explicitly deferred?

## Non-goals

What this domain will *not* do, even though it might seem related. Be
explicit — this is what keeps the domain from growing into something else.

## Data model

The GORM struct(s) that will live in `model.go`. List fields, types, and any
`gorm` tags (indexes, uniqueness) with a one-line reason for each
constraint.

## Business rules

The invariants `service.go` must enforce (uniqueness, required
relationships, valid state transitions). For each rule, name the sentinel
error it produces on violation.

## Endpoints

For each handler method: HTTP verb + path, request shape, success response,
and which sentinel errors map to which `apierror` response.

## Error cases

Table of failure scenario → sentinel error → HTTP status. Anything not
covered here falls through to `apierror.Respond`'s generic 500.

## Acceptance criteria

- [ ] `model.go`, `service.go`, `handler.go`, `routes.go` exist per
      [AGENT.md](./AGENT.md)'s domain layout
- [ ] `service_test.go` covers every business rule above using an
      in-memory sqlite DB (see `internal/user/service_test.go` for the
      pattern)
- [ ] New model registered in
      `internal/database/migrations/migrations.go`
- [ ] `RegisterRoutes` wired into `cmd/api/main.go`
- [ ] Every handler method carries `swaggo/swag` annotations; `docs/`
      regenerated (`go tool swag init -g cmd/api/main.go -o docs
      --parseDependency --parseInternal`)
