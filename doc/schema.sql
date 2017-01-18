CREATE DATABASE IF NOT EXISTS `push_core` DEFAULT CHARACTER SET utf8;

USE `push_core`;

DROP TABLE IF EXISTS `clients`;
CREATE TABLE `clients` (
  `id` VARCHAR(36) NOT NULL, -- 一个客户端对应一个设备ID
  `gate_server_ip` VARCHAR(36) NOT NULL, -- 连接上的Gate服务器的IP
  `user_id` VARCHAR(36) NOT NULL,
  `platform` VARCHAR(30) NOT NULL, -- android,ios
  `status` TINYINT(1) NOT NULL DEFAULT '1', -- 1: oneline,0:offline
  `created_at` BIGINT(20) NOT NULL,
  `updated_at` BIGINT(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `clients_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `offline_msgs`;
CREATE TABLE `offline_msgs` (
  `id` VARCHAR(36) NOT NULL,
  `client_id` VARCHAR(36) NOT NULL,
  `user_id` VARCHAR(36) NOT NULL,
  `kind` INT(4) UNSIGNED NOT NULL, -- 消息类型
  `content` VARCHAR(500) NOT NULL,
  `extra` VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  `created_at` BIGINT(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `offline_msgs_client_id` (`client_id`),
  KEY `offline_msgs_client_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
