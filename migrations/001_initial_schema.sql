CREATE TABLE `users` (
  `id` integer PRIMARY KEY,
  `uuid` varchar(255),
  `username` varchar(255),
  `email` varchar(255),
  `password` varchar(255),
  `created_at` timestamp
);

CREATE TABLE `permissions` (
  `id` integer PRIMARY KEY,
  `application_id` integer,
  `name` varchar(255),
  `slug` varchar(255),
  `description` varchar(255),
  `active` bool
);

CREATE TABLE `roles` (
  `id` integer PRIMARY KEY,
  `application_id` integer,
  `uuid` varchar(255),
  `slug` varchar(255),
  `name` varchar(255),
  `description` varchar(255),
  `active` bool
);

CREATE TABLE `applications` (
  `id` integer PRIMARY KEY,
  `uuid` varchar(255),
  `slug` varchar(255),
  `name` varchar(255),
  `description` varchar(255),
  `active` bool
);

CREATE TABLE `role_permissions` (
  `id` integer PRIMARY KEY,
  `role_id` integer,
  `permission_id` integer
);

CREATE TABLE `user_roles` (
  `id` integer PRIMARY KEY,
  `user_id` integer,
  `role_id` integer
);

-- New unified model: User-Application-Role
-- This table replaces both user_applications and user_roles with context
CREATE TABLE `user_application_role` (
  `id` integer PRIMARY KEY,
  `user_id` integer NOT NULL,
  `application_id` integer NOT NULL,
  `role_id` integer NOT NULL,
  `profile_id` integer,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, application_id, role_id)
);
