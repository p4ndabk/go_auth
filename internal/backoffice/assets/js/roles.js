'use strict';

(function () {
  if (!Auth.require()) return;

  var applications = [];
  var selectedApplicationId = '';

  var content = Layout.mount({
    active: 'roles',
    title: 'Roles',
    subtitle: 'Papéis por aplicação e as permissions que cada um concede',
    actions: [
      {
        label: 'Nova role',
        icon: 'ti-plus',
        onClick: function () {
          openForm(null);
        }
      }
    ]
  });

  function applicationName(id) {
    var app = applications.find(function (a) { return a.id === id; });
    return app ? app.name : '#' + id;
  }

  async function load() {
    content.innerHTML = loadingHtml();
    try {
      if (!applications.length) {
        var appData = await api('/admin/applications');
        applications = appData.applications || [];
      }

      var path = '/admin/roles';
      if (selectedApplicationId) path += '?application_id=' + selectedApplicationId;

      var data = await api(path);
      render(data.roles || []);
    } catch (err) {
      content.innerHTML = errorHtml(err.message);
    }
  }

  function render(roles) {
    var filterOptions = applications
      .map(function (app) {
        var selected = String(app.id) === String(selectedApplicationId) ? ' selected' : '';
        return '<option value="' + app.id + '"' + selected + '>' + escapeHtml(app.name) + '</option>';
      })
      .join('');

    var filter =
      '<div class="card mb-3"><div class="card-body py-2">' +
      '<div class="row align-items-center g-2">' +
      '<div class="col-auto text-secondary">Aplicação</div>' +
      '<div class="col-auto"><select class="form-select" id="app-filter">' +
      '<option value="">Todas</option>' + filterOptions +
      '</select></div></div></div></div>';

    var body;
    if (!roles.length) {
      body = emptyHtml('Nenhuma role encontrada.');
    } else {
      var rows = roles
        .map(function (role) {
          return (
            '<tr>' +
            '<td>' + role.id + '</td>' +
            '<td><span class="fw-bold">' + escapeHtml(role.name) + '</span></td>' +
            '<td><code>' + escapeHtml(role.slug) + '</code></td>' +
            '<td>' + escapeHtml(applicationName(role.application_id)) + '</td>' +
            '<td class="text-secondary">' + escapeHtml(role.description || '—') + '</td>' +
            '<td>' + badge(role.active) + '</td>' +
            '<td class="table-actions">' +
            '<button class="btn btn-sm" data-perms="' + role.id + '">' +
            '<i class="ti ti-key me-1"></i>Permissions</button> ' +
            '<button class="btn btn-sm" data-edit="' + role.id + '"><i class="ti ti-pencil"></i></button> ' +
            '<button class="btn btn-sm btn-ghost-danger" data-delete="' + role.id + '"><i class="ti ti-trash"></i></button>' +
            '</td></tr>'
          );
        })
        .join('');

      body =
        '<div class="card"><div class="table-responsive">' +
        '<table class="table table-vcenter card-table">' +
        '<thead><tr><th>ID</th><th>Nome</th><th>Slug</th><th>Aplicação</th>' +
        '<th>Descrição</th><th>Status</th><th></th></tr></thead>' +
        '<tbody>' + rows + '</tbody></table></div></div>';
    }

    content.innerHTML = filter + body;

    document.getElementById('app-filter').addEventListener('change', function (e) {
      selectedApplicationId = e.target.value;
      load();
    });

    function find(id) {
      return roles.find(function (r) { return String(r.id) === id; });
    }

    content.querySelectorAll('[data-edit]').forEach(function (btn) {
      btn.addEventListener('click', function () { openForm(find(btn.dataset.edit)); });
    });
    content.querySelectorAll('[data-delete]').forEach(function (btn) {
      btn.addEventListener('click', function () { confirmDelete(find(btn.dataset.delete)); });
    });
    content.querySelectorAll('[data-perms]').forEach(function (btn) {
      btn.addEventListener('click', function () { openPermissions(find(btn.dataset.perms)); });
    });
  }

  function openForm(role) {
    var editing = !!role;

    if (!editing && !applications.length) {
      toast('Cadastre uma aplicação antes de criar roles.', 'warning');
      return;
    }

    var fields = [];

    // A role belongs to one application for its whole life — the update
    // endpoint does not accept application_id.
    if (!editing) {
      fields.push({
        name: 'application_id',
        label: 'Aplicação',
        type: 'select',
        required: true,
        placeholder: 'Selecione…',
        defaultValue: selectedApplicationId,
        options: applications.map(function (app) {
          return { value: app.id, label: app.name };
        })
      });
    }

    fields.push(
      { name: 'name', label: 'Nome', required: true, placeholder: 'Administrator' },
      {
        name: 'slug',
        label: 'Slug',
        required: true,
        placeholder: 'admin',
        help: 'Único dentro da aplicação. É o valor que aparece no token JWT.'
      },
      { name: 'description', label: 'Descrição', type: 'textarea' }
    );

    if (editing) {
      fields.push({ name: 'active', label: 'Ativa', type: 'switch' });
    }

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
          values.application_id = Number(values.application_id);
          await api('/admin/roles', { method: 'POST', body: values });
          toast('Role criada.');
        }
        await load();
      }
    });
  }

  function confirmDelete(role) {
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

  // Permission manager: lists every permission of the role's application
  // with a switch reflecting whether the role grants it. Toggling calls the
  // assign/remove endpoints immediately — there is no "save" step.
  function openPermissions(role) {
    var modal = openModal(
      '<div class="modal-header">' +
      '<h5 class="modal-title">Permissions de ' + escapeHtml(role.name) + '</h5>' +
      '<button type="button" class="btn-close" data-modal-close aria-label="Fechar"></button>' +
      '</div>' +
      '<div class="modal-body" data-perm-body>' + loadingHtml() + '</div>' +
      '<div class="modal-footer">' +
      '<button class="btn btn-primary ms-auto" data-modal-close>Concluir</button>' +
      '</div>',
      { size: 'modal-lg' }
    );

    var body = modal.root.querySelector('[data-perm-body]');

    (async function () {
      try {
        var all = await api('/admin/permissions?application_id=' + role.application_id);
        var granted = await api('/admin/roles/' + role.id + '/permissions');

        var grantedIds = (granted.permissions || []).map(function (p) { return p.id; });
        var permissions = all.permissions || [];

        if (!permissions.length) {
          body.innerHTML = emptyHtml(
            'A aplicação "' + applicationName(role.application_id) + '" ainda não tem permissions.'
          );
          return;
        }

        body.innerHTML =
          '<div class="permission-list">' +
          permissions
            .map(function (p) {
              var checked = grantedIds.indexOf(p.id) !== -1 ? ' checked' : '';
              return (
                '<label class="form-check form-switch mb-2">' +
                '<input class="form-check-input" type="checkbox" data-permission="' + p.id + '"' + checked + '>' +
                '<span class="form-check-label">' +
                '<span class="fw-bold">' + escapeHtml(p.name) + '</span> ' +
                '<code class="ms-1">' + escapeHtml(p.slug) + '</code>' +
                (p.description
                  ? '<div class="text-secondary small">' + escapeHtml(p.description) + '</div>'
                  : '') +
                '</span></label>'
              );
            })
            .join('') +
          '</div>';

        body.querySelectorAll('[data-permission]').forEach(function (input) {
          input.addEventListener('change', async function () {
            var permissionId = input.dataset.permission;
            input.disabled = true;
            try {
              if (input.checked) {
                await api('/admin/roles/' + role.id + '/permissions', {
                  method: 'POST',
                  body: { permission_id: Number(permissionId) }
                });
                toast('Permission concedida.');
              } else {
                await api('/admin/roles/' + role.id + '/permissions/' + permissionId, {
                  method: 'DELETE'
                });
                toast('Permission revogada.');
              }
            } catch (err) {
              input.checked = !input.checked; // roll the switch back
              toast(err.message, 'danger');
            } finally {
              input.disabled = false;
            }
          });
        });
      } catch (err) {
        body.innerHTML = errorHtml(err.message);
      }
    })();
  }

  load();
})();
