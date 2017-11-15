CREATE DATABASE IF NOT EXISTS push_core DEFAULT CHARACTER SET utf8;
USE push_core;

DROP TABLE IF EXISTS apps;
CREATE TABLE apps (
  id VARCHAR(36) NOT NULL,
  secret VARCHAR(36) NOT NULL,
  auth_type INT(2) UNSIGNED NOT NULL DEFAULT 1, -- 1:id,secret认证
  name VARCHAR(50) NOT NULL,
  description VARCHAR(256) NOT NULL,
  bundle_id VARCHAR(128) NOT NULL,
  cert TEXT NOT NULL,
  cert_passwd VARCHAR(128) NOT NULL,
  cert_production TEXT NOT NULL,
  cert_passwd_production VARCHAR(128) NOT NULL,
  status TINYINT(1) NOT NULL DEFAULT 1, -- 1: 激活,0:未激活
  created_at BIGINT(20) NOT NULL,
  updated_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- Sharding
DROP TABLE IF EXISTS messages;
CREATE TABLE messages (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  reg_id VARCHAR(36) NOT NULL,
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  status TINYINT(1) NOT NULL DEFAULT 1, -- 1:未读 ,2:已读
  created_at BIGINT(20) NOT NULL,
  updated_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY messages_app_id_reg_id(app_id,reg_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- Sharding
DROP TABLE IF EXISTS registries;
CREATE TABLE registries (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  dev_token VARCHAR(64) NOT NULL,
  kind VARCHAR(16) NOT NULL, -- android,ios
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY registries_app_id(app_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
