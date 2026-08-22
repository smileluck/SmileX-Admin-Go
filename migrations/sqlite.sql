-- SQLite 建表参考（本地开发/测试；种子数据由程序启动时写入）
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  nickname TEXT DEFAULT '',
  email TEXT DEFAULT '',
  status INTEGER DEFAULT 1,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME
);

CREATE TABLE IF NOT EXISTS roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  remark TEXT DEFAULT '',
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME
);

CREATE TABLE IF NOT EXISTS permissions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  type TEXT DEFAULT 'menu', -- menu 菜单 | button 按钮权限（api 已废弃，启动迁移自动转为 button）
  method TEXT DEFAULT '',
  path TEXT DEFAULT '',
  parent_id INTEGER DEFAULT 0,
  icon TEXT DEFAULT '',
  sort INTEGER DEFAULT 0,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME
);

CREATE TABLE IF NOT EXISTS user_roles (
  user_id INTEGER NOT NULL,
  role_id INTEGER NOT NULL,
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id INTEGER NOT NULL,
  permission_id INTEGER NOT NULL,
  PRIMARY KEY (role_id, permission_id)
);

-- 文件元数据表（对象本体在 driver 对应的存储后端）
CREATE TABLE IF NOT EXISTS files (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  driver TEXT NOT NULL,                 -- local | oss | cos | tos | minio（落库时的存储后端）
  object_key TEXT NOT NULL UNIQUE,      -- 服务端生成的对象 key
  name TEXT NOT NULL,                   -- 原始文件名
  ext TEXT DEFAULT '',
  size INTEGER DEFAULT 0,
  content_type TEXT DEFAULT '',
  uploader_id INTEGER DEFAULT 0,
  uploader_name TEXT DEFAULT '',
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME
);
CREATE INDEX IF NOT EXISTS idx_files_driver ON files (driver);
CREATE INDEX IF NOT EXISTS idx_files_ext ON files (ext);
CREATE INDEX IF NOT EXISTS idx_files_uploader_id ON files (uploader_id);
CREATE INDEX IF NOT EXISTS idx_files_deleted_at ON files (deleted_at);
