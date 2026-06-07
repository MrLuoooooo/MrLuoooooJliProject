CREATE DATABASE IF NOT EXISTS im_server CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE im_server;

CREATE TABLE IF NOT EXISTS `accounts` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `account` varchar(45) DEFAULT NULL COMMENT '登录名',
  `password` varchar(45) DEFAULT NULL COMMENT '密码',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '修改时间',
  `state` tinyint DEFAULT '0' COMMENT '状态',
  `role_type` tinyint DEFAULT 0 COMMENT '角色类型',
  `parent_account` varchar(45) DEFAULT NULL COMMENT '父账号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_account` (`account`),
  KEY `idx_parent` (`parent_account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '管理端-账号';
CREATE TABLE IF NOT EXISTS `accountapprels` (
  `id` int NOT NULL AUTO_INCREMENT,
  `app_key` varchar(20) DEFAULT '',
  `account_id` int DEFAULT '0',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  UNIQUE KEY `uniq_app` (`account_id`,`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
CREATE TABLE IF NOT EXISTS `androidpushconfs` (
  `app_key` varchar(20) DEFAULT NULL COMMENT 'appKey',
  `push_channel` varchar(10) DEFAULT NULL COMMENT '推送渠道',
  `package` varchar(100) DEFAULT NULL COMMENT 'package',
  `push_conf` varchar(500) DEFAULT NULL COMMENT '配置参数',
  `push_ext` MEDIUMBLOB DEFAULT NULL COMMENT '其他配置',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  UNIQUE KEY `uniq_channel` (`app_key`,`push_channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'android推送配置';
CREATE TABLE IF NOT EXISTS `appexts` (
  `app_key` varchar(50) DEFAULT NULL COMMENT '应用key',
  `app_item_key` varchar(50) DEFAULT NULL COMMENT '参数key',
  `app_item_value` varchar(2048) DEFAULT NULL COMMENT '参数value',
  UNIQUE KEY `IDX_APPKEY_APPITEMKEY` (`app_key`,`app_item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'app应用-补充配置';
CREATE TABLE IF NOT EXISTS `apps` (
  `app_key` varchar(45) NOT NULL COMMENT '应用key',
  `app_secret` varchar(45) NOT NULL COMMENT '应用密钥 16位',
  `app_secure_key` varchar(45) NOT NULL COMMENT '安全key 16位',
  `app_status` tinyint DEFAULT '0' COMMENT '状态',
  `app_type` tinyint DEFAULT '0' COMMENT '应用类型',
  `app_name` varchar(100) DEFAULT NULL COMMENT '应用名称',
  UNIQUE KEY `uniq_appkey` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'app应用';
CREATE TABLE IF NOT EXISTS `banusers` (
  `user_id` varchar(32) NOT NULL COMMENT '用户id',
  `end_time` bigint DEFAULT '0' COMMENT '结束时间 为0 为永久封禁',
  `scope_key` varchar(20) NOT NULL DEFAULT 'default' COMMENT '封禁范围 default 用户封禁；platform:该用户指定的平台封禁；device:该用户指定的设备封禁;ip:该用户指定的ip封禁',
  `scope_value` varchar(1000) DEFAULT '与scope_key配合使用',
  `ext` varchar(100) DEFAULT NULL COMMENT '封禁时携带的扩展信息',
  `app_key` varchar(20) NOT NULL COMMENT '应用key',
  UNIQUE KEY `uniq_appkey_userid` (`app_key`,`user_id`,`scope_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '禁用用户';
CREATE TABLE IF NOT EXISTS `bc_hismsgs` (
  `conver_id` varchar(100) NOT NULL COMMENT '会话id',
  `sender_id` varchar(32) DEFAULT NULL COMMENT '发送者id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_type` varchar(50) DEFAULT NULL COMMENT '消息类型',
  `msg_id` varchar(20) NOT NULL COMMENT '消息id',
  `send_time` bigint DEFAULT NULL COMMENT '发送时间',
  `msg_seq_no` int DEFAULT NULL COMMENT '消息seq',
  `msg_body` mediumblob COMMENT '消息体',
  KEY `idx_msgid` (`app_key`,`conver_id`,`msg_id`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '广播-历史消息表';
CREATE TABLE IF NOT EXISTS `blocks` (
  `user_id` varchar(32) DEFAULT NULL COMMENT '操作人',
  `block_user_id` varchar(32) DEFAULT NULL   COMMENT '被锁定者id',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  UNIQUE KEY `uniq_appkey_userid` (`app_key`,`user_id`,`block_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '锁定记录';
CREATE TABLE IF NOT EXISTS `brdinboxmsgs` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `app_key` varchar(32) DEFAULT NULL COMMENT '应用key',
  KEY `idx_sendtime` (`app_key`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '广播-收件箱';
CREATE TABLE IF NOT EXISTS `chatroominfos` (
  `chat_id` varchar(32) DEFAULT NULL COMMENT '聊天室id',
  `chat_name` varchar(45) DEFAULT NULL COMMENT '聊天室名称',
  UNIQUE KEY `uniq_chatid` (`app_key`,`chat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '聊天室信息';
CREATE TABLE IF NOT EXISTS `cmdinboxmsgs` (
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `uniq_tag` varchar(100) DEFAULT NULL COMMENT '唯一标签',
  UNIQUE KEY `uniq_tag` (`app_key`,`user_id`,`uniq_tag`),
  KEY `idx_appkey_time` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'cmd收件箱';
CREATE TABLE IF NOT EXISTS `cmdsendboxmsgs` (
  KEY `idx_appkey_userid_time` (`app_key`,`user_id`,`send_time`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'cmd发件箱';
CREATE TABLE IF NOT EXISTS `convercleantimes` (
  `conver_id` varchar(100) DEFAULT NULL COMMENT '会话id',
  `sub_channel` varchar(32) DEFAULT '',
  `channel_type` tinyint DEFAULT '0' COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `clean_time` bigint DEFAULT '0' COMMENT '清除时间',
  UNIQUE KEY `uniq_destroy` (`app_key`,`conver_id`,`sub_channel`,`channel_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '会话清理记录';
CREATE TABLE IF NOT EXISTS `conversations` (
  `latest_msg_id` varchar(20) DEFAULT NULL COMMENT '最新消息id',
  `latest_msg` mediumblob COMMENT '最新消息体',
  `latest_unread_msg_index` int DEFAULT '0' COMMENT '最新未读index',
  `latest_read_msg_index` int DEFAULT '0' COMMENT '最新已读消息index',
  `latest_read_msg_id` varchar(20) DEFAULT NULL COMMENT '最新已读消息id',
  `latest_read_msg_time` bigint DEFAULT '0' COMMENT '最新已读时间',
  `sort_time` bigint DEFAULT '0' COMMENT 'sort time',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除 0 未删除，1已删除',
  `is_top` tinyint DEFAULT '0' COMMENT '是否置顶 0未指定，1置顶',
  `top_updated_time` bigint DEFAULT '0' COMMENT '置顶更新时间',
  `undisturb_type` tinyint DEFAULT '0' COMMENT '免打扰类型：0:取消免打扰；1:普通会话免打扰；',
  `sync_time` bigint DEFAULT '0' COMMENT '同步消息位点',
  `unread_tag` tinyint DEFAULT '0' COMMENT '未读tag',
  `conver_exts` mediumblob,
  UNIQUE KEY `uniq_app_key_user_id_target_id` (`app_key`,`user_id`,`target_id`,`sub_channel`,`channel_type`),
  KEY `idx_sync_time` (`app_key`,`user_id`,`sync_time`),
  KEY `idx_update_time` (`app_key`,`user_id`,`sort_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '会话';
CREATE TABLE IF NOT EXISTS `userconvertags` (
  `user_id` VARCHAR(32) NULL COMMENT '用户id',
  `tag` VARCHAR(50) NULL COMMENT 'tag',
  `tag_name` VARCHAR(50) NULL COMMENT '分组名称',
  `created_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `app_key` VARCHAR(20) NULL COMMENT '应用key',
  UNIQUE INDEX `uniq_tag` (`app_key`, `user_id`, `tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '会话-分组';
CREATE TABLE IF NOT EXISTS `convertagrels` (
  `target_id` VARCHAR(32) NULL COMMENT '目标id',
  `channel_type` TINYINT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  UNIQUE INDEX `uniq_tag_target` (`app_key`, `user_id`, `tag`, `target_id`, `channel_type`, `sub_channel`),
  KEY `idx_target` (`app_key`, `user_id`, `target_id`, `channel_type`, `sub_channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '会话-分组绑定关系表';
CREATE TABLE IF NOT EXISTS `msgstats` (
  `stat_type` TINYINT NULL DEFAULT 0 COMMENT '统计类型 1上行消息，2分发，3下行消息',
  `time_mark` BIGINT NULL COMMENT '标记时间',
  `count` INT NULL COMMENT '数量',
  UNIQUE INDEX `uniq_mark` (`app_key` ASC, `stat_type` ASC, `channel_type` ASC, `time_mark` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '统计-消息统计';
CREATE TABLE IF NOT EXISTS `useractivities` (
  `time_mark` BIGINT NULL COMMENT '标记点',
  UNIQUE INDEX `uniq_userid` (`app_key` ASC, `time_mark` ASC, `user_id` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '统计-用户活跃统计';
CREATE TABLE IF NOT EXISTS `connectcounts` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `connect_type` TINYINT NULL DEFAULT 0,
  `time_mark` BIGINT NULL,
  `count` INT NULL,
  `app_key` VARCHAR(20) NULL,
  UNIQUE INDEX `uniq_mark` (`app_key`, `connect_type`, `time_mark`)
CREATE TABLE IF NOT EXISTS `fileconfs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `app_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '应用key',
  `channel` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'oss供应商',
  `conf` json DEFAULT NULL COMMENT '配置',
  `enable` tinyint(1) DEFAULT '0' COMMENT '是否可用 0不可用，1可用',
  `updated_time` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  UNIQUE KEY `app_key` (`app_key`,`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT = '文件配置';
CREATE TABLE IF NOT EXISTS `g_delhismsgs` (
  `target_id` varchar(32) DEFAULT NULL COMMENT '目标id',
  `msg_time` bigint DEFAULT 0 COMMENT '消息时间',
  `msg_seq` int DEFAULT 0 COMMENT '消息seq',
  `effective_time` bigint DEFAULT 0,
  UNIQUE KEY `uniq_msgid` (`app_key`,`user_id`,`target_id`,`sub_channel`,`msg_id`),
  KEY `idx_target` (`app_key`,`user_id`,`target_id`,`sub_channel`,`msg_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群聊-历史消息已删记录表';
CREATE TABLE IF NOT EXISTS `g_hismsgs` (
  `receiver_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `member_count` int DEFAULT '0' COMMENT '成员数',
  `read_count` int DEFAULT '0' COMMENT '已读数',
  `is_delete` tinyint DEFAULT '0' COMMENT '是否删除 0未删除， 1已删除',
  `is_ext` tinyint DEFAULT '0',
  `is_exset` tinyint DEFAULT '0',
  `msg_ext` mediumblob,
  `msg_exset` mediumblob,
  `destroy_time` bigint DEFAULT 0,
  `life_time_after_read` bigint DEFAULT 0,
  `is_portion` tinyint DEFAULT 0,
  KEY `idx_appkey_converid` (`app_key`,`conver_id`,`sub_channel`,`msg_id`,`send_time`),
  KEY `idx_conver_time` (`app_key`,`conver_id`,`sub_channel`,`send_time`),
  KEY `idx_sender_time` (`app_key`,`conver_id`,`sub_channel`,`sender_id`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群聊-历史消息表';
CREATE TABLE `g_portionrels` (
  `conver_id` VARCHAR(100) NULL,
  `sub_channel` VARCHAR(32) NULL DEFAULT '',
  `channel_type` TINYINT NULL DEFAULT 0,
  `user_id` VARCHAR(32) NULL,
  `msg_id` VARCHAR(32) NULL,
  `msg_time` BIGINT NULL DEFAULT 0,
  `app_key` VARCHAR(20) NULL DEFAULT '',
  UNIQUE INDEX `uniq_msgid` (`app_key`, `conver_id`, `sub_channel`, `user_id`, `msg_id`),
  INDEX `idx_msg_time` (`app_key`, `conver_id`, `sub_channel`, `user_id`, `msg_time`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
CREATE TABLE IF NOT EXISTS `gc_hismsgs` (
  KEY `idx_msg` (`app_key`,`conver_id`,`channel_type`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群聊-公告消息表';
CREATE TABLE IF NOT EXISTS `globalconfs` (
  `conf_key` varchar(50) DEFAULT NULL COMMENT '配置key',
  `conf_value` varchar(2000) DEFAULT NULL COMMENT '配置value',
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建者',
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新者',
  UNIQUE KEY `uniq_key` (`conf_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '全局配置';
CREATE TABLE IF NOT EXISTS `globalconvers` (
  `sender_id` varchar(32) DEFAULT NULL COMMENT '发送者',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者',
  `updated_time` bigint DEFAULT NULL COMMENT '更新时间',
  UNIQUE KEY `uniq_conver` (`app_key`,`conver_id`,`sub_channel`,`channel_type`),
  KEY `idx_time` (`app_key`,`channel_type`,`updated_time`),
  KEY `idx_targetid` (`app_key`,`channel_type`,`target_id`,`updated_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '全局会话表';
CREATE TABLE IF NOT EXISTS `groupinfoexts` (
  `group_id` varchar(32) DEFAULT NULL COMMENT '群id',
  `item_key` varchar(50) DEFAULT NULL COMMENT '参数key',
  `item_value` varchar(100) DEFAULT NULL COMMENT '参数value',
  `item_type` tinyint DEFAULT '0' COMMENT '参数类型',
  UNIQUE KEY `uniq_appkey_groupid` (`app_key`,`group_id`,`item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群补充信息';
CREATE TABLE IF NOT EXISTS `groupinfos` (
  `group_id` varchar(64) DEFAULT NULL COMMENT '群id',
  `group_name` varchar(64) DEFAULT NULL COMMENT '群名称',
  `group_portrait` varchar(200) DEFAULT NULL COMMENT '群头像',
  `is_mute` tinyint DEFAULT '0' COMMENT '是否全局禁言，0：否；1：是；',
  UNIQUE KEY `uniq_appkey_groupid` (`app_key`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群基本信息';
CREATE TABLE IF NOT EXISTS `groupmemberexts` (
  `member_id` varchar(32) DEFAULT NULL COMMENT '群成员id',
  UNIQUE KEY `uniq_item_key` (`app_key`,`group_id`,`member_id`,`item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群成员-补充信息';
CREATE TABLE IF NOT EXISTS `groupmembers` (
  `member_id` varchar(64) DEFAULT NULL COMMENT '成员id',
  `member_type` tinyint DEFAULT '0' COMMENT '成员类型',
  `app_key` varchar(45) DEFAULT NULL COMMENT '应用key',
  `is_allow` tinyint DEFAULT '0' COMMENT '是否白名单 0:非白名单用户；1:白名单用户；',
  `mute_end_at` bigint DEFAULT '0' COMMENT '禁言结束时间戳',
  UNIQUE KEY `uniq_appkey_grpid_memid` (`app_key`,`group_id`,`member_id`),
  KEY `idx_memberid` (`app_key`,`member_id`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群成员';
CREATE TABLE IF NOT EXISTS `grpassistantrels` (
  `id` int NOT NULL COMMENT '主键id',
  `assistant_id` varchar(32) DEFAULT NULL COMMENT '助手id',
  UNIQUE KEY `uniq_target` (`assistant_id`,`target_id`,`channel_type`,`app_key`)
CREATE TABLE IF NOT EXISTS `grpsnapshots` (
  `group_id` varchar(32) NOT NULL COMMENT '群id',
  `created_time` bigint DEFAULT '0' COMMENT '创建时间',
  `snapshot` mediumblob COMMENT '详情',
  KEY `idx_group_id` (`app_key`,`group_id`,`created_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群快照';
CREATE TABLE IF NOT EXISTS `ic_conditions` (
  `channel_type` varchar(100) DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_type` varchar(1000) DEFAULT NULL COMMENT '消息类型',
  `sender_id` varchar(1000) DEFAULT NULL COMMENT '发送者id',
  `receiver_id` varchar(1000) DEFAULT NULL COMMENT '接收者id',
  `interceptor_id` int DEFAULT NULL,
  KEY `idx_icid` (`app_key`,`interceptor_id`)
CREATE TABLE IF NOT EXISTS `inboxmsgs` (
  KEY `IDX_USERID_MSG` (`app_key`,`user_id`,`send_time`),
) ENGINE=InnoDB AUTO_INCREMENT=45612 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '收件箱';
CREATE TABLE IF NOT EXISTS `interceptors` (
  `name` varchar(50) DEFAULT NULL COMMENT '名称',
  `sort` int NOT NULL DEFAULT '0',
  `request_url` varchar(500) DEFAULT NULL COMMENT '请求地址',
  `request_template` text COMMENT '请求模板',
  `succ_template` varchar(200) DEFAULT NULL COMMENT '成功模板',
  `is_async` tinyint DEFAULT '0' COMMENT '是否异步',
  `conf` varchar(2000) DEFAULT NULL,
  `intercept_type` tinyint DEFAULT '0',
  KEY `idx_sort` (`app_key`,`sort`)
CREATE TABLE IF NOT EXISTS `ioscertificates` (
  `package` varchar(100) DEFAULT NULL,
  `certificate` mediumblob,
  `cert_pwd` varchar(50) DEFAULT NULL COMMENT '认证密码',
  `is_product` tinyint DEFAULT '0',
  `cert_path` varchar(255) DEFAULT NULL COMMENT 'cert存入路径',
  `voip_cert` mediumblob COMMENT 'voip certificate',
  `voip_cert_pwd` varchar(50) DEFAULT NULL COMMENT 'voip cert password',
  `voip_cert_path` varchar(255) DEFAULT NULL COMMENT 'voip cert path',
  UNIQUE KEY `uniq_package` (`app_key`)
CREATE TABLE IF NOT EXISTS `mentionmsgs` (
  `mention_type` tinyint DEFAULT NULL COMMENT '类型',
  `msg_time` bigint DEFAULT NULL COMMENT '消息时间',
  `msg_index` int DEFAULT NULL COMMENT '消息index',
  `is_read` tinyint DEFAULT NULL COMMENT '是否已读',
  KEY `idx_uid_tid_type` (`app_key`,`user_id`,`target_id`,`channel_type`,`sub_channel`,`msg_index`,`msg_time`),
  KEY `idx_read` (`app_key`,`user_id`,`target_id`,`channel_type`,`sub_channel`,`is_read`,`msg_time`),
  KEY `idx_user_msgid` (`app_key`,`user_id`,`target_id`,`channel_type`,`sub_channel`,`msg_id`),
  KEY `idx_target_msgid` (`app_key`, `msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '@消息记录';
CREATE TABLE IF NOT EXISTS `mergedmsgs` (
  `parent_msg_id` varchar(20) DEFAULT NULL COMMENT '父消息id',
  `from_id` varchar(32) DEFAULT NULL COMMENT '发送者',
  `msg_time` bigint DEFAULT '0' COMMENT '消息时间',
  UNIQUE KEY `idx_appkey_pmsg` (`app_key`,`parent_msg_id`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '聚合消息';
CREATE TABLE IF NOT EXISTS `msgexts` (
  `key` varchar(50) DEFAULT NULL COMMENT '参数key',
  `value` varchar(1000) DEFAULT NULL COMMENT '参数value',
  `app_key` varchar(45) DEFAULT NULL COMMENT '更新时间',
  UNIQUE KEY `uniq_msgid` (`app_key`,`msg_id`,`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '消息-补充内容';
CREATE TABLE IF NOT EXISTS `msgexsets` (
  `msg_id` VARCHAR(20) NULL COMMENT '消息id',
  `key` VARCHAR(50) NULL COMMENT 'key',
  `item` VARCHAR(100) NULL,
  UNIQUE INDEX `uniq_item` (`app_key`, `msg_id`, `key`, `item`)
CREATE TABLE IF NOT EXISTS `p_delhismsgs` (
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '单聊-历史消息已删记录表';
CREATE TABLE IF NOT EXISTS `p_hismsgs` (
  `receiver_id` varchar(32) DEFAULT NULL COMMENT '接收者',
  `is_read` tinyint DEFAULT '0' COMMENT '是否已读 0未读，1已读',
  `read_time` bigint DEFAULT 0 COMMENT '消息已读时间', 
  `is_delete` tinyint DEFAULT '0' COMMENT '是否删除 0未删除，1已删除',
  KEY `idx_app_key_conver_id` (`app_key`,`conver_id`,`sub_channel`,`msg_id`,`send_time`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '单聊-历史消息';
CREATE TABLE IF NOT EXISTS `pushtokens` (
  `device_id` varchar(200) DEFAULT NULL COMMENT '设备id',
  `platform` varchar(10) DEFAULT NULL COMMENT '平台',
  `package` varchar(200) DEFAULT NULL,
  `push_token` varchar(200) DEFAULT NULL COMMENT '推送token',
  `voip_token` varchar(200) DEFAULT NULL COMMENT 'voip推送token',
  UNIQUE KEY `idx_user_id` (`app_key`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '推送token';
CREATE TABLE IF NOT EXISTS `readinfos` (
  UNIQUE KEY `uniq_member` (`app_key`,`channel_type`,`group_id`,`msg_id`,`member_id`,`sub_channel`),
  KEY `idx_memberid` (`app_key`,`channel_type`,`group_id`,`sub_channel`,`member_id`,`msg_id`),
  KEY `idx_msgid` (`app_key`,`channel_type`,`group_id`,`sub_channel`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群已读表';
CREATE TABLE IF NOT EXISTS `s_hismsgs` (
  KEY `idx_appkey_converid` (`app_key`,`conver_id`,`send_time`)
CREATE TABLE IF NOT EXISTS `sendboxmsgs` (
  `target_id` varchar(45) DEFAULT NULL COMMENT '接收者id',
  KEY `idx_user_id_send_time` (`app_key`,`user_id`,`send_time`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '发件箱';
CREATE TABLE IF NOT EXISTS `sensitivewords` (
  `app_key` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '应用key',
  `word` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '敏感词',
  `word_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '敏感词过滤类型。1：拦截敏感词；2：替换敏感词；',
  UNIQUE KEY `uniq_word` (`app_key`,`word`),
  KEY `idx_appkey` (`app_key`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT = '敏感词';
CREATE TABLE IF NOT EXISTS `subrelations` (
  `subscriber` varchar(32) DEFAULT NULL COMMENT '订阅者',
  UNIQUE KEY `uniq_sub` (`app_key`,`user_id`,`subscriber`)
CREATE TABLE IF NOT EXISTS `usercleantimes` (
  `clean_time` bigint DEFAULT NULL COMMENT '清理时间',
  UNIQUE KEY `uniq_app_key_user_id_target_id` (`app_key`,`user_id`,`target_id`,`sub_channel`,`channel_type`)
CREATE TABLE IF NOT EXISTS `userexts` (
  `item_value` varchar(2000) DEFAULT NULL COMMENT '参数value',
  UNIQUE KEY `uniq_item_key` (`app_key`,`user_id`,`item_key`),
  KEY `idx_item_key` (`app_key`,`item_key`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '用户-补充信息表';
CREATE TABLE IF NOT EXISTS `users` (
  `user_type` tinyint DEFAULT '0' COMMENT '用户类型 0用户, 1机器人',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `user_portrait` varchar(200) DEFAULT NULL COMMENT '用户头像',
  `pinyin` varchar(50) DEFAULT NULL,
  `phone` varchar(50) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `login_account` varchar(50) DEFAULT NULL,
  `login_pass` varchar(50) DEFAULT NULL,
  UNIQUE KEY `uniq_userid` (`app_key`,`user_id`),
  UNIQUE KEY `uniq_phone` (`app_key`,`phone`),
  UNIQUE KEY `uniq_email` (`app_key`,`email`),
  UNIQUE KEY `uniq_account` (`app_key`,`login_account`),
  KEY `idx_userid` (`app_key`,`user_type`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '用户-基本信息表';
CREATE TABLE IF NOT EXISTS `clientlogs` (
  `start` BIGINT NULL COMMENT '开始时间',
  `end` BIGINT NULL COMMENT '结束时间',
  `log` MEDIUMBLOB NULL COMMENT '详情',
  `state` TINYINT NULL DEFAULT 0 COMMENT '状态',
  `platform` VARCHAR(20) NULL COMMENT '客户端类型 ios, android, web',
  `device_id` VARCHAR(100) NULL COMMENT '设备id',
  `log_url` VARCHAR(200) NULL COMMENT '日志链接',
  `trace_id` VARCHAR(50) NULL COMMENT '跟踪id',
  `fail_reason` VARCHAR(100) NULL COMMENT '失败原因',
  `description` VARCHAR(100) NULL COMMENT '描述',
  INDEX `idx_userid` (`app_key` ASC, `user_id` ASC),
  UNIQUE KEY `uniq_msgid` (`app_key`,`msg_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '客户端日志表';
CREATE TABLE IF NOT EXISTS `rtcrooms` (
  `room_id` varchar(50) DEFAULT NULL COMMENT '房间id',
  `room_type` tinyint DEFAULT '0' COMMENT '房间类型',
  `rtc_channel` tinyint DEFAULT '0',
  `rtc_media_type` tinyint DEFAULT '0',
  `owner_id` varchar(32) DEFAULT NULL COMMENT '创建者',
  `ext` varchar(2000) DEFAULT NULL COMMENT '扩展字段',
  `accepted_time` bigint DEFAULT '0' COMMENT '1v1 接通时间',
  UNIQUE KEY `uniq_roomid` (`app_key`,`room_id`),
  KEY `idx_conver` (`app_key`,`conver_id`,`channel_type`,`sub_channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'rtc-房间表';
CREATE TABLE IF NOT EXISTS `rtcmembers` (
  `member_id` varchar(32) DEFAULT NULL COMMENT '成员id',
  `device_id` varchar(50) DEFAULT NULL COMMENT '设备id',
  `rtc_state` tinyint DEFAULT '0'  COMMENT '状态',
  `inviter_id` varchar(32) DEFAULT NULL COMMENT '邀请人id',
  `latest_ping_time` bigint DEFAULT '0' COMMENT '最新ping时间',
  `call_time` bigint DEFAULT '0' COMMENT '拨号时间',
  `connect_time` bigint DEFAULT '0' COMMENT '通话开始时间',
  `hangup_time` bigint DEFAULT '0' COMMENT '通话挂断时间',
  UNIQUE KEY `uniq_member` (`app_key`,`room_id`,`member_id`),
  KEY `idx_room` (`app_key`,`member_id`,`room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'rtc-成员表';
CREATE TABLE IF NOT EXISTS `msgtransconfs` (
  `json_path` varchar(200) DEFAULT NULL,
  `app_key` varchar(20) DEFAULT NULL COMMENT '租户key',
  UNIQUE KEY `uniq_path` (`app_key`,`msg_type`,`json_path`)
CREATE TABLE IF NOT EXISTS `i18nkeys` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `lang` VARCHAR(20) NULL,
  `value` VARCHAR(200) NULL COMMENT 'value',
  `app_key` VARCHAR(20) NULL COMMENT '租户key',
  UNIQUE INDEX `uniq_key` (`app_key`, `lang`, `key`)
CREATE TABLE IF NOT EXISTS `friendrels` (
  `friend_id` varchar(32) DEFAULT NULL COMMENT '朋友userId',
  `display_name` varchar(50) DEFAULT '' COMMENT '好友备注名',
  `order_tag` varchar(20) NULL DEFAULT '' COMMENT '好友的排序标识',
  UNIQUE KEY `uniq_friend` (`app_key`,`user_id`,`friend_id`),
  KEY `idx_order` (`app_key`, `user_id`, `order_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT '好友绑定关系表';
CREATE TABLE IF NOT EXISTS `botconvers` (
  `conver_key` VARCHAR(100) NULL DEFAULT '',
  `conver_type` TINYINT NULL DEFAULT 0,
  `conver_id` VARCHAR(50) NULL DEFAULT '',
  `app_key` VARCHAR(50) NULL DEFAULT '',
  `updated_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE INDEX `uniq_key` (`app_key`, `conver_type`, `conver_key`)
CREATE TABLE IF NOT EXISTS `favoritemsgs` (
  `user_id` VARCHAR(50) NULL,
  `sender_id` VARCHAR(50) NULL,
  `receiver_id` VARCHAR(50) NULL,
  `sub_channel` VARCHAR(32) DEFAULT '',
  `channel_type` TINYINT NULL,
  `msg_id` VARCHAR(50) NULL,
  `msg_time` BIGINT DEFAULT '0',
  `msg_type` VARCHAR(50) NULL,
  `msg_body` MEDIUMBLOB,
  `created_time` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  INDEX `idx_userid` (`app_key`, `user_id`, `created_time`),
  UNIQUE KEY `uniq_msgid` (`app_key`, `user_id`, `msg_id`)
CREATE TABLE IF NOT EXISTS `topmsgs` (
  `conver_id` varchar(100) DEFAULT '',
  `channel_type` tinyint DEFAULT '0',
  `msg_id` varchar(20) DEFAULT NULL,
  `user_id` varchar(32) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE KEY `idx_msg` (`app_key`,`conver_id`,`sub_channel`,`channel_type`)
INSERT IGNORE INTO `globalconfs` (`conf_key`,`conf_value`)VALUES('jimdb_version','20251102');
INSERT IGNORE INTO `accounts`(`account`,`password`)VALUES('admin','7c4a8d09ca3762af61e59520943dc26494f8941b');
-- 自动创建社区论坛应用
INSERT IGNORE INTO apps (app_key, app_secret, app_secure_key, app_status, app_name) VALUES ("community", "community-secret", "community-secure", 0, "论坛");
INSERT IGNORE INTO accountapprels (app_key, account_id) VALUES ("community", 1);
