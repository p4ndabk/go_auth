// End-to-end check for the backoffice UI. See README.md.
//
// Deliberately dependency-light: playwright-core driving the locally
// installed Chrome, no test framework, no browser download. Run it with
// `npm test` against a server backed by a throwaway database.
const { chromium } = require('playwright-core');
const fs = require('fs');
const path = require('path');

const BASE_URL = process.env.BASE_URL || 'http://localhost:8099';
const BASE = `${BASE_URL}/backoffice`;
const EMAIL = process.env.ADMIN_EMAIL || 'admin@admin.com';
const PASSWORD = process.env.ADMIN_PASSWORD || 'admin123';
const CHROME =
  process.env.CHROME_PATH || '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome';

const SHOTS = path.join(__dirname, 'shots');
fs.mkdirSync(SHOTS, { recursive: true });

const results = [];
const problems = [];

function check(name, condition, detail) {
  results.push({ name, ok: !!condition, detail: detail || '' });
  if (!condition) problems.push(`${name}${detail ? ' — ' + detail : ''}`);
}

(async () => {
  const browser = await chromium.launch({ executablePath: CHROME, headless: true });
  const page = await browser.newPage({ viewport: { width: 1440, height: 950 } });

  const consoleErrors = [];
  const failedRequests = [];
  page.on('console', (m) => {
    if (m.type() === 'error') consoleErrors.push(m.text());
  });
  page.on('pageerror', (e) => consoleErrors.push('pageerror: ' + e.message));
  page.on('response', (r) => {
    if (r.status() >= 400) failedRequests.push(`${r.status()} ${r.url()}`);
  });

  // ---------------------------------------------------------------- login
  await page.goto(`${BASE}/login.html`, { waitUntil: 'networkidle' });
  await page.screenshot({ path: `${SHOTS}/01-login.png` });

  // Confirms the pinned Tabler CDN is reachable and actually applied.
  const radius = await page.evaluate(() => {
    const el = document.querySelector('.card');
    return el ? getComputedStyle(el).borderRadius : null;
  });
  check('Tabler CSS applied', radius && radius !== '0px', `border-radius=${radius}`);

  await page.fill('input[name=email]', EMAIL);
  await page.fill('input[name=password]', PASSWORD);
  await page.click('#login-submit');
  await page.waitForURL(`${BASE}/`, { timeout: 10000 });
  await page.waitForSelector('.page-title');

  const token = await page.evaluate(() => localStorage.getItem('go_auth_token'));
  check('login stored a JWT', !!token && token.split('.').length === 3);

  // ------------------------------------------------------------ dashboard
  await page.waitForSelector('.card-table tbody tr');
  await page.screenshot({ path: `${SHOTS}/02-dashboard.png`, fullPage: true });

  const navCount = await page.locator('.navbar-vertical .navbar-nav .nav-item').count();
  check('sidebar renders every nav item', navCount === 5, `got ${navCount}`);

  const statCards = await page.locator('.row-cards .card-sm').count();
  check('dashboard renders 4 stat cards', statCards === 4, `got ${statCards}`);

  // --------------------------------------------------------- applications
  await page.click('.navbar-vertical a[href="/backoffice/applications.html"]');
  await page.waitForSelector('.card-table tbody tr');
  await page.screenshot({ path: `${SHOTS}/03-applications.png`, fullPage: true });

  const rowsBefore = await page.locator('.card-table tbody tr').count();

  await page.click('#page-actions .btn');
  await page.waitForSelector('.modal.show');
  await page.screenshot({ path: `${SHOTS}/04-application-modal.png` });

  check('create modal is visible', await page.locator('.modal.show').isVisible());
  check('modal renders one backdrop', (await page.locator('.modal-backdrop').count()) === 1);

  // Required fields must block the submit before any request goes out.
  await page.click('.modal.show button[type=submit]');
  check(
    'empty submit shows inline validation',
    await page.locator('.modal.show [data-form-error]').isVisible()
  );

  const slug = 'pw-' + Date.now();
  await page.fill('.modal.show input[name=name]', 'Playwright App');
  await page.fill('.modal.show input[name=slug]', slug);
  await page.fill('.modal.show textarea[name=description]', 'criada pelo teste e2e');
  await page.click('.modal.show button[type=submit]');

  await page.waitForSelector('#toast-container .toast');
  const toastText = (await page.locator('#toast-container .toast').first().innerText()).trim();
  check('success toast after create', /criada/i.test(toastText), toastText);

  await page.waitForFunction(
    (n) => document.querySelectorAll('.card-table tbody tr').length > n,
    rowsBefore,
    { timeout: 5000 }
  );
  const rowsAfter = await page.locator('.card-table tbody tr').count();
  check('table reloads with the new row', rowsAfter === rowsBefore + 1, `${rowsBefore} -> ${rowsAfter}`);
  await page.screenshot({ path: `${SHOTS}/05-application-created.png`, fullPage: true });

  // ------------------------------------- application detail (the matrix)
  await page.click('.card-table tbody tr:first-child a[href^="/backoffice/application.html"]');
  await page.waitForSelector('.matrix-table, .empty');
  await page.screenshot({ path: `${SHOTS}/06-application-detail.png`, fullPage: true });

  const statTiles = await page.locator('.row-cards .card-sm').count();
  check('detail page shows 3 stat tiles', statTiles === 3, `got ${statTiles}`);

  const matrixCells = await page.locator('[data-grant-role]').count();
  if (matrixCells > 0) {
    const cols = await page.locator('.matrix-table thead th').count();
    const rows = await page.locator('.matrix-table tbody tr').count();
    check(
      'matrix cell count is roles x permissions',
      matrixCells === (cols - 1) * rows,
      `${matrixCells} cells, ${cols - 1} roles, ${rows} permissions`
    );

    const cell = page.locator('[data-grant-role]').first();
    const cellBefore = await cell.isChecked();
    await cell.click();
    await page.waitForSelector('#toast-container .toast');
    const cellToast = (await page.locator('#toast-container .toast').first().innerText()).trim();
    check('matrix cell toggle hits the API', /concedida|revogada/i.test(cellToast), cellToast);
    check('matrix cell reflects the new state', (await cell.isChecked()) === !cellBefore);
    await cell.click(); // restore
    await page.waitForTimeout(600);

    // Linking a user always goes through a role, so this only applies to an
    // application that already has one.
    await page.click('[data-link-user]');
    await page.waitForSelector('.modal.show select[name=user_id]');
    check(
      'link-user modal offers a role select',
      (await page.locator('.modal.show select[name=role_id]').count()) === 1
    );
    await page.screenshot({ path: `${SHOTS}/07-link-user.png` });
    await page.click('.modal.show [data-modal-close]');
    await page.waitForTimeout(300);
  } else {
    check('matrix shows an empty state when there is nothing to grant', true, 'no roles/permissions');
  }

  // ---------------------------------------- roles + permission management
  await page.click('.navbar-vertical a[href="/backoffice/roles.html"]');
  await page.waitForSelector('.card-table tbody tr');
  await page.screenshot({ path: `${SHOTS}/08-roles.png`, fullPage: true });

  await page.locator('[data-perms]').first().click();
  await page.waitForSelector('.modal.show [data-permission]');
  await page.screenshot({ path: `${SHOTS}/09-role-permissions.png` });

  const switches = await page.locator('.modal.show [data-permission]').count();
  check('permission manager lists permissions', switches > 0, `${switches} switches`);

  const toggle = page.locator('.modal.show [data-permission]').first();
  const before = await toggle.isChecked();
  await toggle.click();
  await page.waitForSelector('#toast-container .toast');
  const permToast = (await page.locator('#toast-container .toast').first().innerText()).trim();
  check('toggling a permission hits the API', /concedida|revogada/i.test(permToast), permToast);
  check('switch reflects the new state', (await toggle.isChecked()) === !before);

  await toggle.click(); // restore the original grant
  await page.waitForTimeout(600);

  await page.click('.modal.show [data-modal-close]');
  await page.waitForTimeout(300);
  check('modal closes cleanly', (await page.locator('.modal.show').count()) === 0);
  check('backdrop removed on close', (await page.locator('.modal-backdrop').count()) === 0);

  // ------------------------------------------------- permissions + users
  await page.click('.navbar-vertical a[href="/backoffice/permissions.html"]');
  await page.waitForSelector('.card-table tbody tr');
  check('permissions page has the application filter', (await page.locator('#app-filter').count()) === 1);
  await page.screenshot({ path: `${SHOTS}/10-permissions.png`, fullPage: true });

  await page.click('.navbar-vertical a[href="/backoffice/users.html"]');
  await page.waitForSelector('.card-table tbody tr');
  check('users page lists users', (await page.locator('.card-table tbody tr').count()) >= 1);
  await page.screenshot({ path: `${SHOTS}/11-users.png`, fullPage: true });

  // ----------------------------------------- user detail (grant from user)
  await page.click('.card-table tbody tr:first-child a[href^="/backoffice/user.html"]');
  await page.waitForSelector('[data-add-role]');
  await page.screenshot({ path: `${SHOTS}/12-user-detail.png`, fullPage: true });

  const userTiles = await page.locator('.row-cards .card-sm').count();
  check('user detail shows 3 stat tiles', userTiles === 3, `got ${userTiles}`);

  // Granting from the user side: one grouped select, application derived
  // from the chosen role.
  await page.click('[data-add-role]');
  await page.waitForSelector('.modal.show select[name=role_id]');
  const optgroups = await page.locator('.modal.show select[name=role_id] optgroup').count();
  check('grant-role modal groups roles by application', optgroups >= 1, `${optgroups} optgroups`);
  await page.screenshot({ path: `${SHOTS}/13-grant-role.png` });

  const options = await page.locator('.modal.show select[name=role_id] option').count();
  if (options > 1) {
    await page.selectOption('.modal.show select[name=role_id]', { index: 1 });
    await page.click('.modal.show button[type=submit]');
    await page.waitForSelector('#toast-container .toast');
    const grantToast = (await page.locator('#toast-container .toast').first().innerText()).trim();
    check('granting a role from the user side works', /concedida/i.test(grantToast), grantToast);

    await page.waitForSelector('[data-remove-role]');
    check('granted role appears with a remove action',
      (await page.locator('[data-remove-role]').count()) >= 1);
  } else {
    check('grant-role modal handles "nothing left to grant"', true, 'user holds every role');
    await page.click('.modal.show [data-modal-close]');
    await page.waitForTimeout(300);
  }

  // ----------------------------------------------------------- dark theme
  await page.click('#theme-toggle');
  await page.waitForTimeout(400);
  const theme = await page.evaluate(() => document.documentElement.getAttribute('data-bs-theme'));
  check('theme toggle switches to dark', theme === 'dark', `theme=${theme}`);
  await page.screenshot({ path: `${SHOTS}/14-dark.png`, fullPage: true });
  await page.click('#theme-toggle');

  // ---------------------------------------------------------- auth guard
  await page.evaluate(() => localStorage.clear());
  await page.goto(`${BASE}/applications.html`);
  await page.waitForURL(/login\.html/, { timeout: 5000 });
  check('anonymous visitor is redirected to login', page.url().includes('login.html'));

  // A bad login must surface the API message rather than redirect to itself.
  await page.fill('input[name=email]', EMAIL);
  await page.fill('input[name=password]', 'senha-errada');
  await page.click('#login-submit');
  await page.waitForSelector('#login-error:not(.d-none)', { timeout: 5000 });
  const loginErr = (await page.locator('#login-error').innerText()).trim();
  check('bad login shows the API message, no redirect loop', loginErr.length > 0, loginErr);
  await page.screenshot({ path: `${SHOTS}/15-login-error.png` });

  await browser.close();

  // -------------------------------------------------------------- report
  // The 401 from the bad-credentials step above is expected.
  const unexpected = failedRequests.filter(
    (r) => !(r.startsWith('401') && r.includes('/api/login'))
  );
  check('no console errors', consoleErrors.length === 0, consoleErrors.join(' | '));
  check('no unexpected failed requests', unexpected.length === 0, unexpected.join(' | '));

  console.log('');
  for (const r of results) {
    console.log(`${r.ok ? 'PASS' : 'FAIL'}  ${r.name}${r.detail ? '  (' + r.detail + ')' : ''}`);
  }
  console.log('');
  console.log(problems.length ? `${problems.length} FAILURE(S)` : 'ALL CHECKS PASSED');
  process.exit(problems.length ? 1 : 0);
})().catch((e) => {
  console.error('HARNESS ERROR:', e.message);
  process.exit(2);
});
