-- PostgreSQL 建表参考（种子数据由程序启动时写入）
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(64) UNIQUE NOT NULL,
  password VARCHAR(128) NOT NULL,
  nickname VARCHAR(64) DEFAULT '',
  email VARCHAR(128) DEFAULT '',
  status INT DEFAULT 1,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS roles (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  remark VARCHAR(255) DEFAULT '',
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS permissions (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  type VARCHAR(16) DEFAULT 'menu', -- menu 菜单 | button 按钮权限（api 已废弃，启动迁移自动转为 button）
  method VARCHAR(16) DEFAULT '',
  path VARCHAR(255) DEFAULT '',
  parent_id BIGINT DEFAULT 0,
  icon VARCHAR(512) DEFAULT '',
  sort INT DEFAULT 0,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_roles (
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id BIGINT NOT NULL,
  permission_id BIGINT NOT NULL,
  PRIMARY KEY (role_id, permission_id)
);

-- 文件元数据表（对象本体在 driver 对应的存储后端）
CREATE TABLE IF NOT EXISTS files (
  id BIGSERIAL PRIMARY KEY,
  driver VARCHAR(16) NOT NULL,          -- local | oss | cos | tos | minio（落库时的存储后端）
  object_key VARCHAR(512) NOT NULL UNIQUE, -- 服务端生成的对象 key
  name VARCHAR(255) NOT NULL,           -- 原始文件名
  ext VARCHAR(16) DEFAULT '',
  size BIGINT DEFAULT 0,
  content_type VARCHAR(128) DEFAULT '',
  uploader_id BIGINT DEFAULT 0,
  uploader_name VARCHAR(64) DEFAULT '',
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_files_driver ON files (driver);
CREATE INDEX IF NOT EXISTS idx_files_ext ON files (ext);
CREATE INDEX IF NOT EXISTS idx_files_uploader_id ON files (uploader_id);
CREATE INDEX IF NOT EXISTS idx_files_deleted_at ON files (deleted_at);
