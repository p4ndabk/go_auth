'use strict';

// Single-screen management for one application: the roles x permissions
// matrix, plus which users hold which role here. Everything on this page
// writes through the same /api/admin endpoints the other pages use.
(function () {
  if (!Auth.require()) return;

  var applicationId = Number(new URLSearchParams(window.location.search).get('id'));

  var application = null;
  var roles = [];
  var permissions = [];
  var grants = {}; // roleId -> { permissionId: true }
  var users = [];
  var userRoles = {}; // userId -> [roleId]

  var content = Layout.mount({
    active: 'applications',
    title: 'Aplicação',
    subtitle: 'Roles, permissions e usuários desta aplicação'
  });

  if (!applicationId) {
    content.innerHTML = errorHtml('Nenhuma aplicação informada. Volte para a lista e escolha uma.');
    return;
  }

  function roleById(id) {
    return roles.find(function (r) { return r.id === id; });
  }

  async function load() {
    content.innerHTML = loadingHtml();
    try {
      var base = await Promise.all([
        api('/admin/applications/' + applicationId),
        api('/admin/roles?application_id=' + applicationId),
        api('/admin/permissions?application_id=' + applicationId),
        api('/admin/applications/' + applicationId + '/users')
      ]);

      application = base[0].application;
      roles = base[1].roles || [];
      permissions = base[2].permissions || [];
      users = base[3].users || [];

      // One request per role and per user — fire them together rather than
      // in series. There is no bulk endpoint for either relationship.
      var grantLists = await Promise.all(
        roles.map(function (r) { return api('/admin/roles/' + r.id + '/permissions'); })
      );
      grants = {};
      roles.forEach(function (r, i) {
        grants[r.id] = {};
        (grantLists[i].permissions || []).forEach(function (p) {
          grants[r.id][p.id] = true;
        });
      });

      var roleLists = await Promise.all(
        users.map(function (u) {
          return api(
            '/admin/application-roles/users/' + u.id + '/applications/' + applicationId + '/roles'
          );
        })
      );
      userRoles = {};
      users.forEach(function (u, i) {
        userRoles[u.id] = (roleLists[i].roles || []).map(function (r) { return r.id; });
      });

      render();
    } catch (err) {
      content.innerHTML = errorHtml(err.message);
    }
  }

  function render() {
    content.innerHTML =
      headerHtml() +
      statsHtml() +
      matrixHtml() +
      usersHtml();

    bind();
  }

  function headerHtml() {
    return (
      '<div class="mb-3">' +
      '<a href="/backoffice/applications.html" class="text-secondary">' +
      '<i class="ti ti-arrow-left me-1"></i>Todas as aplicações</a></div>' +
      '<div class="card mb-3"><div class="card-body">' +
      '<div class="row align-items-center">' +
      '<div class="col">' +
      '<h2 class="mb-1">' + escapeHtml(application.name) + ' ' + badge(application.active) + '</h2>' +
      '<div class="text-secondary"><code>' + escapeHtml(application.slug) + '</code>' +
      (application.description ? ' · ' + escapeHtml(application.description) : '') +
      '</div>' +
      '<div class="text-secondary small mt-1">UUID <code>' + escapeHtml(application.uuid) + '</code></div>' +
      '</div>' +
      '<div class="col-auto">' +
      '<button class="btn" data-edit-app><i class="ti ti-pencil me-1"></i>Editar</button>' +
      '</div></div></div></div>'
    );
  }

  function statsHtml() {
    var tiles = [
      { label: 'Roles', value: roles.length, icon: 'ti-shield-lock', color: 'green' },
      { label: 'Permissions', value: permissions.length, icon: 'ti-key', color: 'azure' },
      { label: 'Usuários', value: users.length, icon: 'ti-users', color: 'orange' }
    ];
    return (
      '<div class="row row-cards mb-3">' +
      tiles
        .map(function (t) {
          return (
            '<div class="col-sm-4"><div class="card card-sm"><div class="card-body d-flex align-items-center">' +
            '<span class="bg-' + t.color + ' text-white avatar me-3"><i class="ti ' + t.icon + '"></i></span>' +
            '<div><div class="h1 mb-0">' + t.value + '</div>' +
            '<div class="text-secondary">' + t.label + '</div></div>' +
            '</div></div></div>'
          );
        })
        .join('') +
      '</div>'
    );
  }

  // The matrix is the centrepiece: permissions down the side, roles across
  // the top, a switch at every intersection.
  function matrixHtml() {
    var header =
      '<div class="card-header">' +
      '<h3 class="card-title">Acessos</h3>' +
      '<div class="card-actions btn-list">' +
      '<button class="btn btn-sm" data-new-role><i class="ti ti-plus me-1"></i>Nova role</button>' +
      '<button class="btn btn-sm" data-new-permission><i class="ti ti-plus me-1"></i>Nova permission</button>' +
      '</div></div>';

    if (!roles.length || !permissions.length) {
      var missing = !roles.length
        ? 'Esta aplicação ainda não tem roles.'
        : 'Esta aplicação ainda não tem permissions.';
      return '<div class="card mb-3">' + header +
        '<div class="card-body">' + emptyHtml(missing) + '</div></div>';
    }

    var cols = roles
      .map(function (r) {
        return (
          '<th class="text-center align-bottom">' +
          '<div class="fw-bold">' + escapeHtml(r.name) + '</div>' +
          '<div class="text-secondary small"><code>' + escapeHtml(r.slug) + '</code></div>' +
          (r.active ? '' : '<div class="small">' + badge(false) + '</div>') +
          '<div class="btn-list justify-content-center mt-1">' +
          '<button class="btn btn-sm btn-icon" data-edit-role="' + r.id + '" title="Editar role">' +
          '<i class="ti ti-pencil"></i></button>' +
          '<button class="btn btn-sm btn-icon btn-ghost-danger" data-delete-role="' + r.id + '" title="Excluir role">' +
          '<i class="ti ti-trash"></i></button>' +
          '</div></th>'
        );
      })
      .join('');

    var rows = permissions
      .map(function (p) {
        var cells = roles
          .map(function (r) {
            var checked = grants[r.id] && grants[r.id][p.id] ? ' checked' : '';
            return (
              '<td class="text-center">' +
              '<input class="form-check-input m-0" type="checkbox"' +
              ' data-grant-role="' + r.id + '" data-grant-permission="' + p.id + '"' + checked + '>' +
              '</td>'
            );
          })
          .join('');

        return (
          '<tr><td class="matrix-row-header">' +
          '<div class="fw-bold">' + escapeHtml(p.name) + '</div>' +
          '<div class="text-secondary small"><code>' + escapeHtml(p.slug) + '</code>' +
          (p.active ? '' : ' ' + badge(false)) + '</div>' +
          '<div class="btn-list mt-1">' +
          '<button class="btn btn-sm btn-icon" data-edit-permission="' + p.id + '" title="Editar permission">' +
          '<i class="ti ti-pencil"></i></button>' +
          '<button class="btn btn-sm btn-icon btn-ghost-danger" data-delete-permission="' + p.id + '" title="Excluir permission">' +
          '<i class="ti ti-trash"></i></button>' +
          '</div></td>' + cells + '</tr>'
        );
      })
      .join('');

    return (
      '<div class="card mb-3">' + header +
      '<div class="table-responsive"><table class="table table-vcenter card-table matrix-table">' +
      '<thead><tr><th class="matrix-row-header">Permission \\ Role</th>' + cols + '</tr></thead>' +
      '<tbody>' + rows + '</tbody></table></div>' +
      '<div class="card-footer text-secondary small">' +
      'Marcar concede a permission àquela role imediatamente. Usuários herdam as permissions das roles que possuem.' +
      '</div></div>'
    );
  }

  function usersHtml() {
    var body;
    if (!users.length) {
      body = '<div class="card-body">' +
        emptyHtml('Nenhum usuário vinculado a esta aplicação.') + '</div>';
    } else {
      var rows = users
        .map(function (u) {
          var badges = (userRoles[u.id] || [])
            .map(function (id) {
              var r = roleById(id);
              return '<span class="badge bg-blue-lt me-1">' +
                escapeHtml(r ? r.name : '#' + id) + '</span>';
            })
            .join('') || '<span class="text-secondary">sem role</span>';

          return (
            '<tr>' +
            '<td><a class="fw-bold text-reset" href="/backoffice/user.html?id=' + u.id + '">' +
            escapeHtml(u.username) + '</a></td>' +
            '<td class="text-secondary">' + escapeHtml(u.email) + '</td>' +
            '<td>' + badges + '</td>' +
            '<td class="table-actions">' +
            '<button class="btn btn-sm" data-user-roles="' + u.id + '">' +
            '<i class="ti ti-shield-lock me-1"></i>Roles</button></td>' +
            '</tr>'
          );
        })
        .join('');

      body =
        '<div class="table-responsive"><table class="table table-vcenter card-table">' +
        '<thead><tr><th>Usuário</th><th>E-mail</th><th>Roles nesta aplicação</th><th></th></tr></thead>' +
        '<tbody>' + rows + '</tbody></table></div>';
    }

    return (
      '<div class="card">' +
      '<div class="card-header"><h3 class="card-title">Usuários</h3>' +
      '<div class="card-actions">' +
      '<button class="btn btn-sm btn-primary" data-link-user>' +
      '<i class="ti ti-user-plus me-1"></i>Vincular usuário</button>' +
      '</div></div>' + body + '</div>'
    );
  }

  function on(selector, handler) {
    content.querySelectorAll(selector).forEach(function (el) {
      el.addEventListener('click', function () { handler(el); });
    });
  }

  function bind() {
    on('[data-edit-app]', function () { editApplication(); });
    on('[data-new-role]', function () { editRole(null); });
    on('[data-new-permission]', function () { editPermission(null); });
    on('[data-edit-role]', function (el) { editRole(roleById(Number(el.dataset.editRole))); });
    on('[data-delete-role]', function (el) { deleteRole(roleById(Number(el.dataset.deleteRole))); });
    on('[data-edit-permission]', function (el) {
      editPermission(permissions.find(function (p) { return p.id === Number(el.dataset.editPermission); }));
    });
    on('[data-delete-permission]', function (el) {
      deletePermission(permissions.find(function (p) { return p.id === Number(el.dataset.deletePermission); }));
    });
    on('[data-link-user]', function () { linkUser(); });
    on('[data-user-roles]', function (el) {
      manageUserRoles(users.find(function (u) { return u.id === Number(el.dataset.userRoles); }));
    });

    content.querySelectorAll('[data-grant-role]').forEach(function (input) {
      input.addEventListener('change', function () { toggleGrant(input); });
    });
  }

  async function toggleGrant(input) {
    var roleId = Number(input.dataset.grantRole);
    var permissionId = Number(input.dataset.grantPermission);

    input.disabled = true;
    try {
      if (input.checked) {
        await api('/admin/roles/' + roleId + '/permissions', {
          method: 'POST',
          body: { permission_id: permissionId }
        });
        grants[roleId][permissionId] = true;
        toast('Permission concedida.');
      } else {
        await api('/admin/roles/' + roleId + '/permissions/' + permissionId, { method: 'DELETE' });
        delete grants[roleId][permissionId];
        toast('Permission revogada.');
      }
    } catch (err) {
      input.checked = !input.checked; // roll the switch back
      toast(err.message, 'danger');
    } finally {
      input.disabled = false;
    }
  }

  function editApplication() {
    formModal({
      title: 'Editar aplicação',
      fields: [
        { name: 'name', label: 'Nome', required: true },
        { name: 'slug', label: 'Slug', required: true },
        { name: 'description', label: 'Descrição', type: 'textarea' },
        { name: 'active', label: 'Ativa', type: 'switch' }
      ],
      values: application,
      onSubmit: async function (values) {
        await api('/admin/applications/' + applicationId, { method: 'PUT', body: values });
        toast('Aplicação atualizada.');
        await load();
      }
    });
  }

  function editRole(role) {
    var editing = !!role;
    var fields = [
      { name: 'name', label: 'Nome', required: true, placeholder: 'Administrator' },
      { name: 'slug', label: 'Slug', required: true, placeholder: 'admin' },
      { name: 'description', label: 'Descrição', type: 'textarea' }
    ];
    if (editing) fields.push({ name: 'active', label: 'Ativa', type: 'switch' });

    formModal({
      title: editing ? 'Editar role' : 'Nova role',
      submitLabel: editing ? 'Salvar' : 'Criar',
      fields: fields,
      values: role || {},
      onSubmit: async function (values) {
        if (editing) {
          await api('/admin/roles/' + role.id, { method: 'PUT', body: values });
          toast('Role atualizada.');
        } else {
          values.application_id = applicationId;
          await api('/admin/roles', { method: 'POST', body: values });
          toast('Role criada.');
        }
        await load();
      }
    });
  }

  function deleteRole(role) {
    confirmModal({
      title: 'Excluir role?',
      message: '"' + role.name + '" será removida. Usuários que a possuem perdem esse acesso.',
      onConfirm: async function () {
        await api('/admin/roles/' + role.id, { method: 'DELETE' });
        toast('Role excluída.');
        await load();
      }
    });
  }

  function editPermission(permission) {
    var editing = !!permission;
    var fields = [
      { name: 'name', label: 'Nome', required: true, placeholder: 'Read Users' },
      { name: 'slug', label: 'Slug', required: true, placeholder: 'read_users' },
      { name: 'description', label: 'Descrição', type: 'textarea' }
    ];
    if (editing) fields.push({ name: 'active', label: 'Ativa', type: 'switch' });

    formModal({
      title: editing ? 'Editar permission' : 'Nova permission',
      submitLabel: editing ? 'Salvar' : 'Criar',
      fields: fields,
      values: permission || {},
      onSubmit: async function (values) {
        if (editing) {
          await api('/admin/permissions/' + permission.id, { method: 'PUT', body: values });
          toast('Permission atualizada.');
        } else {
          values.application_id = applicationId;
          await api('/admin/permissions', { method: 'POST', body: values });
          toast('Permission criada.');
        }
        await load();
      }
    });
  }

  function deletePermission(permission) {
    confirmModal({
      title: 'Excluir permission?',
      message: '"' + permission.name + '" será removida de todas as roles que a concedem.',
      onConfirm: async function () {
        await api('/admin/permissions/' + permission.id, { method: 'DELETE' });
        toast('Permission excluída.');
        await load();
      }
    });
  }

  // Linking a user to an application always goes through a role — that is
  // what user_application_role records. There is no user-without-role state.
  async function linkUser() {
    if (!roles.length) {
      toast('Crie uma role antes de vincular usuários.', 'warning');
      return;
    }

    var all;
    try {
      all = (await api('/admin/users')).users || [];
    } catch (err) {
      toast(err.message, 'danger');
      return;
    }

    formModal({
      title: 'Vincular usuário',
      submitLabel: 'Vincular',
      fields: [
        {
          name: 'user_id',
          label: 'Usuário',
          type: 'select',
          required: true,
          placeholder: 'Selecione…',
          options: all.map(function (u) {
            return { value: u.id, label: u.username + ' (' + u.email + ')' };
          })
        },
        {
          name: 'role_id',
          label: 'Role',
          type: 'select',
          required: true,
          placeholder: 'Selecione…',
          help: 'O usuário herda as permissions concedidas a esta role.',
          options: roles.map(function (r) {
            return { value: r.id, label: r.name };
          })
        }
      ],
      onSubmit: async function (values) {
        await api('/admin/users/' + Number(values.user_id) + '/application-roles', {
          method: 'POST',
          body: { application_id: applicationId, role_id: Number(values.role_id) }
        });
        toast('Usuário vinculado.');
        await load();
      }
    });
  }

  function manageUserRoles(user) {
    var modal = openModal(
      '<div class="modal-header">' +
      '<h5 class="modal-title">Roles de ' + escapeHtml(user.username) + '</h5>' +
      '<button type="button" class="btn-close" data-modal-close aria-label="Fechar"></button>' +
      '</div>' +
      '<div class="modal-body">' +
      '<div class="text-secondary mb-3">Em <strong>' + escapeHtml(application.name) + '</strong></div>' +
      '<div class="permission-list">' +
      roles
        .map(function (r) {
          var has = (userRoles[user.id] || []).indexOf(r.id) !== -1 ? ' checked' : '';
          return (
            '<label class="form-check form-switch mb-2">' +
            '<input class="form-check-input" type="checkbox" data-user-role="' + r.id + '"' + has + '>' +
            '<span class="form-check-label"><span class="fw-bold">' + escapeHtml(r.name) + '</span> ' +
            '<code class="ms-1">' + escapeHtml(r.slug) + '</code></span></label>'
          );
        })
        .join('') +
      '</div></div>' +
      '<div class="modal-footer">' +
      '<button class="btn btn-primary ms-auto" data-modal-close>Concluir</button></div>',
      { size: 'modal-lg' }
    );

    var dirty = false;

    modal.root.querySelectorAll('[data-user-role]').forEach(function (input) {
      input.addEventListener('change', async function () {
        var roleId = Number(input.dataset.userRole);
        input.disabled = true;
        try {
          if (input.checked) {
            await api('/admin/users/' + user.id + '/application-roles', {
              method: 'POST',
              body: { application_id: applicationId, role_id: roleId }
            });
            toast('Role atribuída.');
          } else {
            await api(
              '/admin/users/' + user.id + '/application-roles/' + applicationId + '/' + roleId,
              { method: 'DELETE' }
            );
            toast('Role removida.');
          }
          dirty = true;
        } catch (err) {
          input.checked = !input.checked;
          toast(err.message, 'danger');
        } finally {
          input.disabled = false;
        }
      });
    });

    modal.root.querySelectorAll('[data-modal-close]').forEach(function (btn) {
      btn.addEventListener('click', function () {
        if (dirty) load();
      });
    });
  }

  load();
})();
