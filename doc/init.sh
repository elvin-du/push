#!/bin/sh

###########################
# Configuration
###########################

HOST=127.0.0.1
USER=root
PASSWORD=
TABLE_APPS="
DROP TABLE IF EXISTS apps;
CREATE TABLE apps (
  id VARCHAR(10) NOT NULL,
  secret VARCHAR(36) NOT NULL,
  auth_type INT(2) UNSIGNED NOT NULL, -- 1:id,secret认证
  name VARCHAR(50) NOT NULL,
  description VARCHAR(256) NOT NULL,
  status TINYINT(1) NOT NULL DEFAULT 1, -- 1: 激活,0:未激活
  created_at BIGINT(20) NOT NULL,
  updated_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
"

TABLE_CLIENTS="
DROP TABLE IF EXISTS clients;
CREATE TABLE clients (
  id VARCHAR(36) NOT NULL, -- 一个客户端对应一个设备ID
  platform VARCHAR(30) NOT NULL, -- android,ios
  gate_server_ip VARCHAR(36) NOT NULL, -- 连接上的Gate服务器的IP
  gate_server_port VARCHAR(10) NOT NULL, -- 连接上的Gate服务器的Port
  status TINYINT(1) NOT NULL DEFAULT 1, -- 1: oneline,0:offline
  created_at BIGINT(20) NOT NULL,
  updated_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY clients_platform (platform)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
"

TABLE_OFFLINE_MSG="
DROP TABLE IF EXISTS offline_msgs;
CREATE TABLE offline_msgs (
  id VARCHAR(36) NOT NULL,
  user_id VARCHAR(36) NOT NULL,
  packet_id INT(2) UNSIGNED NOT NULL, -- MQTT协议规定消息ID是16bit的整型数据
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY offline_msgs_client_id_packet_id (client_id,packet_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
"

SQL="
CREATE DATABASE IF NOT EXISTS push_core DEFAULT CHARACTER SET utf8;
USE push_core;

$TABLE_APPS;

CREATE DATABASE IF NOT EXISTS your_doctor_push DEFAULT CHARACTER SET utf8;
USE your_doctor_push;

$TABLE_CLIENTS;
$TABLE_OFFLINE_MSG;

CREATE DATABASE IF NOT EXISTS heartsquare_push DEFAULT CHARACTER SET utf8;
USE heartsquare_push;

$TABLE_CLIENTS;
$TABLE_OFFLINE_MSG;

CREATE DATABASE IF NOT EXISTS hs_for_sales_push DEFAULT CHARACTER SET utf8;
USE hs_for_sales_push;

$TABLE_CLIENTS;
$TABLE_OFFLINE_MSG;

"

usage() {
  echo "usage: init.sh"
  echo ""
  echo "init        create databases and tables"
}

init() {
  mysql -u$USER -p$PASSWORD -h$HOST -e "$SQL"
}

case "$1" in
  init)
    init
    ;;
  *)
    usage
    ;;
esac
