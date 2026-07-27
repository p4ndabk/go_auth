'use strict';

// Transport + session layer. Every page talks to the API through here.

var API_BASE = '/api';
var BACKOFFICE_BASE = '/backoffice';

var TOKEN_KEY = 'go_auth_token';
var USER_KEY = 'go_auth_user';

var Auth = {
  get token() {
    return localStorage.getItem(TOKEN_KEY);
  },

  get user() {
    var raw = localStorage.getItem(USER_KEY);
    if (!raw) return null;
    try {
      return JSON.parse(raw);
    } catch (e) {
      return null;
    }
  },

  save: function (token, user) {
    localStorage.setItem(TOKEN_KEY, token);
    localStorage.setItem(USER_KEY, JSON.stringify(user || null));
  },

  clear: function () {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  },

  logout: function () {
    this.clear();
    window.location.href = BACKOFFICE_BASE + '/login.html';
  },

  // require() sends anonymous visitors to the login page. Page scripts must
  // bail out when it returns false — the redirect is not instantaneous.
  require: function () {
    if (!this.token) {
      window.location.href = BACKOFFICE_BASE + '/login.html';
      return false;
    }
    return true;
  }
};

function ApiError(status, code, message) {
  var err = new Error(message);
  err.name = 'ApiError';
  err.status = status;
  err.code = code;
  return err;
}

// api() returns the parsed JSON body, or throws an ApiError carrying the
// code/message from the API's shared error envelope
// ({"error": {"code", "message"}} — see internal/apierror).
async function api(path, options) {
  var opts = options || {};
  var headers = { Accept: 'application/json' };
  if (opts.body !== undefined) headers['Content-Type'] = 'application/json';
  if (Auth.token) headers['Authorization'] = 'Bearer ' + Auth.token;

  var res = await fetch(API_BASE + path, {
    method: opts.method || 'GET',
    headers: headers,
    body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined
  });

  // A 401 means the token is missing, expired or invalid — but the login
  // call itself also answers 401 for bad credentials. Only bounce to the
  // login screen if we actually had a token, otherwise the login page would
  // redirect to itself and swallow the error message.
  if (res.status === 401 && Auth.token) {
    Auth.logout();
    throw ApiError(401, 'unauthorized', 'Sessão expirada. Faça login novamente.');
  }

  var text = await res.text();
  var data = null;
  if (text) {
    try {
      data = JSON.parse(text);
    } catch (e) {
      data = null;
    }
  }

  if (!res.ok) {
    var detail = (data && data.error) || {};
    throw ApiError(
      res.status,
      detail.code || 'error',
      detail.message || 'Erro inesperado (HTTP ' + res.status + ').'
    );
  }

  return data;
}
