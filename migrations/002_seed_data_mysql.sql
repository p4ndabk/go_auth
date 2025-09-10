-- MySQL Seed Data
-- This file contains sample data for MySQL database

USE auth_db;

-- Insert sample applications
INSERT IGNORE INTO applications (uuid, slug, name, description, active) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'main-app', 'Main Application', 'Principal aplicação do sistema', true),
('550e8400-e29b-41d4-a716-446655440002', 'admin-app', 'Admin Panel', 'Painel administrativo', true);

-- Insert sample roles
INSERT IGNORE INTO roles (application_id, uuid, slug, name, description, active) VALUES
(1, '550e8400-e29b-41d4-a716-446655440011', 'admin', 'Administrator', 'Administrador com acesso total', true),
(1, '550e8400-e29b-41d4-a716-446655440012', 'user', 'User', 'Usuário padrão', true),
(1, '550e8400-e29b-41d4-a716-446655440013', 'moderator', 'Moderator', 'Moderador com permissões limitadas', true),
(2, '550e8400-e29b-41d4-a716-446655440014', 'admin', 'Admin Panel Admin', 'Administrador do painel', true);

-- Insert sample permissions
INSERT IGNORE INTO permissions (application_id, name, slug, description, active) VALUES
(1, 'Read Users', 'read_users', 'Visualizar usuários', true),
(1, 'Write Users', 'write_users', 'Criar/editar usuários', true),
(1, 'Delete Users', 'delete_users', 'Excluir usuários', true),
(1, 'Read Roles', 'read_roles', 'Visualizar roles', true),
(1, 'Write Roles', 'write_roles', 'Criar/editar roles', true),
(1, 'Read Permissions', 'read_permissions', 'Visualizar permissões', true),
(1, 'Write Permissions', 'write_permissions', 'Criar/editar permissões', true),
(2, 'Admin Access', 'admin_access', 'Acesso ao painel admin', true);

-- Insert sample role-permission associations
INSERT IGNORE INTO role_permissions (role_id, permission_id, application_id) VALUES
-- Admin role gets all permissions in main-app (application_id = 1)
(1, 1, 1), (1, 2, 1), (1, 3, 1), (1, 4, 1), (1, 5, 1), (1, 6, 1), (1, 7, 1),
-- User role gets basic read permissions in main-app (application_id = 1)
(2, 1, 1), (2, 4, 1), (2, 6, 1),
-- Moderator role gets read and write for users and roles in main-app (application_id = 1)
(3, 1, 1), (3, 2, 1), (3, 4, 1), (3, 5, 1), (3, 6, 1),
-- Admin Panel Admin gets admin access in admin-app (application_id = 2)
(4, 8, 2);

-- Insert admin user (password: admin123)
INSERT IGNORE INTO users (uuid, username, email, password, created_at) VALUES
('5c1daaf4-b138-4ca1-bbdd-1be4f2871d5c', 'admin', 'admin@admin.com', '$2a$14$LZeQWeCWbfVddA.RzJ1swueyAvl5JBSyLGcNft1buHBJCBHcn8dCy', NOW());

-- Assign admin user to administrator role in main application
INSERT IGNORE INTO user_application_role (user_id, application_id, role_id, profile_id, created_at) VALUES
(1, 1, 1, NULL, NOW());
