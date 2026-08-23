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

-- 异步导出任务记录表（产物本体在 driver 对应的存储后端；无软删，保留期清理为物理删除）
CREATE TABLE IF NOT EXISTS export_records (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT DEFAULT 0,             -- 任务归属用户
  biz VARCHAR(32) DEFAULT '',           -- 业务类型（user / login_log / op_log）
  name VARCHAR(255) DEFAULT '',         -- 展示名（兼作下载文件名）
  params TEXT,                          -- 查询条件快照（JSON）
  driver VARCHAR(16) DEFAULT '',        -- 产物落库时的存储后端
  object_key VARCHAR(512) DEFAULT '',   -- 产物对象 key
  size BIGINT DEFAULT 0,                -- 产物字节数（含 BOM）
  rows INT DEFAULT 0,                   -- 已导出数据行数（不含表头）
  status VARCHAR(16) DEFAULT 'pending', -- pending | running | done | failed
  truncated BOOLEAN DEFAULT FALSE,      -- 触及大小/行数上限被截断
  error VARCHAR(512) DEFAULT '',        -- 失败原因（成功为空）
  created_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ               -- 完成/失败时间（未结束为 NULL）
);
CREATE INDEX IF NOT EXISTS idx_export_records_user_id ON export_records (user_id);
CREATE INDEX IF NOT EXISTS idx_export_records_status ON export_records (status);
CREATE INDEX IF NOT EXISTS idx_export_records_created_at ON export_records (created_at);

-- IP 黑名单表（管理员手工维护的持久化封禁；软删即解封留痕）
CREATE TABLE IF NOT EXISTS ip_blacklist (
  id BIGSERIAL PRIMARY KEY,
  ip VARCHAR(64) NOT NULL UNIQUE,       -- 单个 IP（不支持 CIDR）
  reason VARCHAR(255) DEFAULT '',       -- 封禁原因
  expire_at TIMESTAMPTZ,                -- 过期时间（NULL 为永久封禁）
  creator_id BIGINT DEFAULT 0,
  creator_name VARCHAR(64) DEFAULT '',
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_ip_blacklist_deleted_at ON ip_blacklist (deleted_at);
