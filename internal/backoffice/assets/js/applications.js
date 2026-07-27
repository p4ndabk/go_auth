'use strict';

(function () {
  if (!Auth.require()) return;

  var content = Layout.mount({
    active: 'applications',
    title: 'Aplicações',
    subtitle: 'Sistemas que usam este serviço de autenticação',
    actions: [
      {
        label: 'Nova aplicação',
        icon: 'ti-plus',
        onClick: function () {
          openForm(null);
        }
      }
    ]
  });

  async function load() {
    content.innerHTML = loadingHtml();
    try {
      var data = await api('/admin/applications');
      render(data.applications || []);
    } catch (err) {
      content.innerHTML = errorHtml(err.message);
    }
  }

  function render(apps) {
    if (!apps.length) {
      content.innerHTML = emptyHtml(
        'Nenhuma aplicação cadastrada.',
        '<button class="btn btn-primary" data-create><i class="ti ti-plus me-1"></i>Nova aplicação</button>'
      );
      var createBtn = content.querySelector('[data-create]');
      if (createBtn) createBtn.addEventListener('click', function () { openForm(null); });
      return;
    }

    var rows = apps
      .map(function (app) {
        return (
          '<tr>' +
          '<td>' + app.id + '</td>' +
          '<td><a class="fw-bold text-reset" href="/backoffice/application.html?id=' + app.id + '">' +
          escapeHtml(app.name) + '</a></td>' +
          '<td><code>' + escapeHtml(app.slug) + '</code></td>' +
          '<td class="text-secondary">' + escapeHtml(app.description || '—') + '</td>' +
          '<td>' + badge(app.active) + '</td>' +
          '<td class="table-actions">' +
          '<a class="btn btn-sm btn-primary" href="/backoffice/application.html?id=' + app.id + '">' +
          '<i class="ti ti-settings me-1"></i>Gerenciar</a> ' +
          '<button class="btn btn-sm" data-edit="' + app.id + '"><i class="ti ti-pencil"></i></button> ' +
          '<button class="btn btn-sm btn-ghost-danger" data-delete="' + app.id + '"><i class="ti ti-trash"></i></button>' +
          '</td></tr>'
        );
      })
      .join('');

    content.innerHTML =
      '<div class="card"><div class="table-responsive">' +
      '<table class="table table-vcenter card-table">' +
      '<thead><tr><th>ID</th><th>Nome</th><th>Slug</th><th>Descrição</th><th>Status</th><th></th></tr></thead>' +
      '<tbody>' + rows + '</tbody></table></div></div>';

    content.querySelectorAll('[data-edit]').forEach(function (btn) {
      btn.addEventListener('click', function () {
        var app = apps.find(function (a) { return String(a.id) === btn.dataset.edit; });
        openForm(app);
      });
    });

    content.querySelectorAll('[data-delete]').forEach(function (btn) {
      btn.addEventListener('click', function () {
        var app = apps.find(function (a) { return String(a.id) === btn.dataset.delete; });
        confirmDelete(app);
      });
    });
  }

  function openForm(app) {
    var editing = !!app;

    var fields = [
      { name: 'name', label: 'Nome', required: true, placeholder: 'Main Application' },
      {
        name: 'slug',
        label: 'Slug',
        required: true,
        placeholder: 'main-app',
        help: 'Identificador único da aplicação, usado nas integrações.'
      },
      { name: 'description', label: 'Descrição', type: 'textarea' }
    ];

    // The create endpoint always starts an application as active; only the
    // update endpoint accepts the flag.
    if (editing) {
      fields.push({ name: 'active', label: 'Ativa', type: 'switch' });
    }

    formModal({
      title: editing ? 'Editar aplicação' : 'Nova aplicação',
      submitLabel: editing ? 'Salvar' : 'Criar',
      fields: fields,
      values: app || {},
      onSubmit: async function (values) {
        if (editing) {
          await api('/admin/applications/' + app.id, { method: 'PUT', body: values });
          toast('Aplicação atualizada.');
        } else {
          await api('/admin/applications', { method: 'POST', body: values });
          toast('Aplicação criada.');
        }
        await load();
      }
    });
  }

  function confirmDelete(app) {
    confirmModal({
      title: 'Excluir aplicação?',
      message: '"' + app.name + '" será removida. Roles e permissions vinculadas ficarão órfãs.',
      onConfirm: async function () {
        await api('/admin/applications/' + app.id, { method: 'DELETE' });
        toast('Aplicação excluída.');
        await load();
      }
    });
  }

  load();
})();
