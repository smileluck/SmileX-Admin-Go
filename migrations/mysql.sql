-- MySQL 建表 + 种子数据（程序 AutoMigrate 之外的参考脚本）
-- 默认超级管理员：admin / 123456（上线前务必修改）
CREATE DATABASE IF NOT EXISTS smilex_admin DEFAULT CHARSET utf8mb4;
USE smilex_admin;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(64) NOT NULL,
  password VARCHAR(128) NOT NULL,
  nickname VARCHAR(64) DEFAULT '',
  email VARCHAR(128) DEFAULT '',
  status INT DEFAULT 1,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  UNIQUE KEY uk_username (username),
  KEY idx_deleted (deleted_at)
);

CREATE TABLE IF NOT EXISTS roles (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  remark VARCHAR(255) DEFAULT '',
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  UNIQUE KEY uk_code (code),
  KEY idx_deleted (deleted_at)
);

CREATE TABLE IF NOT EXISTS permissions (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  type VARCHAR(16) DEFAULT 'api',
  method VARCHAR(16) DEFAULT '',
  path VARCHAR(255) DEFAULT '',
  parent_id BIGINT UNSIGNED DEFAULT 0,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  UNIQUE KEY uk_code (code),
  KEY idx_deleted (deleted_at)
);

CREATE TABLE IF NOT EXISTS user_roles (
  user_id BIGINT UNSIGNED NOT NULL,
  role_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id BIGINT UNSIGNED NOT NULL,
  permission_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (role_id, permission_id)
);

-- 种子数据（密码哈希由程序启动时自动补全，此处的哈希对应 123456）
-- INSERT INTO permissions (id, name, code, type, method, path) VALUES (1, '全部权限', 'all', 'api', '*', '*');
-- INSERT INTO roles (id, name, code, remark) VALUES (1, '超级管理员', 'super_admin', '拥有全部权限');
-- INSERT INTO users (id, username, password, nickname, status) VALUES (1, 'admin', '$2a$10$...', '超级管理员', 1);
-- INSERT INTO user_roles VALUES (1, 1);
-- INSERT INTO role_permissions VALUES (1, 1);
