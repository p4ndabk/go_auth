'use strict';

// Small UI kit: escaping, toasts, modals and form modals.
//
// Modals and toasts are driven by adding Tabler/Bootstrap's CSS classes by
// hand instead of going through the bundled Bootstrap JS. The markup is the
// documented one, so it looks native, but nothing here depends on what
// tabler.min.js happens to expose on `window` — one less thing to break on
// a CDN version bump.

function escapeHtml(value) {
  if (value === null || value === undefined) return '';
  return String(value)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;');
}

function toast(message, variant) {
  var container = document.getElementById('toast-container');
  if (!container) {
    container = document.createElement('div');
    container.id = 'toast-container';
    container.className = 'toast-container';
    document.body.appendChild(container);
  }

  var el = document.createElement('div');
  el.className = 'toast align-items-center border-0 show text-bg-' + (variant || 'success');
  el.setAttribute('role', 'alert');
  el.innerHTML =
    '<div class="d-flex">' +
    '<div class="toast-body">' + escapeHtml(message) + '</div>' +
    '<button type="button" class="btn-close me-2 m-auto" aria-label="Fechar"></button>' +
    '</div>';

  var dismiss = function () {
    el.classList.remove('show');
    setTimeout(function () { el.remove(); }, 200);
  };
  el.querySelector('.btn-close').addEventListener('click', dismiss);
  container.appendChild(el);
  setTimeout(dismiss, 4000);
}

// openModal renders `dialogHtml` (a .modal-dialog element's inner markup)
// and returns { root, close }.
function openModal(dialogHtml, options) {
  var opts = options || {};

  var backdrop = document.createElement('div');
  backdrop.className = 'modal-backdrop fade show';

  var modal = document.createElement('div');
  modal.className = 'modal modal-blur fade show';
  modal.style.display = 'block';
  modal.setAttribute('role', 'dialog');
  modal.innerHTML =
    '<div class="modal-dialog ' + (opts.size || '') + ' modal-dialog-centered" role="document">' +
    '<div class="modal-content">' + dialogHtml + '</div></div>';

  function close() {
    modal.remove();
    backdrop.remove();
    document.body.classList.remove('modal-open');
    document.removeEventListener('keydown', onKey);
  }

  function onKey(e) {
    if (e.key === 'Escape') close();
  }

  modal.addEventListener('click', function (e) {
    if (e.target === modal) close();
  });
  document.addEventListener('keydown', onKey);

  document.body.appendChild(backdrop);
  document.body.appendChild(modal);
  document.body.classList.add('modal-open');

  modal.querySelectorAll('[data-modal-close]').forEach(function (btn) {
    btn.addEventListener('click', close);
  });

  return { root: modal, close: close };
}

function fieldHtml(field, values) {
  var value = values[field.name];
  if (value === undefined || value === null) value = field.defaultValue;
  if (value === undefined || value === null) value = field.type === 'switch' ? false : '';

  var required = field.required ? ' required' : '';
  var disabled = field.disabled ? ' disabled' : '';
  var help = field.help
    ? '<small class="form-hint">' + escapeHtml(field.help) + '</small>'
    : '';
  var label =
    '<label class="form-label' + (field.required ? ' required' : '') + '">' +
    escapeHtml(field.label) +
    '</label>';

  if (field.type === 'switch') {
    return (
      '<div class="mb-3"><label class="form-check form-switch">' +
      '<input class="form-check-input" type="checkbox" name="' + field.name + '"' +
      (value ? ' checked' : '') + disabled + '>' +
      '<span class="form-check-label">' + escapeHtml(field.label) + '</span>' +
      '</label>' + help + '</div>'
    );
  }

  if (field.type === 'textarea') {
    return (
      '<div class="mb-3">' + label +
      '<textarea class="form-control" name="' + field.name + '" rows="3" placeholder="' +
      escapeHtml(field.placeholder || '') + '"' + required + disabled + '>' +
      escapeHtml(value) + '</textarea>' + help + '</div>'
    );
  }

  if (field.type === 'select') {
    function optionsHtml(list) {
      return (list || [])
        .map(function (opt) {
          var selected = String(opt.value) === String(value) ? ' selected' : '';
          return '<option value="' + escapeHtml(opt.value) + '"' + selected + '>' +
            escapeHtml(opt.label) + '</option>';
        })
        .join('');
    }

    // `groups` renders <optgroup>s — used to pick a role while showing which
    // application it belongs to.
    var options = field.groups
      ? field.groups
          .map(function (group) {
            return '<optgroup label="' + escapeHtml(group.label) + '">' +
              optionsHtml(group.options) + '</optgroup>';
          })
          .join('')
      : optionsHtml(field.options);

    return (
      '<div class="mb-3">' + label +
      '<select class="form-select" name="' + field.name + '"' + required + disabled + '>' +
      (field.placeholder ? '<option value="">' + escapeHtml(field.placeholder) + '</option>' : '') +
      options + '</select>' + help + '</div>'
    );
  }

  return (
    '<div class="mb-3">' + label +
    '<input type="' + (field.type || 'text') + '" class="form-control" name="' + field.name +
    '" value="' + escapeHtml(value) + '" placeholder="' + escapeHtml(field.placeholder || '') +
    '"' + required + disabled + '>' + help + '</div>'
  );
}

