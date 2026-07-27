'use strict';

(function () {
  if (!Auth.require()) return;

  var content = Layout.mount({
    active: 'dashboard',
    title: 'Dashboard',
    subtitle: 'Visão geral do serviço de autenticação'
  });

  var CARDS = [
    { key: 'applications', label: 'Aplicações', icon: 'ti-apps', color: 'primary', href: '/backoffice/applications.html' },
    { key: 'roles', label: 'Roles', icon: 'ti-shield-lock', color: 'green', href: '/backoffice/roles.html' },
    { key: 'permissions', label: 'Permissions', icon: 'ti-key', color: 'azure', href: '/backoffice/permissions.html' },
    { key: 'users', label: 'Usuários', icon: 'ti-users', color: 'orange', href: '/backoffice/users.html' }
  ];

  async function load() {
    content.innerHTML = loadingHtml();

    try {
      // Four independent reads — fire them together rather than in series.
      var results = await Promise.all([
        api('/admin/applications'),
        api('/admin/roles'),
        api('/admin/permissions'),
        api('/admin/users')
      ]);

      var counts = {
        applications: (results[0].applications || []).length,
        roles: (results[1].roles || []).length,
        permissions: (results[2].permissions || []).length,
        users: (results[3].users || []).length
      };

      render(counts, results[0].applications || []);
    } catch (err) {
      content.innerHTML = errorHtml(err.message);
    }
  }

  function render(counts, applications) {
    var cards = CARDS.map(function (card) {
      return (
        '<div class="col-sm-6 col-lg-3"><div class="card card-sm">' +
        '<a class="card-body d-flex align-items-center text-reset" href="' + card.href + '">' +
        '<span class="bg-' + card.color + ' text-white avatar me-3">' +
        '<i class="ti ' + card.icon + '"></i></span>' +
        '<div><div class="h1 mb-0">' + counts[card.key] + '</div>' +
        '<div class="text-secondary">' + escapeHtml(card.label) + '</div></div>' +
        '</a></div></div>'
      );
    }).join('');

    var appRows = applications.length
      ? applications
          .map(function (app) {
            return (
              '<tr><td><a class="fw-bold text-reset" href="/backoffice/application.html?id=' +
              app.id + '">' + escapeHtml(app.name) + '</a></td>' +
              '<td><code>' + escapeHtml(app.slug) + '</code></td>' +
              '<td>' + badge(app.active) + '</td></tr>'
            );
          })
          .join('')
      : '<tr><td colspan="3" class="text-secondary text-center py-4">Nenhuma aplicação cadastrada.</td></tr>';

    content.innerHTML =
      '<div class="row row-cards mb-3">' + cards + '</div>' +
      '<div class="card">' +
      '<div class="card-header"><h3 class="card-title">Aplicações</h3></div>' +
      '<div class="table-responsive"><table class="table table-vcenter card-table">' +
      '<thead><tr><th>Nome</th><th>Slug</th><th>Status</th></tr></thead>' +
      '<tbody>' + appRows + '</tbody></table></div></div>';
  }

  load();
})();
