CREATE TABLE `file` (
        `id` int NOT NULL AUTO_INCREMENT,
        `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
        `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
        `file_size` bigint DEFAULT '0' COMMENT '文件大小',
        `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
        `create_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
        `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
        `status` tinyint(1) unsigned zerofill NOT NULL DEFAULT '0' COMMENT '状态(0：正常，1：禁用)',
        `is_delete` tinyint(1) unsigned zerofill NOT NULL DEFAULT '0' COMMENT '是否已经删除(0：否，1：是)',
        `ext2` text COMMENT '备用字段2',
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_file_hash` (`file_sha1`),
        KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8mb3;

CREATE TABLE `user` (
        `id` int NOT NULL AUTO_INCREMENT,
        `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
        `user_nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户昵称',
        `user_pwd` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户encoded密码',
        `salt` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码盐',
        `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '用户头像',
        `sex` tinyint(1) NOT NULL DEFAULT '1' COMMENT '性别;0:女,1:男,2:保密',
        `email` varchar(64) DEFAULT '' COMMENT '邮箱',
        `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '手机号',
        `email_validated` tinyint(1) DEFAULT '0' COMMENT '邮箱是否已验证',
        `phone_validated` tinyint(1) DEFAULT '0' COMMENT '手机号是否已验证',
        `sign_up_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
        `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
        `profile` text COMMENT '用户属性',
        `status` int NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_username` (`user_name`),
        KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `user_file` (
         `id` int NOT NULL AUTO_INCREMENT,
         `user_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
         `file_sha1` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文件hash',
         `file_size` bigint DEFAULT '0' COMMENT '文件大小',
         `file_name` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文件名',
         `upload_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
         `last_update` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
         `status` tinyint(1) unsigned zerofill NOT NULL DEFAULT '0' COMMENT '文件状态(0：正常，1：禁用)',
         `is_delete` tinyint(1) unsigned zerofill NOT NULL DEFAULT '0' COMMENT '是否已经删除(0：否，1：是)',
         PRIMARY KEY (`id`),
         UNIQUE KEY `idx_user_file` (`user_name`,`file_sha1`),
         KEY `idx_status` (`status`),
         KEY `idx_user_id` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=54 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;