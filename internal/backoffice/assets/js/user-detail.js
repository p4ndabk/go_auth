'use strict';

// The user's side of the same relationship the application detail page
// manages: which role this user holds in each application, and therefore
// which permissions they inherit.
(function () {
  if (!Auth.require()) return;

  var userId = Number(new URLSearchParams(window.location.search).get('id'));

  var user = null;
  var access = []; // [{ application, roles: [slug], permissions: [slug] }]
  var applications = [];
  var roles = []; // every role, across applications

  var content = Layout.mount({
    active: 'users',
    title: 'Usuário',
    subtitle: 'Acessos deste usuário por aplicação'
  });

  if (!userId) {
    content.innerHTML = errorHtml('Nenhum usuário informado. Volte para a lista e escolha um.');
    return;
  }

  function rolesOfApplication(applicationId) {
    return roles.filter(function (r) { return r.application_id === applicationId; });
  }

  function applicationById(id) {
    return applications.find(function (a) { return a.id === id; });
  }

  async function load() {
    content.innerHTML = loadingHtml();
    try {
      var data = await Promise.all([
        api('/admin/users/' + userId),
        api('/admin/users/' + userId + '/access'),
        api('/admin/applications'),
        api('/admin/roles')
      ]);

      user = data[0].user;
      access = data[1].access || [];
      applications = data[2].applications || [];
      roles = data[3].roles || [];

      render();
    } catch (err) {
      content.innerHTML = errorHtml(err.message);
    }
  }

  function render() {
    var totalRoles = access.reduce(function (sum, a) { return sum + a.roles.length; }, 0);
    var distinctPermissions = new Set();
    access.forEach(function (a) {
      a.permissions.forEach(function (p) { distinctPermissions.add(a.application.id + ':' + p); });
    });

    content.innerHTML =
      '<div class="mb-3">' +
      '<a href="/backoffice/users.html" class="text-secondary">' +
      '<i class="ti ti-arrow-left me-1"></i>Todos os usuários</a></div>' +

      '<div class="card mb-3"><div class="card-body">' +
      '<div class="row align-items-center">' +
      '<div class="col-auto">' +
      '<span class="avatar avatar-lg bg-primary-lt">' +
      escapeHtml((user.username || '?').slice(0, 2).toUpperCase()) + '</span></div>' +
      '<div class="col">' +
      '<h2 class="mb-1">' + escapeHtml(user.username) + '</h2>' +
      '<div class="text-secondary">' + escapeHtml(user.email) + '</div>' +
      '<div class="text-secondary small mt-1">UUID <code>' + escapeHtml(user.uuid) + '</code>' +
      (user.created_at
        ? ' · criado em ' + escapeHtml(new Date(user.created_at).toLocaleString('pt-BR'))
        : '') +
      '</div></div></div></div></div>' +

      '<div class="row row-cards mb-3">' +
      statTile('Aplicações', access.length, 'ti-apps', 'primary') +
      statTile('Roles', totalRoles, 'ti-shield-lock', 'green') +
      statTile('Permissions', distinctPermissions.size, 'ti-key', 'azure') +
      '</div>' +

      '<div class="d-flex align-items-center mb-3">' +
      '<h3 class="mb-0">Acessos</h3>' +
      '<button class="btn btn-primary ms-auto" data-add-role>' +
      '<i class="ti ti-plus me-1"></i>Conceder role</button></div>' +

      (access.length ? access.map(accessCard).join('') : emptyCard());

    bind();
  }

  function statTile(label, value, icon, color) {
    return (
      '<div class="col-sm-4"><div class="card card-sm"><div class="card-body d-flex align-items-center">' +
      '<span class="bg-' + color + ' text-white avatar me-3"><i class="ti ' + icon + '"></i></span>' +
      '<div><div class="h1 mb-0">' + value + '</div>' +
      '<div class="text-secondary">' + label + '</div></div>' +
      '</div></div></div>'
    );
  }

  function emptyCard() {
    return (
      '<div class="card"><div class="card-body">' +
      emptyHtml(
        'Este usuário não tem acesso a nenhuma aplicação.',
        '<button class="btn btn-primary" data-add-role>' +
        '<i class="ti ti-plus me-1"></i>Conceder role</button>'
      ) +
      '</div></div>'
    );
  }

  // One card per application, listing the roles held there and the
  // permissions those roles grant. Permissions are derived, never assigned
  // directly — so they are shown read-only.
  function accessCard(entry) {
    var app = entry.application;

    var roleBadges = entry.roles
      .map(function (slug) {
        var role = rolesOfApplication(app.id).find(function (r) { return r.slug === slug; });
        var name = role ? role.name : slug;
        var id = role ? role.id : '';
        return (
          '<span class="badge bg-blue-lt me-2 mb-2 p-2">' + escapeHtml(name) +
          ' <code class="ms-1">' + escapeHtml(slug) + '</code>' +
          (id !== ''
            ? '<button class="btn btn-sm btn-icon btn-ghost-danger ms-1"' +
              ' data-remove-role="' + id + '" data-remove-app="' + app.id + '"' +
              ' title="Remover role"><i class="ti ti-x"></i></button>'
            : '') +
          '</span>'
        );
      })
      .join('');

    var permissionBadges = entry.permissions.length
      ? entry.permissions
          .map(function (slug) {
            return '<span class="badge bg-azure-lt me-1 mb-1"><code>' + escapeHtml(slug) + '</code></span>';
          })
          .join('')
      : '<span class="text-secondary">Nenhuma — as roles acima não concedem permissions.</span>';

    return (
      '<div class="card mb-3">' +
      '<div class="card-header">' +
      '<h3 class="card-title">' +
      '<a class="text-reset" href="/backoffice/application.html?id=' + app.id + '">' +
      escapeHtml(app.name) + '</a> ' + badge(app.active) + '</h3>' +
      '<div class="card-actions">' +
      '<button class="btn btn-sm" data-add-role-app="' + app.id + '">' +
      '<i class="ti ti-plus me-1"></i>Role nesta aplicação</button></div>' +
      '</div><div class="card-body">' +
      '<div class="mb-3"><div class="form-label">Roles</div>' + roleBadges + '</div>' +
      '<div><div class="form-label">Permissions herdadas</div>' + permissionBadges + '</div>' +
      '</div></div>'
    );
  }

  function bind() {
    content.querySelectorAll('[data-add-role]').forEach(function (el) {
      el.addEventListener('click', function () { grantRole(null); });
    });
    content.querySelectorAll('[data-add-role-app]').forEach(function (el) {
      el.addEventListener('click', function () { grantRole(Number(el.dataset.addRoleApp)); });
    });
    content.querySelectorAll('[data-remove-role]').forEach(function (el) {
      el.addEventListener('click', function () {
        removeRole(Number(el.dataset.removeApp), Number(el.dataset.removeRole));
      });
    });
  }

  // Roles the user already holds, as "applicationId:roleSlug" keys.
  function heldRoleKeys() {
    var held = {};
    access.forEach(function (entry) {
      entry.roles.forEach(function (slug) {
        held[entry.application.id + ':' + slug] = true;
      });
    });
    return held;
  }

  // A single grouped select: the application is derived from the chosen
  // role, since a role always belongs to exactly one application.
  function grantRole(onlyApplicationId) {
    var held = heldRoleKeys();

    var candidates = applications
      .filter(function (app) {
        return onlyApplicationId ? app.id === onlyApplicationId : true;
      })
      .map(function (app) {
        return {
          label: app.name,
          options: rolesOfApplication(app.id)
            .filter(function (r) { return !held[app.id + ':' + r.slug]; })
            .map(function (r) { return { value: r.id, label: r.name + ' (' + r.slug + ')' }; })
        };
      })
      .filter(function (group) { return group.options.length > 0; });

    if (!candidates.length) {
      toast(
        onlyApplicationId
          ? 'O usuário já possui todas as roles desta aplicação.'
          : 'Não há roles disponíveis para conceder. Crie uma role primeiro.',
        'warning'
      );
      return;
    }

    formModal({
      title: onlyApplicationId
        ? 'Conceder role em ' + (applicationById(onlyApplicationId) || {}).name
        : 'Conceder role',
      submitLabel: 'Conceder',
      fields: [
        {
          name: 'role_id',
          label: 'Role',
          type: 'select',
          required: true,
          placeholder: 'Selecione…',
          groups: candidates,
          help: 'O usuário passa a herdar as permissions concedidas a esta role.'
        }
      ],
      onSubmit: async function (values) {
        var roleId = Number(values.role_id);
        var role = roles.find(function (r) { return r.id === roleId; });
        if (!role) throw new Error('Role não encontrada.');

        await api('/admin/users/' + userId + '/application-roles', {
          method: 'POST',
          body: { application_id: role.application_id, role_id: roleId }
        });
        toast('Role concedida.');
        await load();
      }
    });
  }

  function removeRole(applicationId, roleId) {
    var role = roles.find(function (r) { return r.id === roleId; });
    var app = applicationById(applicationId);

    confirmModal({
      title: 'Remover role?',
      message:
        (role ? role.name : 'A role') +
        ' será removida de ' + escapeHtml(user.username) +
        ' em ' + (app ? app.name : 'esta aplicação') +
        '. Ele perde as permissions herdadas dela.',
      confirmLabel: 'Remover',
      onConfirm: async function () {
        await api(
          '/admin/users/' + userId + '/application-roles/' + applicationId + '/' + roleId,
          { method: 'DELETE' }
        );
        toast('Role removida.');
        await load();
      }
    });
  }

  load();
})();