// formModal builds a modal form from a field spec and calls
// onSubmit(values) — which may be async. Throwing an ApiError from
// onSubmit shows its message inline instead of closing the modal.
function formModal(config) {
  var values = config.values || {};
  var fields = config.fields || [];

  var body = fields
    .map(function (f) {
      return fieldHtml(f, values);
    })
    .join('');

  var modal = openModal(
    '<div class="modal-header">' +
    '<h5 class="modal-title">' + escapeHtml(config.title) + '</h5>' +
    '<button type="button" class="btn-close" data-modal-close aria-label="Fechar"></button>' +
    '</div>' +
    '<form novalidate>' +
    '<div class="modal-body">' +
    '<div class="alert alert-danger d-none" data-form-error></div>' +
    body +
    '</div>' +
    '<div class="modal-footer">' +
    '<button type="button" class="btn btn-link link-secondary" data-modal-close>Cancelar</button>' +
    '<button type="submit" class="btn btn-primary ms-auto" data-submit>' +
    escapeHtml(config.submitLabel || 'Salvar') +
    '</button>' +
    '</div></form>'
  );

  var form = modal.root.querySelector('form');
  var errorBox = modal.root.querySelector('[data-form-error]');
  var submitBtn = modal.root.querySelector('[data-submit]');

  var firstInput = form.querySelector('input:not([type=checkbox]), textarea, select');
  if (firstInput) firstInput.focus();

  form.addEventListener('submit', async function (e) {
    e.preventDefault();

    var collected = {};
    fields.forEach(function (f) {
      var input = form.elements[f.name];
      if (!input) return;
      if (f.type === 'switch') {
        collected[f.name] = input.checked;
      } else if (f.type === 'number') {
        collected[f.name] = input.value === '' ? null : Number(input.value);
      } else {
        collected[f.name] = input.value.trim();
      }
    });

    var missing = fields.filter(function (f) {
      return f.required && (collected[f.name] === '' || collected[f.name] === null);
    });
    if (missing.length) {
      errorBox.textContent = 'Preencha os campos obrigatórios: ' +
        missing.map(function (f) { return f.label; }).join(', ') + '.';
      errorBox.classList.remove('d-none');
      return;
    }

    errorBox.classList.add('d-none');
    submitBtn.disabled = true;
    submitBtn.classList.add('btn-loading');

    try {
      await config.onSubmit(collected);
      modal.close();
    } catch (err) {
      errorBox.textContent = err.message || 'Erro inesperado.';
      errorBox.classList.remove('d-none');
      submitBtn.disabled = false;
      submitBtn.classList.remove('btn-loading');
    }
  });

  return modal;
}

function confirmModal(config) {
  var modal = openModal(
    '<div class="modal-body text-center py-4">' +
    '<i class="ti ti-alert-triangle icon mb-2 text-danger icon-lg"></i>' +
    '<h3>' + escapeHtml(config.title) + '</h3>' +
    '<div class="text-secondary">' + escapeHtml(config.message || '') + '</div>' +
    '</div>' +
    '<div class="modal-footer">' +
    '<div class="w-100"><div class="row">' +
    '<div class="col"><button class="btn w-100" data-modal-close>Cancelar</button></div>' +
    '<div class="col"><button class="btn btn-danger w-100" data-confirm>' +
    escapeHtml(config.confirmLabel || 'Excluir') + '</button></div>' +
    '</div></div></div>',
    { size: 'modal-sm' }
  );

  var confirmBtn = modal.root.querySelector('[data-confirm]');
  confirmBtn.addEventListener('click', async function () {
    confirmBtn.disabled = true;
    confirmBtn.classList.add('btn-loading');
    try {
      await config.onConfirm();
      modal.close();
    } catch (err) {
      modal.close();
      toast(err.message || 'Erro inesperado.', 'danger');
    }
  });

  return modal;
}

// Standard placeholders so every page renders loading/empty/error the same.
function loadingHtml() {
  return '<div class="text-center py-5"><div class="spinner-border text-primary"></div></div>';
}

function emptyHtml(message, actionHtml) {
  return (
    '<div class="empty">' +
    '<div class="empty-icon"><i class="ti ti-mood-empty icon"></i></div>' +
    '<p class="empty-title">' + escapeHtml(message) + '</p>' +
    (actionHtml ? '<div class="empty-action">' + actionHtml + '</div>' : '') +
    '</div>'
  );
}

function errorHtml(message) {
  return '<div class="alert alert-danger">' + escapeHtml(message) + '</div>';
}

function badge(active) {
  return active
    ? '<span class="badge bg-success-lt">ativo</span>'
    : '<span class="badge bg-secondary-lt">inativo</span>';
}
