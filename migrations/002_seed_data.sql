-- Inserir aplicações de exemplo
INSERT INTO applications (uuid, slug, name, description, active) VALUES
  ('app-1-uuid', 'main-app', 'Main Application', 'The main application', 1),
  ('app-2-uuid', 'admin-app', 'Admin Panel', 'Administrative panel', 1);

-- Inserir roles de exemplo
INSERT INTO roles (application_id, uuid, slug, name, description, active) VALUES
  (1, 'role-admin-uuid', 'admin', 'Administrator', 'Full system access', 1),
  (1, 'role-user-uuid', 'user', 'User', 'Standard user access', 1),
  (1, 'role-moderator-uuid', 'moderator', 'Moderator', 'Moderate content', 1);

-- Inserir permissions de exemplo
INSERT INTO permissions (application_id, name, slug, description, active) VALUES
  (1, 'Read Users', 'read_users', 'Can view user information', 1),
  (1, 'Write Users', 'write_users', 'Can create and edit users', 1),
  (1, 'Delete Users', 'delete_users', 'Can delete users', 1),
  (1, 'Read Posts', 'read_posts', 'Can view posts', 1),
  (1, 'Write Posts', 'write_posts', 'Can create and edit posts', 1),
  (1, 'Delete Posts', 'delete_posts', 'Can delete posts', 1),
  (1, 'Moderate Content', 'moderate_content', 'Can moderate user content', 1);

-- Associar permissions às roles
-- Admin role gets all permissions in main-app (application_id = 1)
INSERT INTO role_permissions (role_id, permission_id, application_id, created_at) VALUES
  (1, 1, 1, datetime('now')), (1, 2, 1, datetime('now')), (1, 3, 1, datetime('now')), 
  (1, 4, 1, datetime('now')), (1, 5, 1, datetime('now')), (1, 6, 1, datetime('now')), 
  (1, 7, 1, datetime('now'));

-- User role gets basic permissions in main-app (application_id = 1)
INSERT INTO role_permissions (role_id, permission_id, application_id, created_at) VALUES
  (2, 1, 1, datetime('now')), (2, 4, 1, datetime('now')), (2, 5, 1, datetime('now'));

-- Moderator role gets moderation permissions in main-app (application_id = 1)
INSERT INTO role_permissions (role_id, permission_id, application_id, created_at) VALUES
  (3, 1, 1, datetime('now')), (3, 4, 1, datetime('now')), (3, 5, 1, datetime('now')), 
  (3, 7, 1, datetime('now'));

-- Inserir usuário admin (senha: admin123)
INSERT INTO users (uuid, username, email, password, created_at) VALUES
  ('5c1daaf4-b138-4ca1-bbdd-1be4f2871d5c', 'admin', 'admin@admin.com', '$2a$14$LZeQWeCWbfVddA.RzJ1swueyAvl5JBSyLGcNft1buHBJCBHcn8dCy', datetime('now'));

-- Associar usuário admin ao role de administrator na aplicação principal
INSERT INTO user_application_role (user_id, application_id, role_id, profile_id, created_at) VALUES
  (1, 1, 1, NULL, datetime('now'));
