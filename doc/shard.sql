CREATE DATABASE IF NOT EXISTS push_shard1 DEFAULT CHARACTER SET utf8;
USE push_shard1;

DROP TABLE IF EXISTS offline_msgs;
CREATE TABLE offline_msgs (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  reg_id VARCHAR(36) NOT NULL,
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY offline_msgs_app_id_reg_id(app_id,reg_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE DATABASE IF NOT EXISTS push_shard2 DEFAULT CHARACTER SET utf8;
USE push_shard2;

DROP TABLE IF EXISTS offline_msgs;
CREATE TABLE offline_msgs (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  reg_id VARCHAR(36) NOT NULL, 
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY offline_msgs_app_id_reg_id(app_id,reg_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE DATABASE IF NOT EXISTS push_shard3 DEFAULT CHARACTER SET utf8;
USE push_shard3;

DROP TABLE IF EXISTS offline_msgs;
CREATE TABLE offline_msgs (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  reg_id VARCHAR(36) NOT NULL,
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY offline_msgs_app_id_reg_id(app_id,reg_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE DATABASE IF NOT EXISTS push_shard4 DEFAULT CHARACTER SET utf8;
USE push_shard4;

DROP TABLE IF EXISTS offline_msgs;
CREATE TABLE offline_msgs (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  reg_id VARCHAR(36) NOT NULL,
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY offline_msgs_app_id_reg_id(app_id,reg_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE DATABASE IF NOT EXISTS push_shard5 DEFAULT CHARACTER SET utf8;
USE push_shard5;

DROP TABLE IF EXISTS offline_msgs;
CREATE TABLE offline_msgs (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  reg_id VARCHAR(36) NOT NULL, 
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY offline_msgs_app_id_reg_id(app_id,reg_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE DATABASE IF NOT EXISTS push_shard6 DEFAULT CHARACTER SET utf8;
USE push_shard6;

DROP TABLE IF EXISTS offline_msgs;
CREATE TABLE offline_msgs (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  reg_id VARCHAR(36) NOT NULL,
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY offline_msgs_app_id_reg_id(app_id,reg_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE DATABASE IF NOT EXISTS push_shard7 DEFAULT CHARACTER SET utf8;
USE push_shard7;
DROP TABLE IF EXISTS offline_msgs;
CREATE TABLE offline_msgs (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  reg_id VARCHAR(36) NOT NULL,
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY offline_msgs_app_id_reg_id(app_id,reg_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE DATABASE IF NOT EXISTS push_shard8 DEFAULT CHARACTER SET utf8;
USE push_shard8;

DROP TABLE IF EXISTS offline_msgs;
CREATE TABLE offline_msgs (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  reg_id VARCHAR(36) NOT NULL,
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY offline_msgs_app_id_reg_id(app_id,reg_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE DATABASE IF NOT EXISTS push_shard9 DEFAULT CHARACTER SET utf8;
USE push_shard9;

DROP TABLE IF EXISTS offline_msgs;
CREATE TABLE offline_msgs (
  id VARCHAR(36) NOT NULL,
  app_id VARCHAR(36) NOT NULL,
  reg_id VARCHAR(36) NOT NULL,
  kind INT(4) UNSIGNED NOT NULL, -- 消息类型
  content VARCHAR(500) NOT NULL,
  extra VARCHAR(500) NOT NULL, -- json格式，例如：{"order_id":"123"}
  created_at BIGINT(20) NOT NULL,
  PRIMARY KEY (id),
  KEY offline_msgs_app_id_reg_id(app_id,reg_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


