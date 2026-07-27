# Frontend (backoffice)

Static admin UI for managing applications, roles and permissions, built on
[Tabler](https://tabler.io) and served by the API binary itself at
`/backoffice`. Backend conventions live in [AGENT.md](./AGENT.md) — this
file only covers the UI.

## Stack

- **Tabler 1.4.0** via jsDelivr CDN (`@tabler/core`), plus
  `@tabler/icons-webfont` for icons — no local copy, no build step.
- **Plain HTML + vanilla JS.** No framework, no bundler, no npm, no
  `node_modules`. There is nothing to install and nothing to compile: edit a
  file, reload the page.
- **`go:embed`** — the whole `assets/` tree is compiled into the binary, so
  `cmd/api` stays self-contained and the Docker image (which copies only
  binaries) needs no extra `COPY`.

This is deliberately the smallest thing that works, matching AGENT.md's
"keep it simple" rule. A SPA framework would add a toolchain, a build
artifact, and a dependency tree to a UI that is a handful of CRUD tables.
Reach for one only if the UI actually grows past that — not preemptively.

## Layout

```
internal/backoffice/
  embed.go            //go:embed all:assets
  routes.go           RegisterRoutes(r *gin.Engine) — mounts /backoffice
  assets/
    login.html        standalone (no sidebar)
    index.html        dashboard
    applications.html list of applications
    application.html  ?id=N — single-application management (the main screen)
    roles.html
    permissions.html
    users.html        list of users
    user.html         ?id=N — one user's access across applications
    css/app.css       small additions on top of Tabler
    js/
      api.js          Auth (token/session) + api() fetch wrapper
      ui.js           escapeHtml, toast, modals, form builder, placeholders
      layout.js       Layout.mount() — sidebar + page header shell
      login.js        one script per page, same name as the page
      dashboard.js
      applications.js
      roles.js
      permissions.js
      users.js
      application-detail.js
      user-detail.js
  e2e/
    backoffice.spec.js  optional browser check (not part of `go test`)
```

`backoffice` is an infra package like `health`/`docs`/`apierror`: no
`model.go`, no `service.go`. It holds **no business logic** — every page
calls the same public `/api` endpoints any other client would.

### The two detail screens

`user_application_role` is one relationship, and each side of it gets a
screen, because the two jobs are genuinely different:

- **`application.html`** — "who can do what in this system". Roles and
  permissions as a matrix (permissions down, roles across, a checkbox at each
  intersection), plus the users holding each role.
- **`user.html`** — "what can this person do". One card per application the
  user has access to, listing the roles held there and the permissions
  inherited from them.

Both write through the same endpoints, so a change on one is visible on the
other after a reload. The flat `roles.html` / `permissions.html` pages remain
for cross-application listing and filtering.

**Permissions are never assigned to a user directly.** A user gets a role
within an application and inherits that role's permissions — so `user.html`
renders permissions read-only, and every "grant" flow picks a *role*. On
`user.html` the role select is grouped by application (`<optgroup>`) and the
`application_id` is derived from the chosen role, since a role belongs to
exactly one application. That keeps it to a single select instead of two
dependent ones.

### Request fan-out

`application.html` issues one request per role
(`/admin/roles/:id/permissions`) and one per user
(`.../users/:id/applications/:id/roles`), with `Promise.all`, because the API
has no bulk endpoint for either relationship.

`user.html` avoids this: `GET /admin/users/:id/access` returns the whole
picture in one call, reusing the same aggregation the JWT is built from. When
a screen starts fanning out, prefer adding an endpoint like that one over
looping in the browser.

## Routing exception

`RegisterRoutes` takes the `*gin.Engine` and mounts on the root, **not** on
the `/api` group that AGENT.md requires for domains. That is intentional:
the backoffice is a UI, not part of the API surface, and its HTML/CSS/JS
should not sit under `/api`. It is the only package allowed to do this.

## The `.html` files are templates, not content

Every page except `login.html` is the same shell: CDN links, a theme
bootstrap script, and four ordered `<script defer>` tags. The only
differences are `<title>` and the last script. Copy an existing page rather
than writing one from scratch.

Scripts are `defer`red, so they run in document order after parsing — page
scripts can rely on `api.js`, `ui.js` and `layout.js` globals being ready.

## Adding a page

1. Copy `applications.html`, change `<title>` and the final `<script src>`.
2. Add the entry to `NAV_ITEMS` in `js/layout.js`.
3. Write `js/<page>.js` following this shape:

```js
(function () {
  if (!Auth.require()) return;          // bail out; the redirect is async

  var content = Layout.mount({
    active: 'things',                    // must match the NAV_ITEMS key
    title: 'Coisas',
    subtitle: '…',
    actions: [{ label: 'Nova', icon: 'ti-plus', onClick: function () {} }]
  });

  async function load() {
    content.innerHTML = loadingHtml();
    try {
      var data = await api('/admin/things');
      render(data.things || []);
    } catch (err) {
      content.innerHTML = errorHtml(err.message);
    }
  }

  load();
})();
```

Each page script is an IIFE — no module system, so keep the global surface
to the shared helpers and don't leak page state.

## Conventions

- **Always escape interpolated data.** Markup is built by string
  concatenation, so every value from the API goes through `escapeHtml()`.
  This is the one rule that bites: a missed call is an XSS hole, and the
  linter won't catch it.
- **Never read the DOM as state.** Re-fetch and re-render after a mutation
  (`await load()`), rather than patching table rows by hand.
- **Use `ui.js` helpers** — `loadingHtml()`, `emptyHtml()`, `errorHtml()`,
  `badge()`, `toast()`, `formModal()`, `confirmModal()` — so every page
  handles loading, empty and error states identically.
- **Modals and toasts are driven by CSS classes, not Tabler's JS.** The
  markup is Tabler/Bootstrap's, but `ui.js` adds `.show` / builds the
  backdrop itself. Nothing depends on what `tabler.min.js` exposes on
  `window`, so a CDN bump can't silently break dialogs. Keep it that way.
- **Portuguese in the UI, English in the code.** Labels, toasts and error
  copy are pt-BR; identifiers, comments and commits are English, matching
  the backend.
- **Pin the CDN version.** `@tabler/core@1.4.0`, never `@latest` — an
  upstream release should never change the UI without a commit.

## Talking to the API

`api(path, opts)` prefixes `/api`, attaches the bearer token, parses JSON,
and translates failures into a thrown error carrying `status`, `code` and
`message` taken from the backend's shared envelope
(`{"error": {"code", "message"}}` — see `internal/apierror`).

So handlers only need `try { … } catch (err) { toast(err.message) }`; the
message shown to the user is the one the API chose. `formModal` does this
already — throwing inside its `onSubmit` renders the message inline and
keeps the dialog open.

A `401` while a token is present clears the session and bounces to
`login.html`. A `401` *without* a token is left alone, so the login form can
show "invalid credentials" instead of redirecting to itself.

## Session

The JWT and the user object go in `localStorage` (`go_auth_token`,
`go_auth_user`). `Auth.require()` guards each page; `Auth.logout()` clears
and redirects. Theme choice persists under `go_auth_theme` and is applied by
an inline script in `<head>` to avoid a flash of the wrong theme.

localStorage is readable by any script on the origin. That's an accepted
trade-off for an internal tool with a short-lived (24h) token; if this ever
becomes internet-facing, move to an `HttpOnly` cookie and add CSRF
protection.

## Known gap: the UI is not an authorization boundary

`AuthRequired` only validates the JWT — **it does not check for an admin
role**. Any authenticated user can call `/api/admin/*` directly with curl,
whether or not the UI shows them a button. Hiding controls in the browser
would be decoration, so the backoffice doesn't pretend otherwise.

Fixing this belongs in the backend: a middleware that checks the token's
claims for the required role before the admin group. Until that exists,
treat every account as an admin account.

## Testing

No unit tests and no JS test framework: the pages are thin wrappers over API
calls that already have `service_test.go` coverage, so unit-testing them
would mostly assert that string concatenation concatenates.

What *is* worth testing is the part Go tests can't reach — that the page
still works once the JS actually runs. `e2e/backoffice.spec.js` drives a real
browser through login, the three CRUDs, the permission manager and the auth
guard, and fails on any console error. It is deliberately outside
`go test ./...` because it needs a running server and a browser:

```bash
node --check internal/backoffice/assets/js/*.js   # fast syntax check

cd internal/backoffice/e2e && npm install && npm test   # full run, see e2e/README.md
```

Run the e2e check after touching `ui.js` or `layout.js` — those are shared by
every page, so a regression there is invisible until something is clicked.

If page logic ever grows past "fetch, render, submit", that is the signal to
reconsider — both the no-unit-tests rule and the no-framework rule.
