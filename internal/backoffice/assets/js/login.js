'use strict';

(function () {
  // Already signed in? Skip the form.
  if (Auth.token) {
    window.location.href = BACKOFFICE_BASE + '/';
    return;
  }

  var form = document.getElementById('login-form');
  var errorBox = document.getElementById('login-error');
  var submit = document.getElementById('login-submit');

  form.addEventListener('submit', async function (e) {
    e.preventDefault();

    errorBox.classList.add('d-none');
    submit.disabled = true;
    submit.classList.add('btn-loading');

    try {
      var data = await api('/login', {
        method: 'POST',
        body: {
          email: form.elements.email.value.trim(),
          password: form.elements.password.value
        }
      });
      Auth.save(data.token, data.user);
      window.location.href = BACKOFFICE_BASE + '/';
    } catch (err) {
      errorBox.textContent = err.message || 'Não foi possível entrar.';
      errorBox.classList.remove('d-none');
      submit.disabled = false;
      submit.classList.remove('btn-loading');
    }
  });
})();
