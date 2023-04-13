SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for example
-- ----------------------------
DROP TABLE IF EXISTS `example`;
CREATE TABLE `example` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id，自增序列',
    `uuid` varchar(255) NOT NULL DEFAULT '' COMMENT '一标识',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
    `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '类型，1:类型1 2:类型2',
    `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
    `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uniq_uuid` (`uuid`) USING BTREE,
    UNIQUE KEY `uniq_name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='示例表';

SET FOREIGN_KEY_CHECKS = 1;