CREATE DATABASE IF NOT EXISTS `push_core` DEFAULT CHARACTER SET utf8;

USE `push_core`;

DROP TABLE IF EXISTS `clients`;
CREATE TABLE `clients` (
  `id` VARCHAR(36) NOT NULL, -- 一个客户端对应一个设备ID
  `user_id` VARCHAR(36) NOT NULL,
  `platform` VARCHAR(50) NOT NULL,
  `status` TINYINT(1) NOT NULL DEFAULT '1',
  `created_at` BIGINT(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `clients_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `offline_msgs`;
CREATE TABLE `offline_msgs` (
  `id` VARCHAR(36) NOT NULL,
  `client_id` VARCHAR(36) NOT NULL,
  `content` VARCHAR(500) NOT NULL,
  `created_at` BIGINT(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `clients_client_id` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
