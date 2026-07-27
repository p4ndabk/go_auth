'use strict';

(function () {
  if (!Auth.require()) return;

  var applications = [];
  var selectedApplicationId = '';

  var content = Layout.mount({
    active: 'permissions',
    title: 'Permissions',
    subtitle: 'Permissões concedidas através das roles de cada aplicação',
    actions: [
      {
        label: 'Nova permission',
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

      var path = '/admin/permissions';
      if (selectedApplicationId) path += '?application_id=' + selectedApplicationId;

      var data = await api(path);
      render(data.permissions || []);
    } catch (err) {
      content.innerHTML = errorHtml(err.message);
    }
  }

  function render(permissions) {
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
    if (!permissions.length) {
      body = emptyHtml('Nenhuma permission encontrada.');
    } else {
      var rows = permissions
        .map(function (p) {
          return (
            '<tr>' +
            '<td>' + p.id + '</td>' +
            '<td><span class="fw-bold">' + escapeHtml(p.name) + '</span></td>' +
            '<td><code>' + escapeHtml(p.slug) + '</code></td>' +
            '<td>' + escapeHtml(applicationName(p.application_id)) + '</td>' +
            '<td class="text-secondary">' + escapeHtml(p.description || '—') + '</td>' +
            '<td>' + badge(p.active) + '</td>' +
            '<td class="table-actions">' +
            '<button class="btn btn-sm" data-edit="' + p.id + '"><i class="ti ti-pencil"></i></button> ' +
            '<button class="btn btn-sm btn-ghost-danger" data-delete="' + p.id + '"><i class="ti ti-trash"></i></button>' +
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

    content.querySelectorAll('[data-edit]').forEach(function (btn) {
      btn.addEventListener('click', function () {
        openForm(permissions.find(function (p) { return String(p.id) === btn.dataset.edit; }));
      });
    });

    content.querySelectorAll('[data-delete]').forEach(function (btn) {
      btn.addEventListener('click', function () {
        confirmDelete(permissions.find(function (p) { return String(p.id) === btn.dataset.delete; }));
      });
    });
  }

  function openForm(permission) {
    var editing = !!permission;

    if (!editing && !applications.length) {
      toast('Cadastre uma aplicação antes de criar permissions.', 'warning');
      return;
    }

    var fields = [];

    // application_id is only accepted on create — a permission cannot be
    // moved between applications, so it is not editable afterwards.
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
      { name: 'name', label: 'Nome', required: true, placeholder: 'Read Users' },
      {
        name: 'slug',
        label: 'Slug',
        required: true,
        placeholder: 'read_users',
        help: 'Único dentro da aplicação. É o valor que aparece no token JWT.'
      },
      { name: 'description', label: 'Descrição', type: 'textarea' }
    );

    if (editing) {
      fields.push({ name: 'active', label: 'Ativa', type: 'switch' });
    }

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
          values.application_id = Number(values.application_id);
          await api('/admin/permissions', { method: 'POST', body: values });
          toast('Permission criada.');
        }
        await load();
      }
    });
  }

  function confirmDelete(permission) {
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

  load();
})();
