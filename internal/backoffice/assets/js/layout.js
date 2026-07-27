'use strict';

// Renders the Tabler shell (sidebar + page header) so each page script only
// has to fill in its own content.

var NAV_ITEMS = [
  { key: 'dashboard', href: '/backoffice/', icon: 'ti-home', label: 'Dashboard' },
  { key: 'applications', href: '/backoffice/applications.html', icon: 'ti-apps', label: 'Aplicações' },
  { key: 'roles', href: '/backoffice/roles.html', icon: 'ti-shield-lock', label: 'Roles' },
  { key: 'permissions', href: '/backoffice/permissions.html', icon: 'ti-key', label: 'Permissions' },
  { key: 'users', href: '/backoffice/users.html', icon: 'ti-users', label: 'Usuários' }
];

var THEME_KEY = 'go_auth_theme';

var Theme = {
  get: function () {
    return localStorage.getItem(THEME_KEY) || 'light';
  },
  apply: function (value) {
    localStorage.setItem(THEME_KEY, value);
    document.documentElement.setAttribute('data-bs-theme', value);
  },
  toggle: function () {
    this.apply(this.get() === 'dark' ? 'light' : 'dark');
  }
};

var Layout = {
  // mount renders the shell and returns the element page scripts write into.
  mount: function (config) {
    var user = Auth.user || {};
    var initials = (user.username || '?').slice(0, 2).toUpperCase();

    var nav = NAV_ITEMS.map(function (item) {
      var active = item.key === config.active ? ' active' : '';
      return (
        '<li class="nav-item' + active + '">' +
        '<a class="nav-link" href="' + item.href + '">' +
        '<span class="nav-link-icon"><i class="ti ' + item.icon + '"></i></span>' +
        '<span class="nav-link-title">' + escapeHtml(item.label) + '</span>' +
        '</a></li>'
      );
    }).join('');

    document.body.innerHTML =
      '<div class="page">' +
      '<aside class="navbar navbar-vertical navbar-expand-lg d-print-none">' +
      '<div class="container-fluid">' +
      '<button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#sidebar-menu" aria-label="Menu">' +
      '<span class="navbar-toggler-icon"></span></button>' +
      '<h1 class="navbar-brand navbar-brand-autodark">' +
      '<a href="/backoffice/"><i class="ti ti-lock-square-rounded me-2"></i>Go Auth</a></h1>' +
      '<div class="collapse navbar-collapse" id="sidebar-menu">' +
      '<ul class="navbar-nav pt-lg-3">' + nav + '</ul>' +
      '<div class="sidebar-footer mt-lg-auto pb-3 px-2">' +
      '<div class="d-flex align-items-center gap-2 mb-2">' +
      '<span class="avatar avatar-sm bg-primary-lt">' + escapeHtml(initials) + '</span>' +
      '<div class="flex-fill overflow-hidden">' +
      '<div class="text-truncate">' + escapeHtml(user.username || 'usuário') + '</div>' +
      '<div class="text-secondary text-truncate small">' + escapeHtml(user.email || '') + '</div>' +
      '</div></div>' +
      '<div class="d-flex gap-2">' +
      '<button class="btn btn-sm flex-fill" id="theme-toggle" title="Alternar tema">' +
      '<i class="ti ti-contrast"></i></button>' +
      '<button class="btn btn-sm flex-fill" id="logout-btn"><i class="ti ti-logout me-1"></i>Sair</button>' +
      '</div></div></div></div></aside>' +
      '<div class="page-wrapper">' +
      '<div class="page-header d-print-none"><div class="container-xl">' +
      '<div class="row g-2 align-items-center">' +
      '<div class="col">' +
      '<h2 class="page-title">' + escapeHtml(config.title) + '</h2>' +
      (config.subtitle
        ? '<div class="text-secondary mt-1">' + escapeHtml(config.subtitle) + '</div>'
        : '') +
      '</div>' +
      '<div class="col-auto ms-auto d-print-none" id="page-actions"></div>' +
      '</div></div></div>' +
      '<div class="page-body"><div class="container-xl" id="page-content"></div></div>' +
      '</div></div>';

    document.getElementById('logout-btn').addEventListener('click', function () {
      Auth.logout();
    });
    document.getElementById('theme-toggle').addEventListener('click', function () {
      Theme.toggle();
    });

    var actions = document.getElementById('page-actions');
    (config.actions || []).forEach(function (action) {
      var btn = document.createElement('button');
      btn.className = 'btn btn-' + (action.variant || 'primary');
      btn.innerHTML =
        (action.icon ? '<i class="ti ' + action.icon + ' me-1"></i>' : '') +
        escapeHtml(action.label);
      btn.addEventListener('click', action.onClick);
      actions.appendChild(btn);
    });

    return document.getElementById('page-content');
  }
};
