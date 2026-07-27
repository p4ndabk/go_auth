'use strict';

(function () {
  if (!Auth.require()) return;

  var content = Layout.mount({
    active: 'users',
    title: 'Usuários',
    subtitle: 'Contas registradas. Clique em um usuário para gerenciar seus acessos.'
  });

  async function load() {
    content.innerHTML = loadingHtml();
    try {
      var data = await api('/admin/users');
      render(data.users || []);
    } catch (err) {
      content.innerHTML = errorHtml(err.message);
    }
  }

  function render(users) {
    if (!users.length) {
      content.innerHTML = emptyHtml('Nenhum usuário cadastrado.');
      return;
    }

    var rows = users
      .map(function (user) {
        var created = user.created_at ? new Date(user.created_at).toLocaleString('pt-BR') : '—';
        var href = '/backoffice/user.html?id=' + user.id;
        return (
          '<tr>' +
          '<td>' + user.id + '</td>' +
          '<td><a class="fw-bold text-reset" href="' + href + '">' +
          escapeHtml(user.username) + '</a></td>' +
          '<td>' + escapeHtml(user.email) + '</td>' +
          '<td class="text-secondary"><code>' + escapeHtml(user.uuid) + '</code></td>' +
          '<td class="text-secondary">' + escapeHtml(created) + '</td>' +
          '<td class="table-actions">' +
          '<a class="btn btn-sm btn-primary" href="' + href + '">' +
          '<i class="ti ti-settings me-1"></i>Acessos</a></td>' +
          '</tr>'
        );
      })
      .join('');

    content.innerHTML =
      '<div class="card"><div class="table-responsive">' +
      '<table class="table table-vcenter card-table">' +
      '<thead><tr><th>ID</th><th>Usuário</th><th>E-mail</th><th>UUID</th>' +
      '<th>Criado em</th><th></th></tr></thead>' +
      '<tbody>' + rows + '</tbody></table></div></div>';
  }

  load();
})();
