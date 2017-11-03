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

DROP TABLE IF EXISTS offline_msgs;
CREATE TABLE offline_msgs (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  client_id VARCHAR(36) NOT NULL,
  -- packet_id INT(2) UNSIGNED NOT NULL, -- MQTT协议规定消息ID是16bit的整型数据
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY offline_msgs_app_id_client_id(app_id,client_id),
  KEY offline_msgs_packet_id(packet_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
