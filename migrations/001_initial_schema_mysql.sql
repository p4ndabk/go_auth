-- MySQL Schema
-- This file contains the MySQL-specific schema for the auth system

CREATE DATABASE IF NOT EXISTS auth_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE auth_db;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_username (username),
    INDEX idx_uuid (uuid)
);

-- Applications table
CREATE TABLE IF NOT EXISTS applications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_slug (slug),
    INDEX idx_active (active)
);

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    application_id INT NOT NULL,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    UNIQUE KEY unique_role_per_app (application_id, slug),
    INDEX idx_application_id (application_id),
    INDEX idx_slug (slug),
    INDEX idx_active (active)
);

-- Permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    application_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    UNIQUE KEY unique_permission_per_app (application_id, slug),
    INDEX idx_application_id (application_id),
    INDEX idx_slug (slug),
    INDEX idx_active (active)
);

-- Role-Permission relationship table
CREATE TABLE IF NOT EXISTS role_permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    application_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    UNIQUE KEY unique_role_permission_app (role_id, permission_id, application_id),
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id),
    INDEX idx_application_id (application_id)
);

-- User-Role relationship table (legacy - kept for compatibility)
CREATE TABLE IF NOT EXISTS user_roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    role_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_role (user_id, role_id)
);

-- User-Application-Role unified table (recommended approach)
CREATE TABLE IF NOT EXISTS user_application_role (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    application_id INT NOT NULL,
    role_id INT NOT NULL,
    profile_id INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_app_role (user_id, application_id, role_id),
    INDEX idx_user_id (user_id),
    INDEX idx_application_id (application_id),
    INDEX idx_role_id (role_id)
);
