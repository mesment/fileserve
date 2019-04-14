CREATE TABLE `filemeta` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `file_md5` char(64) NOT NULL DEFAULT '''''' COMMENT '文件md5',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_addr` varchar(512) NOT NULL COMMENT '文件存放地址',
  `file_size` bigint(20) NOT NULL COMMENT '文件大小',
  `create_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '0-正常 1-删除',
  `memo` varchar(100) DEFAULT '' COMMENT '备注',
  `resv` varchar(100) DEFAULT '' COMMENT '预留',
  PRIMARY KEY (`id`),
  UNIQUE KEY `file_md5` (`file_md5`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_passwd` varchar(256) NOT NULL DEFAULT '' COMMENT '加密密码',
  `email` varchar(64) DEFAULT '''''' COMMENT '邮箱',
  `phone` varchar(32) DEFAULT '''''' COMMENT '手机号',
  `email_validated` tinyint(1) DEFAULT '0' COMMENT '邮箱是否已验证0-否 1-是',
  `phone_validated` tinyint(1) DEFAULT '0' COMMENT '手机是否已验证0-否 1-是',
  `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
  `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
  `profile` text COMMENT '用户资料',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '启用状态0-正常 1-禁用 2-锁定 3-标记删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;